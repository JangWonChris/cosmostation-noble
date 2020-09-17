package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/cosmostation/cosmostation-cosmos/wallet/config"
	"github.com/cosmostation/cosmostation-cosmos/wallet/db"
	"github.com/cosmostation/cosmostation-cosmos/wallet/handler"

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

	// Parse config from configuration file (config.yaml).
	config := config.ParseConfig()

	// Create connection with PostgreSQL database and
	// Ping database to verify connection is success.
	db := db.Connect(config.DB)
	err := db.Ping()
	if err != nil {
		zap.L().Error("failed to ping database", zap.Error(err))
		return
	}

	r := mux.NewRouter()
	r = r.PathPrefix("/v1").Subrouter()
	r.HandleFunc("/account/update", handler.RegisterOrUpdateAccount).Methods("POST")
	r.HandleFunc("/app/version/{device_type}", handler.GetAppVersion).Methods("GET")
	r.HandleFunc("/app/version", handler.SetAppVersion).Methods("POST")
	r.HandleFunc("/sign/moonpay", handler.SignSignature).Methods("POST")

	// NOT USED APIs
	r.HandleFunc("/alarm/push", handler.PushNotification).Methods("POST")
	r.HandleFunc("/account/delete", handler.DeleteAccount).Methods("DELETE")

	sm := &http.Server{
		Addr:         ":" + config.Web.Port,
		Handler:      handler.Middleware(r, config, db),
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
