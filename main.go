package main

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/renevo/gateway/config"
	"github.com/renevo/gateway/logging"
	"github.com/renevo/gateway/server"
)

func main() {
	cfgFile := flag.String("config", "", "path to configuration file")
	flag.Parse()

	var gatewayConfig *config.Configuration

	if *cfgFile != "" {
		logging.Infof("Loading configuration file %s", *cfgFile)
		f, err := os.Open(*cfgFile)
		if err != nil {
			panic(fmt.Errorf("failed to read configuration file %s: %v", *cfgFile, err))
		}

		cfg, err := config.LoadConfiguration(f)
		f.Close()
		if err != nil {
			panic(fmt.Errorf("failed to read configuration file %s: %v", *cfgFile, err))
		}

		gatewayConfig = cfg
	} else {
		logging.Info("Loading default configuration")
		gatewayConfig = config.DefaultConfiguration()
	}

	// build our server up
	server := server.New(
		server.MountSite(gatewayConfig.Site.Content.Path),
	)

	for _, listener := range gatewayConfig.Site.Listeners {
		listenerAddress, err := listener.Address.URL()
		if err != nil {
			panic(fmt.Errorf("failed to parse listener address %q: %v", listener.Address, err))
		}

		// TODO: handle tls vs non-tls
		go func(addr *url.URL) {
			panic(server.Listen(addr))
		}(listenerAddress)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	server.Shutdown(ctx)
}
