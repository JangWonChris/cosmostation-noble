package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cosmostation/cosmostation-cosmos/mintscan/client"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/db"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/handler"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/handler/common"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/handler/custom"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/handler/mobile"
	cfg "github.com/cosmostation/mintscan-backend-library/config"

	"go.uber.org/zap"

	"github.com/gorilla/mux"
)

var (
	// Version is a project's version string.
	Version = "Development"

	// Commit is commit hash of this project.
	Commit = ""
)

func init() {
	l, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(l)
	defer l.Sync()
}

func main() {
	config := cfg.ParseConfig()

	client := client.NewClient(&config.Client)

	db := db.Connect(&config.DB)
	err := db.Ping()
	if err != nil {
		zap.L().Error("failed to ping database", zap.Error(err))
		return
	}

	s := handler.SetSession(client, db)

	r := mux.NewRouter()
	r = r.PathPrefix("/v1").Subrouter()
	common.RegisterHandlers(s, r)
	custom.RegisterHandlers(s, r)
	mobile.RegisterHandlers(s, r)

	// r.HandleFunc("/tx/{hash}", common.GetLegacyTransactionFromDB).Methods("GET")                // /tx?hash={hash}

	sm := &http.Server{
		Addr:         ":" + config.Web.Port,
		Handler:      handler.Middleware(r),
		ReadTimeout:  10 * time.Second, // max time to read request from the client
		WriteTimeout: 10 * time.Second, // max time to write response to the client
	}

	go func() {
		for {
			if err := common.SetStatus(); err != nil {
				time.Sleep(1 * time.Second)
				continue
			}
			time.Sleep(5 * time.Second)
		}
	}()

	// Start the Mintscan API server.
	go func() {
		zap.S().Infof("Server is running on http://localhost:%s", config.Web.Port)
		zap.S().Infof("Version: %s | Commit: %s", Version, Commit)

		err := sm.ListenAndServe()
		if err != nil {
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
