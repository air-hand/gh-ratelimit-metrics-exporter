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
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal("Failed to generate rsa private key.")
	}

	pemData := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	pemDataString := string(pemData)

	tests := []struct {
		testName       string
		appID          string
		installationID string
		privateKey     string
		expectNil      bool
	}{
		{
			testName:       "No env at all",
			appID:          "",
			installationID: "",
			privateKey:     "",
			expectNil:      true,
		},
		{
			testName:       "Lack of app_id",
			appID:          "",
			installationID: "1008",
			privateKey:     pemDataString,
			expectNil:      true,
		},
		{
			testName:       "Not integer app_id",
			appID:          "abc",
			installationID: "1008",
			privateKey:     pemDataString,
			expectNil:      true,
		},
		{
			testName:       "Lack of installation_id",
			appID:          "100",
			installationID: "",
			privateKey:     pemDataString,
			expectNil:      true,
		},
		{
			testName:       "Not integer installation_id",
			appID:          "100",
			installationID: "def",
			privateKey:     pemDataString,
			expectNil:      true,
		},
		{
			testName:       "Lack of private_key",
			appID:          "100",
			installationID: "1008",
			privateKey:     "",
			expectNil:      true,
		},
		{
			testName:       "Broken private key",
			appID:          "100",
			installationID: "1008",
			privateKey:     "foobarbaz",
			expectNil:      true,
		},
		{
			testName:       "Ok",
			appID:          "100",
			installationID: "1008",
			privateKey:     pemDataString,
			expectNil:      false,
		},
	}

	logger := NewNullLogger()

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			t.Setenv("GH_APP_ID", tt.appID)
			t.Setenv("GH_INSTALLATION_ID", tt.installationID)
			t.Setenv("GH_PRIVATE_KEY", tt.privateKey)

			if tt.expectNil {
				assert.Nil(t, newClientWithGitHubApp(logger))
			} else {
				assert.NotNil(t, newClientWithGitHubApp(logger))
			}
		})
	}
}
