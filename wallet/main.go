package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/cosmostation/cosmostation-cosmos/wallet/client"
	"github.com/cosmostation/cosmostation-cosmos/wallet/config"
	"github.com/cosmostation/cosmostation-cosmos/wallet/db"

	"go.uber.org/zap"

	"github.com/gorilla/mux"
)

var (
	// Version is a project's version string.
	Version = "Development"

	// Commit is commit hash of this project.
	Commit = ""
)

func main() {
	l, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(l)

	config := config.NewConfig()

	// Parse config from configuration file (config.yaml).
	config := config.ParseConfig()

	// Create new client with node configruation.
	// Client is used for requesting any type of network data from RPC full node and REST Server.
	client, err := client.NewClient(config.Node, config.Market)
	if err != nil {
		zap.L().Error("failed to create new client", zap.Error(err))
		return
	}

	// Create connection with PostgreSQL database and
	// Ping database to verify connection is success.
	db := db.Connect(config.DB)
	err = db.Ping()
	if err != nil {
		zap.L().Error("failed to ping database", zap.Error(err))
		return
	}

	r := mux.NewRouter()
	r = r.PathPrefix("/v1").Subrouter()
	r.HandleFunc("/account/update", handler.RegisterOrUpdate).Methods("POST")
	r.HandleFunc("/account/delete", handler.Delete).Methods("DELETE")
	r.HandleFunc("/app/version/{deviceType}", handler.GetVersion).Methods("GET")
	r.HandleFunc("/app/version", handler.SetVersion).Methods("POST")

	// These APIs are not used at this moment.
	r.HandleFunc("/alarm/push", handler.PushNotification).Methods("POST")

	sm := &http.Server{
		Addr:         ":" + config.Web.Port,
		Handler:      handler.Middleware(r, client, db),
		ReadTimeout:  10 * time.Second, // max time to read request from the client
		WriteTimeout: 10 * time.Second, // max time to write response to the client
	}

	// Start the Mintscan API server.
	go func() {
		zap.S().Infof("Server is running on http://localhost:%s", config.Web.Port)
		zap.S().Infof("Network Type: %s | Version: %s | Commit: %s", config.Node.NetworkType, Version, Commit)

		err := sm.ListenAndServe()
		if err != nil {
			os.Exit(1)
		}
	}()

	TrapSignal(sm)
}

// TrapSignal traps sigterm or interupt and gracefully shutdown the server.
func TrapSignal(sm *http.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c

	// Gracefully shutdown the server, waiting max 30 seconds for current operations to complete.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	sm.Shutdown(ctx)

	zap.S().Infof("Gracefully shutting down the server: %s", sig)
}
