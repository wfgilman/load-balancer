package utils

import (
	"net/http"
)

const (
	Attempts int = iota
	Retries
)

func CountAttempts(req *http.Request) int {
	if attempts, ok := req.Context().Value(Attempts).(int); ok {
		return attempts
	}
	return 1
}

func CountRetries(req *http.Request) int {
	if retries, ok := req.Context().Value(Retries).(int); ok {
		return retries
	}
	return 0
}
