package main

import (
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/larspensjo/config"
	"github.com/urfave/cli"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	// DefaultCount address count
	DefaultCount = 200000
	// StartPath HD wallet path
	StartPath = "/0/77"
	// UseSigWitness sigwitness
	UseSigWitness = true
	// DefaultBind defalut bind ip
	DefaultBind = ":10087"
	// DefaultGRPCHost defalut grpc host
	DefaultGRPCHost = "localhost:10089"
	addressPath     = "./addresses.txt"
)

// Configure config
type Configure struct {
	ExtendKey     string
	Start         string
	Net           *chaincfg.Params
	UseSigWitness bool
	Count         int // 地址数量
	Bind          string
	GRPCHost      string
	WalletConfig  struct {
		URL      string
		User     string
		Password string
	}
}

// Application application
type Application struct {
	app       *cli.App
	config    *Configure
	addresses map[string]string
	lastErr   error
	processor *Processor
	machine   *ObserverMechine
	scanBegin int
	dataDir   string
	db        *Database
}

// Bind retrun bind ip
func (s *Application) Bind() string {
	return s.config.Bind
}

// Addresses return manager addresses
func (s *Application) Addresses() map[string]string {
	return s.addresses
}

// Run run
func (s *Application) Run(args []string) error {
	return s.app.Run(args)
}

// Before run before application action
func (s *Application) Before(c *cli.Context) error {
	// load config
	cfgPath := c.String("config")
	path, err := Expand(cfgPath)
	if err != nil {
		s.lastErr = err
		return nil
	}
	if err := s.loadConfig(path); err != nil {
		s.lastErr = err
		return nil
	}
	// generate address
	addrs, err := GenerateAddress(s.config.ExtendKey, s.config.Count, s.config.Start, s.config.Net, s.config.UseSigWitness)
	if err != nil {
		s.lastErr = err
		return nil
	}
	s.addresses = make(map[string]string, len(addrs))
	// write file
	file, err := os.Create(addressPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	for i, v := range addrs {
		_, err := file.WriteString(strings.TrimSpace(v) + "\n")
		if err != nil {
			panic(err)
		}
		s.addresses[v] = strconv.Itoa(i)
	}
	// load database
	tmp, err := Expand(s.dataDir)
	if err != nil {
		panic(err)
	}
	db, err := NewDatabase(tmp)
	if err != nil {
		s.lastErr = err
		return nil
	}
	s.db = db
	return nil
}

// Action action
func (s *Application) Action(c *cli.Context) error {
	if s.lastErr != nil {
		log.Println(s.lastErr)
	} else {
		//start fms
		config := &rpcclient.ConnConfig{
			Host:         s.config.WalletConfig.URL,
			User:         s.config.WalletConfig.User,
			Pass:         s.config.WalletConfig.Password,
			HTTPPostMode: true,
			DisableTLS:   true,
		}
		s.machine = NewObserverMechine(s.scanBegin, config, s.db, s.config.Net, s, s.config.GRPCHost)
		ch := make(chan error, 1)
		s.machine.Run(&ch)
		if err := <-ch; err != nil {
			log.Println(err)
		}
		return nil
	}
	return nil
}

func (s *Application) loadConfig(path string) error {
	cfg, err := config.ReadDefault(path)
	if err != nil {
		return err
	}
	if cfg.HasOption("app", "bind") {
		str, err := cfg.String("app", "bind")
		if err != nil {
			return err
		}
		s.config.Bind = str
	} else {
		s.config.Bind = DefaultBind
	}
	if cfg.HasOption("app", "key") {
		key, err := cfg.String("app", "key")
		if err != nil {
			return err
		}
		s.config.ExtendKey = key
	} else {
		return fmt.Errorf("cant find extend key in config file %s", path)
	}
	if cfg.HasOption("app", "count") {
		value, err := cfg.Int("app", "count")
		if err != nil {
			return err
		}
		s.config.Count = value
	} else { // default
		s.config.Count = DefaultCount
	}
	if cfg.HasOption("app", "start") {
		value, err := cfg.String("app", "start")
		if err != nil {
			return err
		}
		s.config.Start = value
	} else {
		s.config.Start = StartPath
	}
	if cfg.HasOption("app", "testnet") {
		value, err := cfg.Bool("app", "testnet")
		if err != nil {
			return err
		}
		if value {
			s.config.Net = &chaincfg.TestNet3Params
		} else {
			s.config.Net = &chaincfg.MainNetParams
		}
	} else {
		s.config.Net = &chaincfg.MainNetParams
	}
	if cfg.HasOption("app", "sigwitness") {
		value, err := cfg.Bool("app", "sigwitness")
		if err != nil {
			return err
		}
		s.config.UseSigWitness = value
	} else {
		s.config.UseSigWitness = UseSigWitness
	}
	// load wallet option
	if cfg.HasOption("wallet", "url") {
		value, err := cfg.String("wallet", "url")
		if err != nil {
			return err
		}
		s.config.WalletConfig.URL = value
	} else {
		s.config.WalletConfig.URL = "localhost:8332"
	}
	if cfg.HasOption("wallet", "user") {
		value, err := cfg.String("wallet", "user")
		if err != nil {
			return err
		}
		s.config.WalletConfig.User = value
	} else {
		s.config.WalletConfig.User = "obs"
	}
	if cfg.HasOption("wallet", "password") {
		value, err := cfg.String("wallet", "password")
		if err != nil {
			return err
		}
		s.config.WalletConfig.Password = value
	} else {
		s.config.WalletConfig.Password = "obs"
	}

	// load grpc
	if cfg.HasOption("grpc", "host") {
		value, err := cfg.String("grpc", "host")
		if err != nil {
			return err
		}
		s.config.GRPCHost = value
	} else {
		s.config.GRPCHost = DefaultGRPCHost
	}
	return nil
}

// Stop stop application
func (s *Application) Stop() {
	s.processor.StopGRPC()
	s.machine.Stop()
}

// NewApp new application
func NewApp() *Application {
	app := &Application{
		app:    cli.NewApp(),
		config: &Configure{},
	}
	app.processor = NewProcessor(app)

	app.app.Usage = "observer blockchain transaction"
	app.app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "load configuration from `FILE`",
			Value: "./obs.conf",
		},
		cli.StringFlag{
			Name:        "datadir, d",
			Usage:       "secify data path",
			Value:       "~/.obs",
			Destination: &app.dataDir,
		},
		cli.IntFlag{
			Name:        "height",
			Usage:       "secify scan block height",
			Value:       -1,
			Destination: &app.scanBegin,
		},
	}
	app.app.Before = app.Before
	app.app.Action = app.Action
	return app
}
