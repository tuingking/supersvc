package httpserver

import (
	"context"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/rs/zerolog/log"
)

type HttpServer interface {
	Start() error
	Stop(ctx context.Context) error
}

type httpServer struct {
	opt    Option
	server *http.Server
}

type Option struct {
	Port              string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ReadHeaderTimeout time.Duration
}

func NewHttpServer(opt Option, h http.Handler) HttpServer {
	server := &http.Server{
		ReadTimeout:       opt.ReadTimeout,
		WriteTimeout:      opt.WriteTimeout,
		IdleTimeout:       opt.IdleTimeout,
		ReadHeaderTimeout: opt.ReadHeaderTimeout,
		Handler:           h,
	}

	return &httpServer{
		opt:    opt,
		server: server,
	}
}

func (s *httpServer) Start() error {
	// add sysinfo
	buildInfo, _ := debug.ReadBuildInfo()

	l, err := net.Listen("tcp", s.opt.Port)
	if err != nil {
		log.Fatal().Msgf("failed start http server. err=%s", err)
		return err
	}
	log.Info().Int("pid", os.Getpid()).Str("go_version", buildInfo.GoVersion).Msgf("server running on port %s", s.opt.Port)
	return s.server.Serve(l)
}

func (s *httpServer) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
