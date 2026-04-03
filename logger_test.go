package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSetupLogger(t *testing.T) {
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
	orig, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer func() { _ = os.Chdir(orig) }()

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
	if contains(string(content), "old log data") {
		t.Error("old log data should have been truncated")
	}
	if !contains(string(content), "new session") {
		t.Error("new session log should be present")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
