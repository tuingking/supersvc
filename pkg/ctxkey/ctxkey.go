package ctxkey

type CtxKey string

const (
	// http header
	UserAgent     CtxKey = "User-Agent"
	ContentLength CtxKey = "Content-Length"
	ContentType   CtxKey = "Content-Type"

	// x field
	XRequestID CtxKey = "x-request-id"

	// custom
	ZeroLogSubLogger    CtxKey = "x-zerolog"     // type: zerolog.Logger
	ZeroLogSubLoggerCtx CtxKey = "x-zerolog-ctx" // type: context.Context
)

func (k CtxKey) String() string {
	return string(k)
}
