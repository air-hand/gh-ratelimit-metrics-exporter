package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClientWithGitHubApp(t *testing.T) {
	private_key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal("Failed to generate rsa private key.")
	}

	pem_data := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(private_key),
	})
	pem_data_s := string(pem_data)

	tests := []struct {
		test_name       string
		app_id          string
		installation_id string
		private_key     string
		expectNil       bool
	}{
		{
			test_name:       "No env all",
			app_id:          "",
			installation_id: "",
			private_key:     "",
			expectNil:       true,
		},
		{
			test_name:       "Lack of app_id",
			app_id:          "",
			installation_id: "1008",
			private_key:     pem_data_s,
			expectNil:       true,
		},
		{
			test_name:       "Not integer app_id",
			app_id:          "abc",
			installation_id: "1008",
			private_key:     pem_data_s,
			expectNil:       true,
		},
		{
			test_name:       "Lack of installation_id",
			app_id:          "100",
			installation_id: "",
			private_key:     pem_data_s,
			expectNil:       true,
		},
		{
			test_name:       "Not integer installation_id",
			app_id:          "100",
			installation_id: "def",
			private_key:     pem_data_s,
			expectNil:       true,
		},
		{
			test_name:       "Lack of private_key",
			app_id:          "100",
			installation_id: "1008",
			private_key:     "",
			expectNil:       true,
		},
		{
			test_name:       "Broken private key",
			app_id:          "100",
			installation_id: "1008",
			private_key:     "foobarbaz",
			expectNil:       true,
		},
		{
			test_name:       "Ok",
			app_id:          "100",
			installation_id: "1008",
			private_key:     pem_data_s,
			expectNil:       false,
		},
	}

	logger := NewNullLogger()

	for _, tt := range tests {
		t.Run(tt.test_name, func(t *testing.T) {
			t.Setenv("GH_APP_ID", tt.app_id)
			t.Setenv("GH_INSTALLATION_ID", tt.installation_id)
			t.Setenv("GH_PRIVATE_KEY", tt.private_key)

			if tt.expectNil {
				assert.Nil(t, newClientWithGitHubApp(logger))
			} else {
				assert.NotNil(t, newClientWithGitHubApp(logger))
			}
		})
	}
}
