package main

import (
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/looplab/fsm"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	pb "observer/finance"
	"time"
)

//
const (
	ESTART      = "start"
	ERETRY      = "retry"
	EHEIGHT     = "height"
	EBLOCK      = "block"
	EFETCHERROR = "fetcherror"

	SSTOPED      = "stoped"
	SFETCHHEIGHT = "fetchheight"
	SFETCHBLOCK  = "fetchblock"
	SFETCHERROR  = "fetcherror"
	SENTER       = "enter_state"

	DInterval      = 2
	DRetryInterval = 2
)

// ObserverMechine comment
type ObserverMechine struct {
	Finite        *fsm.FSM
	curHeight     int64
	desHeight     int64
	nextEvent     string
	addrMgr       AddressManager
	rpcConfig     *rpcclient.ConnConfig
	isStop        bool
	net           *chaincfg.Params
	GRPCHost      string
	db            *Database
	interval      int //如果已经扫描到最新块高，间隔多长时间(s)进行扫描
	retryInterval int //请求失败重试时间间隔(s)
}

// DepositTask deposit task
type DepositTask struct {
	Symbol    string `json:"symbol"`
	AddressTo string `json:"addressto"`
	Txid      string `json:"txid"`
	Amount    int64  `json:"amount"`
	BlockHash string `json:"blockhash"`
}

// Run 启动状态机
func (s *ObserverMechine) Run(ch *chan error) {
	go func() {
		if s.curHeight == -1 {
			height, err := s.db.GetLastHeight()
			if err != nil {
				*ch <- err
			}
			s.curHeight = int64(height)
		}
		// load task

		s.nextEvent = ESTART
		for {
			s.Finite.Event(s.nextEvent)
			//time.Sleep(time.Duration(s.interval) * time.Second)
			if s.isStop {
				break
			}
		}
		// stoping
		if err := s.flush(); err != nil {
			log.Println(err)
		}
		s.db.Close()
		*ch <- nil
	}()
}

// Stop stop machine
func (s *ObserverMechine) Stop() {
	s.isStop = true
}

func (s *ObserverMechine) flush() error {
	return s.db.UpdataLastHeight(int(s.curHeight))
}

func (s *ObserverMechine) fetchHeight() (int64, error) {
	client, err := rpcclient.New(s.rpcConfig, nil)
	if err != nil {
		return -1, err
	}
	defer client.Shutdown()
	return client.GetBlockCount()
}

func (s *ObserverMechine) fetchBlock(height int64) error {
	client, err := rpcclient.New(s.rpcConfig, nil)
	if err != nil {
		return err
	}
	defer client.Shutdown()
	hash, err := client.GetBlockHashAsync(height).Receive()
	if err != nil {
		return err
	}
	block, err := client.GetBlockAsync(hash).Receive()
	if err != nil {
		return err
	}
	log.Printf("fetch block, hash=%s, height=%d\n", block.BlockHash().String(), height)
	// logger.WithFields(logrus.Fields{
	// 	"hash":   block.BlockHash().String(),
	// 	"height": height,
	// }).Info("fetch block")
	s.scanBlock(block)
	// notify
	s.sendDepositNotify()
	return nil
}

func (s *ObserverMechine) sendDepositNotify() {
	urlTemp := "http://%s:%s@%s"
	url := fmt.Sprintf(urlTemp, s.rpcConfig.User, s.rpcConfig.Pass, s.rpcConfig.Host)
	notifyTxs, err := s.db.GetTaskList()
	if err != nil {
		panic(err)
	}
	for _, task := range notifyTxs {
		exblock, err := exGetBlock(url, task.BlockHash)
		if err != nil {
			log.Println(err)
			continue
		}
		request := pb.DepositRequest{
			Timestamp: time.Now().Unix(),
			Symbol:    "btc",
			Txid:      task.Txid,
			AddressTo: task.AddressTo,
			Amount:    btcutil.Amount(task.Amount).String(),
			Confirm:   int32(exblock.Confirmations),
		}
		conn, err := grpc.Dial(s.GRPCHost, grpc.WithInsecure())
		if err != nil {
			log.Println(err)
			continue
		}
		defer conn.Close()
		client := pb.NewFinanceServiceClient(conn)
		// Contact the server and print out its response.
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err := client.DepositNotify(ctx, &request)
		if err != nil {
			log.Println(err)
			continue
		}
		if r.Interrupt {
			if err := s.db.DeleteTask(task); err != nil {
				panic(err)
			}
		}
	}
}

