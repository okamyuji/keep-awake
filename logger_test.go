package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSetupLogger(t *testing.T) {
	dir := t.TempDir()

	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(orig); err != nil {
			t.Errorf("failed to restore cwd: %v", err)
		}
	})

	logger, cleanup := setupLogger()
	defer cleanup()

	if logger == nil {
		t.Fatal("logger should not be nil")
	}
	logger.Println("test message")
}

func TestSetupLogger_TruncatesOnStartup(t *testing.T) {
	dir := t.TempDir()

	// Create a pre-existing log file with old content
	path := filepath.Join(dir, logFileName)
	if err := os.WriteFile(path, []byte("old log data\n"), 0644); err != nil {
		t.Fatalf("failed to create log file: %v", err)
	}

	// Change to temp dir so setupLogger creates log there
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(orig); err != nil {
			t.Errorf("failed to restore cwd: %v", err)
		}
	})

	logger, cleanup := setupLogger()
	defer cleanup()

	logger.Println("new session")

	// Read the log file and verify old content is gone
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}
	if string(content) == "" {
		t.Fatal("log file should not be empty")
	}
	if strings.Contains(string(content), "old log data") {
		t.Error("old log data should have been truncated")
	}
	if !strings.Contains(string(content), "new session") {
		t.Error("new session log should be present")
	}
}
