//go:build windows

package main

import (
	"fmt"
	"log"
	"sync"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

type caller interface {
	Call(a ...uintptr) (uintptr, uintptr, error)
}

var (
	user32                  = windows.NewLazyDLL("user32.dll")
	procGetCursorPos caller = user32.NewProc("GetCursorPos")
	procSetCursorPos caller = user32.NewProc("SetCursorPos")
)

type POINT struct {
	X int32
	Y int32
}

func getCursorPos() (x, y int32, err error) {
	var pt POINT
	ret, _, sysErr := procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
	if ret == 0 {
		return 0, 0, fmt.Errorf("GetCursorPos の呼び出しに失敗しました: %v", sysErr)
	}
	return pt.X, pt.Y, nil
}

func setCursorPos(x, y int32) error {
	ret, _, sysErr := procSetCursorPos.Call(uintptr(x), uintptr(y))
	if ret == 0 {
		return fmt.Errorf("SetCursorPos の呼び出しに失敗しました: %v", sysErr)
	}
	return nil
}

type mouseMoveKeeper struct {
	interval int
	maxMove  int
	done     chan struct{}
	logger   *log.Logger
	mu       sync.Mutex
}

func (k *mouseMoveKeeper) Name() string { return "mouse-move" }

func (k *mouseMoveKeeper) Start() error {
	k.mu.Lock()
	defer k.mu.Unlock()
	if k.done != nil {
		return nil
	}
	k.done = make(chan struct{})
	done := k.done
	go func() {
		ticker := time.NewTicker(time.Duration(k.interval) * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				x, y, err := getCursorPos()
				if err != nil {
					k.logger.Printf("カーソル位置の取得に失敗: %v\n", err)
					continue
				}
				move := int32(k.maxMove)
				if move <= 0 {
					move = 1
				}
				if err := setCursorPos(x+move, y); err != nil {
					k.logger.Printf("カーソル移動に失敗: %v\n", err)
					continue
				}
				time.Sleep(100 * time.Millisecond)
				if err := setCursorPos(x, y); err != nil {
					k.logger.Printf("カーソル復帰に失敗: %v\n", err)
				}
				k.logger.Printf("マウスを移動: (%d, %d) -> %dpx右 -> 元の位置\n", x, y, move)
			}
		}
	}()
	return nil
}

func (k *mouseMoveKeeper) Stop() error {
	k.mu.Lock()
	defer k.mu.Unlock()
	if k.done == nil {
		return nil
	}
	close(k.done)
	k.done = nil
	return nil
}

func platformKeepers(interval, maxMove int, logger *log.Logger) []Keeper {
	return []Keeper{&mouseMoveKeeper{interval: interval, maxMove: maxMove, logger: logger}}
}
