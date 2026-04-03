//go:build windows

package main

import (
	"testing"
	"time"
	"unsafe"
)

type MockProc struct {
	CallFunc func(...uintptr) (uintptr, uintptr, error)
}

func (m *MockProc) Call(a ...uintptr) (uintptr, uintptr, error) {
	return m.CallFunc(a...)
}

type MockDLL struct {
	Procs map[string]*MockProc
}

func (dll *MockDLL) NewProc(name string) *MockProc {
	if proc, ok := dll.Procs[name]; ok {
		return proc
	}
	return &MockProc{}
}

func TestMouseMoveKeeper_Name(t *testing.T) {
	k := &mouseMoveKeeper{interval: 1}
	if k.Name() != "mouse-move" {
		t.Errorf("expected 'mouse-move', got '%s'", k.Name())
	}
}

func TestMouseMoveKeeper_StartStop(t *testing.T) {
	mockDLL := &MockDLL{
		Procs: map[string]*MockProc{
			"GetCursorPos": {
				CallFunc: func(a ...uintptr) (uintptr, uintptr, error) {
					pt := (*POINT)(unsafe.Pointer(a[0]))
					pt.X, pt.Y = 100, 200
					return 0, 0, nil
				},
			},
			"SetCursorPos": {
				CallFunc: func(a ...uintptr) (uintptr, uintptr, error) {
					return 0, 0, nil
				},
			},
		},
	}

	procGetCursorPos = mockDLL.NewProc("GetCursorPos")
	procSetCursorPos = mockDLL.NewProc("SetCursorPos")

	k := &mouseMoveKeeper{interval: 1}
	if err := k.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	<-time.After(1500 * time.Millisecond)

	if err := k.Stop(); err != nil {
		t.Fatalf("Stop failed: %v", err)
	}
}

func TestPlatformKeepers_Windows(t *testing.T) {
	keepers := platformKeepers(180, 5)
	if len(keepers) == 0 {
		t.Fatal("expected at least one keeper for windows")
	}
	if keepers[0].Name() != "mouse-move" {
		t.Errorf("expected first keeper to be 'mouse-move', got '%s'", keepers[0].Name())
	}
}
