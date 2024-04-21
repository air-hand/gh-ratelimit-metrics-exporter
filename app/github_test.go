package main

import (
	"testing"

	"github.com/google/go-github/v61/github"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestNewClientWithEnv_NotNil(t *testing.T) {
	tests := []struct {
		test_name string
		funcs     []newClientFunc
	}{
		{
			test_name: "empty funcs",
			funcs:     []newClientFunc{},
		},
		{
			test_name: "a func returns nil",
			funcs: []newClientFunc{
				func(*zerolog.Logger) *github.Client {
					return nil
				},
			},
		},
		{
			test_name: "all funcs return nil",
			funcs: []newClientFunc{
				func(*zerolog.Logger) *github.Client {
					return nil
				},
				func(*zerolog.Logger) *github.Client {
					return nil
				},
			},
		},
	}

	logger := NewNullLogger()

	for _, tt := range tests {
		t.Run(tt.test_name, func(t *testing.T) {
			assert.NotNilf(t, newClientWithEnv(tt.funcs, logger), "test_name: %s", tt.test_name)
		})
	}
}
