package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"
	"webchat-basedon-pubsub/internal/config"
	"webchat-basedon-pubsub/internal/logger"
	"webchat-basedon-pubsub/internal/server"
)

var log *logger.Logger

func init() {
	log = logger.NewLogger("webchat")
}

func main() {
	log.I("Starting service ...")

	err := Initialize()
	if err != nil {
		log.E("Failed to initialize service: %v", err)
		os.Exit(1)
	}

	err = StartService()
	if err != nil {
		log.E("Failed to start service: %v", err)
		os.Exit(1)
	}
	log.E("Exiting service ...")
}

// Initialize initializes DB and updates DB tables.
func Initialize() error {
	log.I("initiate the service...")

	// Configuration loading
	var configFileName string = "configs/config.json"

	conf := config.GetInstance()
	if !conf.Load(configFileName) {
		log.E("Failed to load config file: %s", configFileName)
		os.Exit(1)
	}
	log.D("Configuration has been loaded.")

	// Setup log level
	logger.SetupLogger(conf.Logging.Enable, conf.Logging.Level)

	// Setup signal handlers for interruption and termination
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for sig := range sigCh {
			if sig == syscall.SIGINT || sig == syscall.SIGTERM {
				log.D("Graceful Termination Time = %d", conf.GracefulTermTimeMillis)
				time.Sleep(time.Duration(conf.GracefulTermTimeMillis) * time.Millisecond)
				Finalize()
				os.Exit(ExitFailure)
			}
		}
	}()

	return nil
}

// StartService starts all the component of this service.
func StartService() error {
	log.I("start the service...")

	conf := config.GetInstance()

	var err error
	if err = server.InitServer(conf); err != nil {
		log.E("Failed to start the chat server: err:[%v]", err)
	}

	return nil
}

// Finalize and clean up the service
func Finalize() {
	//	db.Close()
	//	client.Close()
	log.E("Shutdown service...")
}

//ExitSuccess is exit code 0 and ExitFailure is exit code 1
const (
	ExitSuccess = iota
	ExitFailure
)
