package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/vxxvvxxv/hashcash/internal/db"
	"github.com/vxxvvxxv/hashcash/internal/env"
	"github.com/vxxvvxxv/hashcash/internal/logger"
	"github.com/vxxvvxxv/hashcash/internal/server"
	"github.com/vxxvvxxv/hashcash/internal/server/handlers"
)

var (
	logLevel         = env.GetString("LOG_LEVEL", logger.InfoLevel)
	serverAddr       = env.GetString("SERVER_ADDR", "localhost:8080")
	serverConnTTL    = env.GetDuration("SERVER_TTL_CONNECTION", time.Minute)
	headerDifficulty = env.GetInt("HEADER_DIFFICULTY", 5)
	headerTTL        = env.GetDuration("HEADER_TTL", time.Minute*10)
)

func main() {
	loggerInstance := logger.NewLogger(logLevel)

	dbInstance, errDB := db.NewDBService(loggerInstance)
	if errDB != nil {
		loggerInstance.Fatal("can't create db service: ", errDB)
	}
	if errDB = dbInstance.FillTestData(); errDB != nil {
		loggerInstance.Fatal("can't fill test data: ", errDB)
	}

	serverInstance, errServer := server.NewServer(
		serverAddr,
		dbInstance,
		loggerInstance,
		serverConnTTL,
		handlers.NewHandlerManager(dbInstance, loggerInstance, serverAddr, headerTTL, headerDifficulty),
	)
	if errServer != nil {
		loggerInstance.Fatal("can't create server service: " + errServer.Error())
	}

	ctxCancel, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := sync.WaitGroup{}

	// Catch signals and shutdowns
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	// Catch signals and shutdowns
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-c:
			cancel()
		case <-ctxCancel.Done():
		}
	}()

	// Shutdown server
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-ctxCancel.Done():
			loggerInstance.Info("Stopping the server...")

			if err := serverInstance.Shutdown(); err != nil {
				loggerInstance.Error("can't shutdown server: ", err)
			}

			loggerInstance.Info("Server is stopped...")
		}
	}()

	// Listen and serve
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancel()

		loggerInstance.Info("Starting the server on ", serverAddr)

		if errListen := serverInstance.ListenAndServe(ctxCancel); errListen != nil {
			loggerInstance.Error("can't listen and serve: ", errListen)
		}
	}()

	loggerInstance.Info("Waiting for requests...")

	wg.Wait()
	loggerInstance.Info("Finish")

}
