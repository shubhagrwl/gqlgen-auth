package gracefulshutdown

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"todo/internal/app/service/logger"
)

var (
	shutdownWaitGroup   *sync.WaitGroup
	shutdownChannel     chan int
	initVariablesOnce   sync.Once
	exitChannelInstance chan struct{}
)

// Init : Init
func Init(exitChannel chan struct{}) {
	if shutdownWaitGroup == nil {
		initVariablesOnce.Do(func() {
			shutdownWaitGroup = new(sync.WaitGroup)
			shutdownChannel = make(chan int)
			exitChannelInstance = exitChannel
		})
	}
}

// GetVariables : GetVariables
func GetVariables() (*sync.WaitGroup, chan int) {
	return shutdownWaitGroup, shutdownChannel
}

// Shutdown : Shutdown
func Shutdown(srv *http.Server) {

	// QuitChannel & Signals Map
	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// When Quit Signal received, send shutdown
	<-quitChannel
	logger.SugarLogger.Info("Quit signal received....")

	// HTTP Context Shutdown
	contextTimeoutInSeconds := time.Duration(10)
	// Wait for interrupt signal to gracefully shutdown the server
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeoutInSeconds*time.Second)
	logger.SugarLogger.Info("Quit signal received, sending shutdown and waiting on HTTP calls...")
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.SugarLogger.Fatal("Error Occurred")
	}
	logger.SugarLogger.Info("HTTP Server, shutdown gracefully.")

	// Go Routines Shutdown
	logger.SugarLogger.Info("Quit signal received, sending shutdown and waiting on goroutines...")
	close(shutdownChannel)
	// Go Routines shutdownWaitGroup
	shutdownWaitGroup.Wait()
	logger.SugarLogger.Info("All go routines shutdown gracefully.")

	// Actual shutdown trigger.
	logger.SugarLogger.Info("main goroutine shutdown triggering...")
	close(exitChannelInstance)
}
