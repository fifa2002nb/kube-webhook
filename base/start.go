package base

import (
	"crypto/tls"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"kubewebhook/config"
	"kubewebhook/utils/webhook"
	"net/http"
	"os"
	"os/signal"
)

func Start(c *cli.Context) {
	var (
		err     error
		options *config.Options
	)
	options, err = config.ParseConf(c)
	if nil != err {
		log.Fatal(err.Error())
		os.Exit(1)
	}
	sidecarConfig, err := webhook.LoadConfig(options.SidecarCfgFile)
	if nil != err {
		log.Fatal(err.Error())
		os.Exit(1)
	}
	pair, err := tls.LoadX509KeyPair(options.TlsCertFile, options.TlsKeyFile)
	if nil != err {
		log.Fatal(err.Error())
		os.Exit(1)
	}
	whsvr := &webhook.WebhookServer{
		SidecarConfig: sidecarConfig,
		Server: &http.Server{
			Addr:      fmt.Sprintf(":%v", options.Port),
			TLSConfig: &tls.Config{Certificates: []tls.Certificate{pair}},
		},
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/mutate", whsvr.Serve)
	whsvr.Server.Handler = mux

	go func() {
		if err := whsvr.Server.ListenAndServeTLS("", ""); err != nil {
			log.Errorf("Filed to listen and serve webhook server: %v", err)
		}
	}()

	waitingForExit()
}

func waitingForExit() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	killing := false
	for range sc {
		if killing {
			log.Info("Second interrupt: exiting")
			os.Exit(1)
		}
		killing = true
		go func() {
			log.Info("Interrupt: closing down...")
			log.Info("done")
			os.Exit(1)
		}()
	}
}
