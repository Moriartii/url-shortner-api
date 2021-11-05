package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Moriartii/url-shortner-api/app"
	"github.com/Moriartii/url-shortner-api/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	log       *zap.SugaredLogger
	ipAddress string
)

func init() {
	log = logger.GetLogger().Named("application (main.go)").Sugar()
}

//type config struct {
//TODO ipAddress string
//TODO dbName
//TODO dbUser
//TODO can use https://github.com/caarlos0/env
//}

func main() {
	defer log.Sync()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	router := gin.Default()
	ipAddress = os.Getenv("ip_address")
	srv := &http.Server{Addr: ipAddress, Handler: router, ReadTimeout: 5 * time.Second, WriteTimeout: 10 * time.Second}

	go app.StartApplication(srv, router)
	log.Infof("[MAIN-APP] Server start information: %+v", srv)

	new_signal := <-quit
	log.Infof("[MAIN-APP] receive interrupt signal: %+v", new_signal)

	if err := srv.Close(); err != nil {
		log.Errorf("[MAIN-APP] ERROR when close server: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		//extra handling here (for example close db session etc)
		log.Infof("[MAIN-APP] BEGIN CALL: cancel() from context.WithTimeout()")
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Errorf("[MAIN-APP] Server forced to shutdown: %s", err)
	}

	log.Info("[MAIN-APP] SUCCESS CALLED srv.Shutdown()")
	log.Info("[MAIN-APP] Server exited")

}
