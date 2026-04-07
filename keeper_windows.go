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
	ret, _, _ := procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
	if ret == 0 {
		return 0, 0, fmt.Errorf("GetCursorPos の呼び出しに失敗しました")
	}
	return pt.X, pt.Y, nil
}

func setCursorPos(x, y int32) {
	procSetCursorPos.Call(uintptr(x), uintptr(y))
}

type mouseMoveKeeper struct {
	interval int
	maxMove  int
	done     chan struct{}
	logger   *log.Logger
	once     sync.Once
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
				x, y, err := getCursorPos()
				if err != nil {
					k.logger.Printf("カーソル位置の取得に失敗: %v\n", err)
					continue
				}
				move := int32(k.maxMove)
				if move <= 0 {
					move = 1
				}
				setCursorPos(x+move, y)
				time.Sleep(100 * time.Millisecond)
				setCursorPos(x, y)
				k.logger.Printf("マウスを移動: (%d, %d) -> %dpx右 -> 元の位置\n", x, y, move)
			}
		}
	}()
	return nil
}

func (k *mouseMoveKeeper) Stop() error {
	k.once.Do(func() {
		if k.done != nil {
			close(k.done)
		}
	})
	return nil
}

func platformKeepers(interval, maxMove int, logger *log.Logger) []Keeper {
	return []Keeper{&mouseMoveKeeper{interval: interval, maxMove: maxMove, logger: logger}}
}
