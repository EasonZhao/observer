package main

import (
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/larspensjo/config"
	"github.com/urfave/cli"
	"strconv"
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
)

// Configure config
type Configure struct {
	ExtendKey     string
	Start         string
	Net           *chaincfg.Params
	UseSigWitness bool
	Count         int // 地址数量
	Bind          string
}

// Application application
type Application struct {
	app       *cli.App
	config    *Configure
	addresses map[string]string
	lastErr   error
	processor *Processor
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
	for i, v := range addrs {
		s.addresses[v] = strconv.Itoa(i)
	}
	return nil
}

// Action action
func (s *Application) Action(c *cli.Context) error {
	if s.lastErr != nil {
		fmt.Println(s.lastErr)
	} else {
		//start grpc
		ch := make(chan error, 1)
		if err := s.processor.AsyncGRPC(&ch); err != nil {
			fmt.Println(err)
		} else {
			err := <-ch
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("application exit")
			}
		}
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
	return nil
}

// Stop stop application
func (s *Application) Stop() {
	s.processor.StopGRPC()
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
			Usage: "Load configuration from `FILE`",
			Value: "./obs.conf",
		},
	}
	app.app.Before = app.Before
	app.app.Action = app.Action
	return app
}
