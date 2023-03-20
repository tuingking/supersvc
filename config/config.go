package config

import (
	"log"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/tuingking/supersvc/pkg/httpserver"
	"github.com/tuingking/supersvc/pkg/logger"
	"github.com/tuingking/supersvc/pkg/mysql"
	"github.com/tuingking/supersvc/svc/user"
)

type Config struct {
	HttpServer httpserver.Option
	MySQL      map[string]mysql.Option
	Logger     logger.Option

	// service config
	User user.Option
}

func InitConfig(opts ...Option) *Config {
	var cfg Config

	// default option
	opt := DefaultOption()

	// override option
	for _, fn := range opts {
		fn(opt)
	}

	v := viper.New()
	v.AddConfigPath(opt.configPath)
	v.SetConfigName(opt.configName)
	v.SetConfigType(opt.configType)

	if err := v.ReadInConfig(); err != nil {
		log.Fatal(errors.Wrap(err, "read in config"))
	}
	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatal(errors.Wrap(err, "unmarshal config"))
	}

	log.Printf("config file: %s", opt.configPath+"/"+opt.configName+"."+opt.configType)

	return &cfg
}

type option struct {
	configPath string
	configName string
	configType string
}

func DefaultOption() *option {
	return &option{
		configPath: "./config/",
		configName: "config",
		configType: "yaml",
	}
}

type Option func(*option)

func WithConfigName(configName string) Option {
	return func(o *option) {
		o.configName = configName
	}
}

func WithConfigPath(configPath string) Option {
	return func(o *option) {
		o.configPath = configPath
	}
}

func WithConfigType(configType string) Option {
	return func(o *option) {
		o.configType = configType
	}
}
