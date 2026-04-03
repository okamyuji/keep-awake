package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

const (
	maxLogSize   = 1 * 1024 * 1024 // 1MB
	logFileName  = "keep-awake.log"
	backupSuffix = ".1"
)

type rotatingWriter struct {
	mu       sync.Mutex
	file     *os.File
	filePath string
	size     int64
}

func newRotatingWriter(dir string) (*rotatingWriter, error) {
	path := filepath.Join(dir, logFileName)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("ログファイルを開けません: %w", err)
	}
	info, err := f.Stat()
	if err != nil {
		_ = f.Close()
		return nil, err
	}
	return &rotatingWriter{
		file:     f,
		filePath: path,
		size:     info.Size(),
	}, nil
}

func (w *rotatingWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.size+int64(len(p)) > maxLogSize {
		if err := w.rotate(); err != nil {
			return 0, err
		}
	}
	n, err = w.file.Write(p)
	w.size += int64(n)
	return n, err
}

func (w *rotatingWriter) rotate() error {
	_ = w.file.Close()
	backup := w.filePath + backupSuffix
	_ = os.Remove(backup)
	_ = os.Rename(w.filePath, backup)
	f, err := os.OpenFile(w.filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	w.file = f
	w.size = 0
	return nil
}

func (w *rotatingWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.file != nil {
		return w.file.Close()
	}
	return nil
}

// setupLogger creates a logger that writes to both stdout and a rotating log file.
// The log file is created in the same directory as the executable.
// Returns the logger and a cleanup function.
func setupLogger() (*log.Logger, func()) {
	// Use executable directory for log file
	dir, err := os.Getwd()
	if err != nil {
		dir = "."
	}

	rw, err := newRotatingWriter(dir)
	if err != nil {
		// Fall back to stdout-only
		fmt.Fprintf(os.Stderr, "警告: ログファイルの作成に失敗: %v (stdout のみに出力します)\n", err)
		return log.New(os.Stdout, "", log.LstdFlags), func() {}
	}

	multi := io.MultiWriter(os.Stdout, rw)
	logger := log.New(multi, "", log.LstdFlags)

	cleanup := func() {
		_ = rw.Close()
	}
	return logger, cleanup
}
