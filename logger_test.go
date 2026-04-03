package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRotatingWriter_Write(t *testing.T) {
	dir := t.TempDir()
	rw, err := newRotatingWriter(dir)
	if err != nil {
		t.Fatalf("failed to create rotating writer: %v", err)
	}
	defer func() { _ = rw.Close() }()

	msg := "test log message\n"
	n, err := rw.Write([]byte(msg))
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	if n != len(msg) {
		t.Errorf("expected %d bytes written, got %d", len(msg), n)
	}

	// Verify file content
	content, err := os.ReadFile(filepath.Join(dir, logFileName))
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}
	if string(content) != msg {
		t.Errorf("expected %q, got %q", msg, string(content))
	}
}

func TestRotatingWriter_Rotation(t *testing.T) {
	dir := t.TempDir()
	rw, err := newRotatingWriter(dir)
	if err != nil {
		t.Fatalf("failed to create rotating writer: %v", err)
	}
	defer func() { _ = rw.Close() }()

	// Write enough data to trigger rotation (just over 1MB)
	chunk := strings.Repeat("x", 1024) + "\n"
	for i := 0; i < 1025; i++ {
		_, err := rw.Write([]byte(chunk))
		if err != nil {
			t.Fatalf("Write failed at iteration %d: %v", i, err)
		}
	}

	// Verify backup file exists
	backupPath := filepath.Join(dir, logFileName+backupSuffix)
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		t.Fatal("backup file should exist after rotation")
	}

	// Verify current log file is small (just the data written after rotation)
	info, err := os.Stat(filepath.Join(dir, logFileName))
	if err != nil {
		t.Fatalf("failed to stat log file: %v", err)
	}
	if info.Size() > maxLogSize {
		t.Errorf("log file should be smaller than maxLogSize after rotation, got %d", info.Size())
	}
}

func TestRotatingWriter_BackupOverwrite(t *testing.T) {
	dir := t.TempDir()

	// Create a pre-existing backup file
	backupPath := filepath.Join(dir, logFileName+backupSuffix)
	if err := os.WriteFile(backupPath, []byte("old backup"), 0644); err != nil {
		t.Fatalf("failed to create backup file: %v", err)
	}

	rw, err := newRotatingWriter(dir)
	if err != nil {
		t.Fatalf("failed to create rotating writer: %v", err)
	}
	defer func() { _ = rw.Close() }()

	// Trigger rotation
	chunk := strings.Repeat("x", 1024) + "\n"
	for i := 0; i < 1025; i++ {
		_, err := rw.Write([]byte(chunk))
		if err != nil {
			t.Fatalf("Write failed at iteration %d: %v", i, err)
		}
	}

	// Verify backup was overwritten (not "old backup")
	content, err := os.ReadFile(backupPath)
	if err != nil {
		t.Fatalf("failed to read backup: %v", err)
	}
	if string(content) == "old backup" {
		t.Error("backup file should have been overwritten")
	}
}

func TestSetupLogger(t *testing.T) {
	logger, cleanup := setupLogger()
	defer cleanup()

	if logger == nil {
		t.Fatal("logger should not be nil")
	}

	// Just verify it doesn't panic
	logger.Println("test message")
}
