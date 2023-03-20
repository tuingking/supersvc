package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/log"

	"github.com/tuingking/supersvc/config"
	"github.com/tuingking/supersvc/handler/api"
	"github.com/tuingking/supersvc/handler/mux"
	"github.com/tuingking/supersvc/pkg/httpserver"
	"github.com/tuingking/supersvc/pkg/logger"
	"github.com/tuingking/supersvc/pkg/mysql"
	"github.com/tuingking/supersvc/svc/user"
)

var (
	ServiceName string
)

func main() {
	// config
	cfg := config.InitConfig()
	logger.Setup(cfg.Logger, ServiceName)

	// infra
	dbLocal := cfg.MySQL["localhost"]
	db := mysql.NewMySQL(dbLocal)

	// service
	userRepo := user.NewRepository(cfg.User.Repository, db)
	usersvc := user.NewService(cfg.User, userRepo)

	// handler
	apiHandler := api.NewHandler(cfg, usersvc)
	httpHandler := mux.NewMux(apiHandler)

	// server
	httpServer := httpserver.NewHttpServer(cfg.HttpServer, httpHandler)
	run(httpServer)
}

func run(server httpserver.HttpServer) {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGHUP,
		syscall.SIGQUIT,
		syscall.SIGINT,
	)

	go server.Start()

	s := <-signalChannel // graceful shutdown
	log.Info().Any("signal", s).Msg("receiving terminate signal")
	signal.Stop(signalChannel)
	close(signalChannel)

	if err := server.Stop(context.Background()); err != nil {
		log.Error().Msg("failed to stop http server")
		return
	}
	log.Info().Msg("app stopped")
}
