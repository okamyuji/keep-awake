//go:build windows

package main

import (
	"fmt"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

type caller interface {
	Call(a ...uintptr) (uintptr, uintptr, error)
}

var (
	user32           = windows.NewLazyDLL("user32.dll")
	procGetCursorPos caller = user32.NewProc("GetCursorPos")
	procSetCursorPos caller = user32.NewProc("SetCursorPos")
)

type POINT struct {
	X int32
	Y int32
}

func getCursorPos() (x, y int32) {
	var pt POINT
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
	return pt.X, pt.Y
}

func setCursorPos(x, y int32) {
	procSetCursorPos.Call(uintptr(x), uintptr(y))
}

type mouseMoveKeeper struct {
	interval int
	done     chan struct{}
}

func (k *mouseMoveKeeper) Name() string { return "mouse-move" }

func (k *mouseMoveKeeper) Start() error {
	k.done = make(chan struct{})
	go func() {
		ticker := time.NewTicker(time.Duration(k.interval) * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-k.done:
				return
			case <-ticker.C:
				x, y := getCursorPos()
				setCursorPos(x+1, y)
				time.Sleep(100 * time.Millisecond)
				setCursorPos(x, y)
				fmt.Printf("マウスを移動: (%d, %d) -> 1px右 -> 元の位置\n", x, y)
			}
		}
	}()
	return nil
}

func (k *mouseMoveKeeper) Stop() error {
	if k.done != nil {
		close(k.done)
	}
	return nil
}

func platformKeepers(interval, maxMove int) []Keeper {
	return []Keeper{&mouseMoveKeeper{interval: interval}}
}
