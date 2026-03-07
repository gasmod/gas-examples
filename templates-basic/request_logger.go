package main

import (
	"log/slog"

	"github.com/gasmod/gas"
	log "github.com/gasmod/gas-log"
)

type RequestLogger interface {
	gas.Logger
}

func RequestLoggerCtor() func() gas.Logger {
	return func() gas.Logger {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		return log.NewSlogLogger(log.WithSlogInstance(slog.Default()))()
	}
}
