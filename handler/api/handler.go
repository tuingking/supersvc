package api

import (
	"github.com/rs/zerolog/log"
	"github.com/tuingking/supersvc/config"
	"github.com/tuingking/supersvc/svc/user"
)

type Handler struct {
	cfg *config.Config

	// service
	user user.Service
}

func NewHandler(cfg *config.Config, user user.Service) Handler {
	h := Handler{
		cfg:  cfg,
		user: user,
	}
	log.Debug().Msg("api handler initalized")

	return h
}
