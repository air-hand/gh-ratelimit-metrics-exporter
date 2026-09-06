package main

import (
	"testing"

	"github.com/google/go-github/v91/github"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestNewClientFromEnv_Nil(t *testing.T) {
	tests := []struct {
		testName string
		funcs    []newClientOptsFunc
	}{
		{
			testName: "empty funcs",
			funcs:    []newClientOptsFunc{},
		},
		{
			testName: "a func returns nil",
			funcs: []newClientOptsFunc{
				func(*zerolog.Logger) ([]github.ClientOptionsFunc, error) {
					return nil, nil
				},
			},
		},
		{
			testName: "all funcs return nil",
			funcs: []newClientOptsFunc{
				func(*zerolog.Logger) ([]github.ClientOptionsFunc, error) {
					return nil, nil
				},
				func(*zerolog.Logger) ([]github.ClientOptionsFunc, error) {
					return nil, nil
				},
			},
		},
	}

	logger := NewNullLogger()

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			c, err := newClientFromEnv(tt.funcs, logger)
			assert.Nilf(t, c, "test_name: %s", tt.testName)
			assert.Error(t, err)
		})
	}
}
