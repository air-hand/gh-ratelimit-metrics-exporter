package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStdoutLogger_WriteToStdout(t *testing.T) {
	// Capture stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to pipe: %v", err)
	}

	stdout := os.Stdout
	// Restore stdout after the test
	defer func() {
		os.Stdout = stdout
	}()
	os.Stdout = w

	logger := NewStdoutLogger()
	logger.Info().Msg("foo bar baz")

	w.Close()

	// Read from the pipe
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatal(err)
	}

	s := buf.String()
	assert.Contains(t, s, "foo bar baz")
}
