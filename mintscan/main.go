package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cosmostation/cosmostation-cosmos/chain-config/custom"
	pg "github.com/go-pg/pg/v10"

	"github.com/cosmostation/cosmostation-cosmos/mintscan/client"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/db"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/handler"
	commonhandler "github.com/cosmostation/cosmostation-cosmos/mintscan/handler/common"
	customhandler "github.com/cosmostation/cosmostation-cosmos/mintscan/handler/custom"
	mobilehandler "github.com/cosmostation/cosmostation-cosmos/mintscan/handler/mobile"
	cfg "github.com/cosmostation/mintscan-backend-library/config"

	"go.uber.org/zap"

	"github.com/gorilla/mux"
)

var (
	Version = "Development"
	Commit  = ""
)

type dbLogger struct{}

func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	query, err := q.FormattedQuery()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(query))
	return nil
}

func init() {
	if !custom.IsSetAppConfig() {
		panic(fmt.Errorf("appconfig was not set"))
	}

	l, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(l)
	defer l.Sync()
}

func main() {
	fileBaseName := "mintscan"
	config := cfg.ParseConfig(fileBaseName)

	client := client.NewClient(&config.Client)

	db := db.Connect(&config.DB)
	err := db.Ping()
	if err != nil {
		zap.L().Error("failed to ping database", zap.Error(err))
		return
	}

	db.AddQueryHook(dbLogger{}) // debugging 용

	s := handler.SetSession(client, db)
	handler.SetChainID()
	handler.SetMessageInfo()

	r := mux.NewRouter()
	r = r.PathPrefix("/v1").Subrouter()
	commonhandler.RegisterHandlers(s, r)
	customhandler.RegisterHandlers(s, r)
	mobilehandler.RegisterHandlers(s, r)
	db.PrepareStmt()

	// r.HandleFunc("/tx/{hash}", common.GetLegacyTransactionFromDB).Methods("GET")                // /tx?hash={hash}

	sm := &http.Server{
		Addr:         ":" + config.Web.Port,
		Handler:      handler.Middleware(r),
		ReadTimeout:  10 * time.Second, // max time to read request from the client
		WriteTimeout: 10 * time.Second, // max time to write response to the client
	}

	go func() {
		for {
			if err := commonhandler.SetStatus(); err != nil {
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
