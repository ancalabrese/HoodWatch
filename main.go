package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/ancalabrese/hoodwatch/core"
	"github.com/hashicorp/go-hclog"
)

var (
	debug      = false
	configPath = "./config.yml"
	config     = (*core.Config)(nil)
)

func init() {
	flag.BoolVar(&debug, "d", debug, "Enable debug logging")
	flag.StringVar(&configPath, "config", configPath, "configuration file path")
}

func main() {
	var err error
	flag.Parse()
	l := hclog.Default()
	if debug {
		l.SetLevel(hclog.Debug)
	}

	config, err = core.Load(configPath)
	if err != nil {
		l.Error("Configuration file", "error", err)
		os.Exit(1)
	}

	for _, h := range config.Hoods {
		hq := core.NewHq(l)
		err = h.HoodInit(hq, l)
		if err != nil {
			l.Error("Couldn't initialise hood", "Name", h.Name, "Error", err)
			continue
		}
		h.Watch()
	}

	// Handle interrupt command
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	sig := <-sigChan
	l.Info("Got system signal", "sig", sig)
}
