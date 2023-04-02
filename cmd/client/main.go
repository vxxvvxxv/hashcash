package main

import (
	"context"
	"errors"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/vxxvvxxv/hashcash/internal/client"
	"github.com/vxxvvxxv/hashcash/internal/env"
	"github.com/vxxvvxxv/hashcash/internal/logger"
)

var (
	logLevel            = env.GetString("LOG_LEVEL", logger.InfoLevel)
	serverAddr          = env.GetString("SERVER_ADDR", "localhost:8080")
	clientTimeout       = env.GetDuration("CLIENT_TIMEOUT", time.Second*5)
	headerMaxIterations = env.GetInt("HEADER_MAX_ITERATIONS", -1)
	headerTTL           = env.GetDuration("HEADER_TTL", time.Minute*10)
	isDDOS              = env.GetBool("CLIENT_DDOS_MODE", false)
	ddosCountClient     = env.GetInt("CLIENT_DDOS_CLIENT_COUNT", 100)
	ddosTimeout         = env.GetDuration("CLIENT_DDOS_TIMEOUT", time.Minute)
	ddosWaitIfError     = env.GetDuration("CLIENT_DDOS_WAITING", time.Second)
)

func main() {
	loggerInstance := logger.NewLogger(logLevel)

	ctxCancel, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := sync.WaitGroup{}

	// Catch signals and shutdowns
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-c:
			cancel()
		case <-ctxCancel.Done():
		}
	}()

	if isDDOS {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer cancel()

			ddosClients := make([]client.Client, 0, ddosCountClient)

			for i := 0; i < ddosCountClient; i++ {
				clientInstance, errClient := client.NewClient(
					serverAddr,
					clientTimeout,
					loggerInstance,
					headerMaxIterations,
					headerTTL,
				)
				if errClient != nil {
					loggerInstance.Error("can't create client: ", errClient)
					continue
				}

				ddosClients = append(ddosClients, clientInstance)
			}

			if len(ddosClients) == 0 {
				loggerInstance.Fatal("can't create clients:")
			}

			ctxTimeout, cancelTimeout := context.WithDeadline(ctxCancel, time.Now().Add(ddosTimeout))
			defer cancelTimeout()

			loggerInstance.Info("CLIENT: Start DDOS mode")
			loggerInstance.Info("CLIENT: Client count: ", len(ddosClients))
			loggerInstance.Info("CLIENT: Timeout: ", ddosTimeout)
			loggerInstance.Info("CLIENT: -----------------------------")

			wgClients := sync.WaitGroup{}

			for _, clientI := range ddosClients {
				wgClients.Add(1)
				go func(clientInstance client.Client) {
					defer wgClients.Done()
					for {
						select {
						case <-ctxCancel.Done():
							loggerInstance.Info("CLIENT: Context is closed")
							return
						case <-ctxTimeout.Done():
							loggerInstance.Info("CLIENT: Timeout is reached")
							return
						default:
							if err := clientInstance.Connect(ctxCancel); err != nil {
								time.Sleep(ddosWaitIfError)
							}
							if err := clientInstance.Start(ctxCancel); err != nil {
								time.Sleep(ddosWaitIfError)
							}
						}
					}
				}(clientI)
			}

			wgClients.Wait()
		}()
	} else {
		clientInstance, errClient := client.NewClient(
			serverAddr,
			clientTimeout,
			loggerInstance,
			headerMaxIterations,
			headerTTL,
		)
		if errClient != nil {
			loggerInstance.Fatal("can't create client: ", errClient)
		}

		// Shutdown server
		wg.Add(1)
		go func() {
			defer wg.Done()
			select {
			case <-ctxCancel.Done():
				if err := clientInstance.Stop(); err != nil {
					if !errors.Is(err, net.ErrClosed) {
						loggerInstance.Error("can't stop client: ", err)
					}
				}
			}
		}()

		// Listen and serve
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer cancel()

			if err := clientInstance.Connect(ctxCancel); err != nil {
				if !errors.Is(err, io.EOF) && !errors.Is(err, net.ErrClosed) && !errors.Is(err, io.ErrClosedPipe) {
					loggerInstance.Error("can't connect to server: ", err)
				}
			}

			if err := clientInstance.Start(ctxCancel); err != nil {
				if !errors.Is(err, io.EOF) && !errors.Is(err, net.ErrClosed) && !errors.Is(err, io.ErrClosedPipe) {
					loggerInstance.Error("can't start to server: ", err)
				}
			}
		}()
	}

	loggerInstance.Info("Client is connected to ", serverAddr)

	wg.Wait()
	loggerInstance.Info("Client is stopped...")
}
