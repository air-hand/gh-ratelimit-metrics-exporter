package main

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

func NewNullLogger() *zerolog.Logger {
	return newLogger(io.Discard)
}

func NewStdoutLogger() *zerolog.Logger {
	return newLogger(os.Stdout)
}

func newLogger(w io.Writer) *zerolog.Logger {
	logger := zerolog.New(w)
	return &logger
}
