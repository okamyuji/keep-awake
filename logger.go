package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

const logFileName = "keep-awake.log"

// setupLogger creates a logger that writes to both stdout and a log file.
// The log file is truncated on each startup (previous logs are cleared).
// Returns the logger and a cleanup function.
func setupLogger() (*log.Logger, func()) {
	dir, err := os.Getwd()
	if err != nil {
		dir = "."
	}

	path := filepath.Join(dir, logFileName)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "警告: ログファイルの作成に失敗: %v (stdoutのみに出力します)\n", err)
		return log.New(os.Stdout, "", log.LstdFlags), func() {}
	}

	multi := io.MultiWriter(os.Stdout, f)
	logger := log.New(multi, "", log.LstdFlags)

	cleanup := func() {
		_ = f.Close()
	}
	return logger, cleanup
}
