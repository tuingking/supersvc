package logger

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type ZeroLogKey int

const (
	ZeroLogSubLogger    ZeroLogKey = iota // type: zerolog.Logger
	ZeroLogSubLoggerCtx                   // type: context.Context
)

func Setup(opt Option, svcname string) {
	// set time format
	switch opt.TimeFormat {
	case "unix":
		//  UNIX Time is faster and smaller than most timestamps
		zerolog.TimeFieldFormat = "unix"
	case "":
		zerolog.TimeFieldFormat = "unix"
	default:
		zerolog.TimeFieldFormat = opt.TimeFormat
	}

	// set output
	if opt.Output == "stdout" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// set log level
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if lv, err := zerolog.ParseLevel(opt.Level); err == nil {
		zerolog.SetGlobalLevel(lv)
	}

	// show caller
	if opt.ShowCaller {
		log.Logger = log.With().Str("svc", svcname).Caller().Logger()
	}

	// get root dir
	_, b, _, _ := runtime.Caller(0)
	rootPath := filepath.Join(filepath.Dir(b), "../..") // remove -> /Users/username/go/src/github.com/tuingking/supersvc

	// customize global caller
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		sort := strings.ReplaceAll(file, rootPath, "")
		return strings.TrimPrefix(sort, "/") + ":" + strconv.Itoa(line)
	}

	log.Debug().Msg("zerolog initialized")
}

func Get(ctx context.Context) zerolog.Logger {
	return ctx.Value(ZeroLogSubLogger).(zerolog.Logger)
}
