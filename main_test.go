package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

type MockProc struct {
	CallFunc func(...uintptr) (uintptr, uintptr, error)
}

func (m *MockProc) Call(a ...uintptr) (uintptr, uintptr, error) {
	return m.CallFunc(a...)
}

// モックを使用して getCursorPos と setCursorPos の動作を検証
type MockDLL struct {
	Procs map[string]*MockProc
}

func (dll *MockDLL) NewProc(name string) *MockProc {
	if proc, ok := dll.Procs[name]; ok {
		return proc
	}
	return &MockProc{}
}

func TestMainLogic(t *testing.T) {
	// モックDLLのセットアップ
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

	// DLLをモックに置き換え
	user32 = &windows.LazyDLL{}
	procGetCursorPos = mockDLL.NewProc("GetCursorPos")
	procSetCursorPos = mockDLL.NewProc("SetCursorPos")

	// コマンドライン引数を設定
	os.Args = []string{"cmd", "-interval=1", "-maxmove=1"}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// シグナルのモック
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)
	defer signal.Stop(sigChan)

	// タイマーを短縮してテスト
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	go func() {
		time.Sleep(2 * time.Second)
		sigChan <- syscall.SIGINT
	}()

	// メインロジックのテスト
	mainLoop := func() {
		for {
			select {
			case <-ticker.C:
				x, y := getCursorPos()
				if x != 100 || y != 200 {
					t.Fatalf("unexpected cursor position: (%d, %d)", x, y)
				}
				setCursorPos(x+1, y)
				setCursorPos(x, y)
			case <-sigChan:
				return
			}
		}
	}

	mainLoop()
}

func TestGetCursorPos(t *testing.T) {
	testX, testY := getCursorPos()
	if testX != 100 || testY != 200 {
		t.Errorf("Expected cursor position (100, 200), got (%d, %d)", testX, testY)
	}
}

func TestSetCursorPos(t *testing.T) {
	setCursorPos(300, 400)
	// マウスカーソルの設定結果を検証するためのモックログを確認する手段がある場合に追加する
}
