package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClientWithToken(t *testing.T) {
	tests := []struct {
		test_name        string
		gh_token_env     string
		github_token_env string
		expectNil        bool
	}{
		{
			test_name:        "No env both",
			gh_token_env:     "",
			github_token_env: "",
			expectNil:        true,
		},
		{
			test_name:        "GH_TOKEN only",
			gh_token_env:     "foo",
			github_token_env: "",
			expectNil:        false,
		},
		{
			test_name:        "GITHUB_TOKEN only",
			gh_token_env:     "",
			github_token_env: "bar",
			expectNil:        false,
		},
		{
			test_name:        "Both envs",
			gh_token_env:     "foo",
			github_token_env: "bar",
			expectNil:        false,
		},
	}

	logger := NewNullLogger()

	for _, tt := range tests {
		t.Run(tt.test_name, func(t *testing.T) {
			t.Setenv("GH_TOKEN", tt.gh_token_env)
			t.Setenv("GITHUB_TOKEN", tt.github_token_env)

			if tt.expectNil {
				assert.Nil(t, newClientWithToken(logger))
			} else {
				assert.NotNil(t, newClientWithToken(logger))
			}
		})
	}
}
