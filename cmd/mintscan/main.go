package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cosmostation/cosmostation-cosmos/app"
	"github.com/cosmostation/cosmostation-cosmos/handler"
	commonhandler "github.com/cosmostation/cosmostation-cosmos/handler/common"
	customhandler "github.com/cosmostation/cosmostation-cosmos/handler/custom"

	"go.uber.org/zap"

	"github.com/gorilla/mux"
)

var (
	Version = "Development"
	Commit  = ""
)

func main() {
	fileBaseName := "mintscan"
	mintscan := app.NewApp(fileBaseName)

	mintscan.SetChainID()
	mintscan.SetMessageInfo()

	r := mux.NewRouter()
	r = r.PathPrefix("/v1").Subrouter()
	commonhandler.RegisterHandlers(mintscan, r)
	customhandler.RegisterHandlers(mintscan, r)
	mintscan.DB.PrepareStmt()

	sm := &http.Server{
		Addr:         ":" + mintscan.Config.Web.Port,
		Handler:      handler.Middleware(r),
		ReadTimeout:  10 * time.Second, // max time to read request from the client
		WriteTimeout: 10 * time.Second, // max time to write response to the client
	}

	go func() {
		for {
			if err := commonhandler.SetStatus(mintscan); err != nil {
				time.Sleep(1 * time.Second)
				continue
			}
			time.Sleep(5 * time.Second)
		}
	}()

	// Start the Mintscan API server.
	go func() {
		zap.S().Infof("Server is running on http://localhost:%s", mintscan.Config.Web.Port)
		zap.S().Infof("Version: %s | Commit: %s", Version, Commit)

		err := sm.ListenAndServe()
		if err != nil {
			zap.S().Fatal(err)
			os.Exit(1)
		}
	}()

	TrapSignal(sm)
}

//TrapSignal trap the signal from os
func TrapSignal(sm *http.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c,
		syscall.SIGINT,  // interrupt
		syscall.SIGKILL, // kill
		syscall.SIGTERM, // terminate
		syscall.SIGHUP)  // hangup(reload)

	for {
		sig := <-c // Block until a signal is received.
		switch sig {
		case syscall.SIGHUP:
			// cfg.ReloadConfig()
		default:
			terminate(sm, sig)
			break
		}
	}
}

func terminate(sm *http.Server, sig os.Signal) {
	// Gracefully shutdown the server, waiting max 30 seconds for current operations to complete.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	sm.Shutdown(ctx)

	zap.S().Infof("Gracefully shutting down the server: %s", sig)
}
