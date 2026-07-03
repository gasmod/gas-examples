package main

import (
	"github.com/gasmod/gas"
)

// RequestLogger is a distinct type so the DI container can register a scoped
// logger separately from the singleton gas.Logger. The request logger
// middleware (gas.RequestLogger) mutates the logger in-place via
// SetBaseFields().Apply() to stamp per-request fields (request ID, method,
// path). If the singleton gas.Logger were used directly, those mutations
// would corrupt the shared instance. A separate scoped type avoids this.
type RequestLogger interface {
	gas.Logger
}

// requestLogger is the DI constructor for RequestLogger. It receives the
// singleton gas.Logger and returns a clone via With().Logger(). The clone
// shares the same output/handler but is an independent instance safe to
// mutate per-request.
func requestLogger(l gas.Logger) RequestLogger {
	return l.With().Logger()
}
