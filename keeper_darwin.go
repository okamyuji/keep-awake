//go:build darwin

package main

import (
	"fmt"
	"log"
	"os/exec"
	"sync"
	"syscall"
)

type caffeinateKeeper struct {
	cmd    *exec.Cmd
	logger *log.Logger
	once   sync.Once
}

func (k *caffeinateKeeper) Name() string { return "caffeinate" }

func (k *caffeinateKeeper) Start() error {
	if k.cmd != nil && k.cmd.Process != nil {
		return fmt.Errorf("caffeinate は既に起動済みです")
	}
	path, err := exec.LookPath("caffeinate")
	if err != nil {
		return fmt.Errorf("caffeinate が見つかりません: %w", err)
	}
	k.cmd = exec.Command(path, "-di")
	if err := k.cmd.Start(); err != nil {
		return fmt.Errorf("caffeinate の起動に失敗: %w", err)
	}
	k.logger.Println("caffeinate プロセスを起動しました (PID:", k.cmd.Process.Pid, ")")
	return nil
}

func (k *caffeinateKeeper) Stop() error {
	var stopErr error
	k.once.Do(func() {
		if k.cmd != nil && k.cmd.Process != nil {
			_ = k.cmd.Process.Signal(syscall.SIGTERM)
			err := k.cmd.Wait()
			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					if status, ok := exitErr.Sys().(syscall.WaitStatus); ok && status.Signaled() {
						k.logger.Println("caffeinate プロセスを停止しました")
						return
					}
				}
				stopErr = err
				return
			}
			k.logger.Println("caffeinate プロセスを停止しました")
		}
	})
	return stopErr
}

func platformKeepers(interval, maxMove int, logger *log.Logger) []Keeper {
	return []Keeper{&caffeinateKeeper{logger: logger}}
}
