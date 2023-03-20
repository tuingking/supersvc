package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/tuingking/supersvc/pkg/ctxkey"
)

func (h *Handler) HandleRoot(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value(ctxkey.ZeroLogSubLogger).(zerolog.Logger)
	logger.Info().Msg("HandleRoot called!")
	fmt.Fprintf(w, "hello world")
}

func (h *Handler) GetAppConfig(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value(ctxkey.ZeroLogSubLogger).(zerolog.Logger)
	logger.Info().Msg("GetAppConfig called!")
	json.NewEncoder(w).Encode(h.cfg)
}
