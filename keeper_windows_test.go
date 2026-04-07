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

func TestPlatformKeepers_Windows(t *testing.T) {
	logger := log.New(io.Discard, "", 0)
	keepers := platformKeepers(180, 5, logger)
	if len(keepers) == 0 {
		t.Fatal("expected at least one keeper for windows")
	}
	if keepers[0].Name() != "mouse-move" {
		t.Errorf("expected first keeper to be 'mouse-move', got '%s'", keepers[0].Name())
	}
}
