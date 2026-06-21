//go:build windows

package main

import (
	"io"
	"log"
	"testing"
	"time"
	"unsafe"
)

type mockProc struct {
	callFunc func(...uintptr) (uintptr, uintptr, error)
}

func (m *mockProc) Call(a ...uintptr) (uintptr, uintptr, error) {
	return m.callFunc(a...)
}

func TestMouseMoveKeeper_Name(t *testing.T) {
	k := &mouseMoveKeeper{interval: 1, maxMove: 5, logger: log.New(io.Discard, "", 0)}
	if k.Name() != "mouse-move" {
		t.Errorf("expected 'mouse-move', got '%s'", k.Name())
	}
}

func TestMouseMoveKeeper_StartStop(t *testing.T) {
	oldGet := procGetCursorPos
	oldSet := procSetCursorPos
	t.Cleanup(func() {
		procGetCursorPos = oldGet
		procSetCursorPos = oldSet
	})

	procGetCursorPos = &mockProc{callFunc: func(a ...uintptr) (uintptr, uintptr, error) {
		pt := (*POINT)(unsafe.Pointer(a[0]))
		pt.X, pt.Y = 100, 200
		return 1, 0, nil
	}}
	procSetCursorPos = &mockProc{callFunc: func(a ...uintptr) (uintptr, uintptr, error) {
		return 1, 0, nil
	}}

	k := &mouseMoveKeeper{interval: 1, maxMove: 5, logger: log.New(io.Discard, "", 0)}
	if err := k.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	<-time.After(1500 * time.Millisecond)

	if err := k.Stop(); err != nil {
		t.Fatalf("Stop failed: %v", err)
	}
}

func TestExecutionStateKeeper_Name(t *testing.T) {
	k := &executionStateKeeper{logger: log.New(io.Discard, "", 0)}
	if k.Name() != "execution-state" {
		t.Errorf("expected 'execution-state', got '%s'", k.Name())
	}
}

func TestExecutionStateKeeper_StartStop(t *testing.T) {
	oldProc := procSetThreadExecutionState
	t.Cleanup(func() { procSetThreadExecutionState = oldProc })

	var lastFlags uintptr
	procSetThreadExecutionState = &mockProc{callFunc: func(a ...uintptr) (uintptr, uintptr, error) {
		lastFlags = a[0]
		return 1, 0, nil
	}}

	k := &executionStateKeeper{logger: log.New(io.Discard, "", 0)}
	if err := k.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	wantStart := esContinuous | esSystemRequired | esDisplayRequired
	if lastFlags != wantStart {
		t.Errorf("Start flags = 0x%X, want 0x%X", lastFlags, wantStart)
	}

	if err := k.Stop(); err != nil {
		t.Fatalf("Stop failed: %v", err)
	}
	if lastFlags != esContinuous {
		t.Errorf("Stop flags = 0x%X, want 0x%X", lastFlags, esContinuous)
	}
}

func TestExecutionStateKeeper_StopWithoutStart(t *testing.T) {
	k := &executionStateKeeper{logger: log.New(io.Discard, "", 0)}
	if err := k.Stop(); err != nil {
		t.Fatalf("Stop without Start should not error: %v", err)
	}
}

func TestPlatformKeepers_Windows(t *testing.T) {
	logger := log.New(io.Discard, "", 0)
	keepers := platformKeepers(180, 5, logger)
	if len(keepers) < 2 {
		t.Fatal("expected at least two keepers for windows")
	}
	if keepers[0].Name() != "execution-state" {
		t.Errorf("expected first keeper to be 'execution-state', got '%s'", keepers[0].Name())
	}
	if keepers[1].Name() != "mouse-move" {
		t.Errorf("expected second keeper to be 'mouse-move', got '%s'", keepers[1].Name())
	}
}
