package svc

import (
	"Windy-API/config"
	zl "github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var version string
var build string

func Init(cfg config.Config) {
	version = cfg.Version
	build = cfg.Build
}

func Logger() *zl.Logger {
	l := log.Logger.With().Str("service", ServiceName).Str("version", version).Str("build", build).Logger()
	return &l
}
