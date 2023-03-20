package logger

type Logger interface {
}

type Option struct {
	Level      string
	TimeFormat string
	Output     string
	ShowCaller bool
}
