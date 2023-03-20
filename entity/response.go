package entity

import (
	"net/http"
	"time"

	"github.com/go-chi/render"
)

type HttpResponse struct {
	Code       int         `json:"code"`
	Data       interface{} `json:"data,omitempty"`
	Error      *Error      `json:"error,omitempty"`
	Message    string      `json:"message"`
	ServerTime int64       `json:"serverTime"`
	Pagination interface{} `json:"pagination,omitempty"`
}

type Error struct {
	Status bool   `json:"status" example:"false"` // true if we have error
	Msg    string `json:"msg" example:" "`        // error message
	Code   int    `json:"code" example:"0"`       // application error code for tracing
}

// Render writes the http response to the client
func (res *HttpResponse) Render(w http.ResponseWriter, r *http.Request) {
	if res.Code == 0 {
		res.Code = http.StatusOK
	}

	if res.ServerTime == 0 {
		res.ServerTime = time.Now().Unix()
	}

	render.Status(r, res.Code)
	render.JSON(w, r, res)
}

// SetError set the response to return the given error.
// code is http status code, http.StatusInternalServerError is the default value
func (res *HttpResponse) SetError(err error, code ...int) {
	if len(code) > 0 {
		res.Code = code[0]
	} else {
		res.Code = http.StatusInternalServerError
	}

	if err != nil {
		res.Error = &Error{
			Msg:    err.Error(),
			Status: true,
		}
	}

}