func (s *ObserverMechine) fetchHeightState(e *fsm.Event) {
	height, err := s.fetchHeight()
	if err != nil {
		// TODO logger
		// logger.WithFields(logrus.Fields{
		// 	"height": height,
		// 	"error":  err,
		// }).Warning("fetch height")
		log.Printf("fetch height, height=%d, err=%s\n", height, err)
		s.nextEvent = EFETCHERROR
		return
	}
	if height >= s.curHeight {
		s.desHeight = height
		//logger.Infof("%d blocks need fetch", height-s.curHeight)
		log.Printf("%d blocks need fetch\n", height-s.curHeight)
		s.nextEvent = EBLOCK
		return
	}
	time.Sleep(10 * time.Second)
}

func (s *ObserverMechine) fetchBlockState(e *fsm.Event) {
	for s.curHeight <= s.desHeight && !s.isStop {
		if err := s.fetchBlock(s.curHeight); err != nil {
			// logger.WithFields(logrus.Fields{
			// 	"err": err,
			// }).Warning("fetch block failure")
			log.Printf("fetch block failure, err=%s\n", err)
			s.nextEvent = EFETCHERROR
			return
		}
		s.curHeight++
	}
	s.nextEvent = EHEIGHT
}

func (s *ObserverMechine) fetchErrorState(e *fsm.Event) {
	// retry
	time.Sleep((time.Duration)(s.retryInterval) * time.Second)
	s.nextEvent = ERETRY
}

func (s *ObserverMechine) debugState(e *fsm.Event) {
	// logger.WithFields(logrus.Fields{
	// 	"event": e.Event,
	// 	"src":   e.Src,
	// 	"dst":   e.Dst,
	// }).Debugln("state change")
}

// 目前版本只提供ScriptHashTy监控
func (s *ObserverMechine) scanBlock(block *wire.MsgBlock) {
	// scan
	for _, tx := range block.Transactions {
		for i, out := range tx.TxOut {
			sc, addrs, signCount, err := txscript.ExtractPkScriptAddrs(out.PkScript, s.net)
			if err != nil {
				// logger.WithFields(logrus.Fields{
				// 	"tx":    tx.TxHash().String(),
				// 	"index": i,
				// }).Error("ExtractPkScriptAddrs failure")
				log.Printf("ExtractPkScriptAddrs failure, tx=%s, index=%d, err=%s\n", tx.TxHash().String(), i, err)
				continue
			}

			if sc == txscript.ScriptHashTy && signCount == 1 && len(addrs) == 1 {
				userAddrs := s.addrMgr.Addresses()
				if _, ok := userAddrs[addrs[0].EncodeAddress()]; ok {
					log.Println(tx.TxHash().String(), i, out.Value)
					task := DepositTask{
						Symbol:    "btc",
						AddressTo: addrs[0].EncodeAddress(),
						Txid:      tx.TxHash().String(),
						Amount:    out.Value,
						BlockHash: block.BlockHash().String(),
					}
					if ok, err := s.db.Exist(task); err != nil {
						panic(err)
					} else if !ok {
						if err := s.db.AddTask(task); err != nil {
							panic(err)
						}
					} else {
						// TODO log
					}
				}
			} else {
				// TODO log
			}
		}
	}
}

// NewObserverMechine 创建状态机
func NewObserverMechine(height int, config *rpcclient.ConnConfig, db *Database, net *chaincfg.Params, mgr AddressManager, host string) *ObserverMechine {
	o := ObserverMechine{
		rpcConfig:     config,
		curHeight:     int64(height),
		desHeight:     0,
		interval:      DInterval,
		retryInterval: DRetryInterval,
		isStop:        false,
		net:           net,
		addrMgr:       mgr,
		GRPCHost:      host,
		db:            db,
	}
	o.Finite = fsm.NewFSM(
		SSTOPED,
		fsm.Events{
			{Name: ESTART, Src: []string{SSTOPED}, Dst: SFETCHHEIGHT},
			{Name: EBLOCK, Src: []string{SFETCHHEIGHT}, Dst: SFETCHBLOCK},
			{Name: EHEIGHT, Src: []string{SFETCHBLOCK}, Dst: SFETCHHEIGHT},
			{Name: EFETCHERROR, Src: []string{SFETCHHEIGHT, SFETCHBLOCK}, Dst: SFETCHERROR},
			{Name: ERETRY, Src: []string{SFETCHERROR}, Dst: SFETCHHEIGHT},
		},
		fsm.Callbacks{
			SFETCHHEIGHT: func(e *fsm.Event) { o.fetchHeightState(e) },
			SFETCHBLOCK:  func(e *fsm.Event) { o.fetchBlockState(e) },
			SFETCHERROR:  func(e *fsm.Event) { o.fetchErrorState(e) },
			SENTER:       func(e *fsm.Event) { o.debugState(e) },
		},
	)
	return &o
}
