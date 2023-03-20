package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/tuingking/supersvc/entity"
	"github.com/tuingking/supersvc/pkg/ctxkey"
	"github.com/tuingking/supersvc/pkg/parser"
	"github.com/tuingking/supersvc/svc/user"
)

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value(ctxkey.ZeroLogSubLogger).(zerolog.Logger)

	var resp entity.HttpResponse
	defer resp.Render(w, r)

	var p user.GetUserParam
	par := parser.InitParamParser()
	err := par.Decode(&p, r.URL.Query())
	if err != nil {
		logger.Err(err).Msg("failed: parser.Decode param")
		fmt.Fprintf(w, "err: decode param")
		return
	}

	user, pagination, err := h.user.GetUser(r.Context(), p)
	if err != nil {
		logger.Err(err).Msg("failed: user.GetUser")
		fmt.Fprintf(w, "err: get user")
		return
	}

	resp.Data = user
	resp.Pagination = pagination
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value(ctxkey.ZeroLogSubLogger).(zerolog.Logger)

	var resp entity.HttpResponse
	defer resp.Render(w, r)

	var req user.User
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Err(err).Msg("err: decode request body")
		fmt.Fprintf(w, "err: decode request body")
		return
	}

	user, err := h.user.CreateUser(r.Context(), req)
	if err != nil {
		logger.Err(err).Msg("err: create user")
		fmt.Fprintf(w, "err: create user")
		return
	}

	resp.Data = user
	resp.Message = "user created"
}
