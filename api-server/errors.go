package main

import (
	"net/http"

	"github.com/gasmod/gas"
)

// statusError is checked by the error handler via interface assertion.
// Any error type with a StatusCode() method controls the HTTP response.
type statusError interface {
	error
	StatusCode() int
}

// errorHandler converts handler errors into JSON responses. If the error
// implements statusError, its status and message are used. Otherwise, a
// generic 500 is returned.
func errorHandler(ctx gas.Context, err error) {
	status := http.StatusInternalServerError
	message := "internal server error"

	if se, ok := err.(statusError); ok {
		status = se.StatusCode()
		message = se.Error()
	}

	ctx.JSON(status, map[string]string{"error": message})
}
