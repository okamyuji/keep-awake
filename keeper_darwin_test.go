//go:build darwin

package main

import (
	"os/exec"
	"testing"
)

func TestCaffeinateKeeper_Name(t *testing.T) {
	k := &caffeinateKeeper{}
	if k.Name() != "caffeinate" {
		t.Errorf("expected 'caffeinate', got '%s'", k.Name())
	}
}

func TestCaffeinateKeeper_StartStop(t *testing.T) {
	if _, err := exec.LookPath("caffeinate"); err != nil {
		t.Skip("caffeinate not found, skipping")
	}

	k := &caffeinateKeeper{}
	if err := k.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	if k.cmd == nil || k.cmd.Process == nil {
		t.Fatal("caffeinate process not started")
	}
	pid := k.cmd.Process.Pid
	if pid <= 0 {
		t.Fatalf("invalid pid: %d", pid)
	}

	if err := k.Stop(); err != nil {
		t.Fatalf("Stop failed: %v", err)
	}
}

func TestCaffeinateKeeper_StopWithoutStart(t *testing.T) {
	k := &caffeinateKeeper{}
	if err := k.Stop(); err != nil {
		t.Fatalf("Stop without Start should not error: %v", err)
	}
}

func TestPlatformKeepers_Darwin(t *testing.T) {
	keepers := platformKeepers(180, 5)
	if len(keepers) == 0 {
		t.Fatal("expected at least one keeper for darwin")
	}
	if keepers[0].Name() != "caffeinate" {
		t.Errorf("expected first keeper to be 'caffeinate', got '%s'", keepers[0].Name())
	}
}
