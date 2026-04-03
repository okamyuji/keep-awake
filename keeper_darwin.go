//go:build darwin

package main

import (
	"fmt"
	"os/exec"
	"syscall"
)

type caffeinateKeeper struct {
	cmd *exec.Cmd
}

func (k *caffeinateKeeper) Name() string { return "caffeinate" }

func (k *caffeinateKeeper) Start() error {
	path, err := exec.LookPath("caffeinate")
	if err != nil {
		return fmt.Errorf("caffeinate が見つかりません: %w", err)
	}
	k.cmd = exec.Command(path, "-di")
	if err := k.cmd.Start(); err != nil {
		return fmt.Errorf("caffeinate の起動に失敗: %w", err)
	}
	return nil
}

func (k *caffeinateKeeper) Stop() error {
	if k.cmd != nil && k.cmd.Process != nil {
		_ = k.cmd.Process.Signal(syscall.SIGTERM)
		err := k.cmd.Wait()
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				if status, ok := exitErr.Sys().(syscall.WaitStatus); ok && status.Signaled() {
					return nil
				}
			}
			return err
		}
		return nil
	}
	return nil
}

func platformKeepers(interval, maxMove int) []Keeper {
	return []Keeper{&caffeinateKeeper{}}
}
