package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32           = windows.NewLazyDLL("user32.dll")
	procGetCursorPos = user32.NewProc("GetCursorPos")
	procSetCursorPos = user32.NewProc("SetCursorPos")
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

func main() {
	// コマンドライン引数の設定
	interval := flag.Int("interval", 180, "マウス移動の間隔(秒)")
	maxMove := flag.Int("maxmove", 5, "最大移動ピクセル数")
	flag.Parse()

	// 現在の設定を表示
	fmt.Printf("設定:\n - 間隔: %d秒\n - 最大移動距離: %dピクセル\n", *interval, *maxMove)
	fmt.Println("Ctrl+Cで終了します")

	// シグナルハンドリングの設定
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// タイマーの設定
	ticker := time.NewTicker(time.Duration(*interval) * time.Second)
	defer ticker.Stop()

	// メインループ
	for {
		select {
		case <-ticker.C:
			// 現在のマウス位置を取得
			x, y := getCursorPos()

			// 1ピクセル右に移動して戻す
			setCursorPos(x+1, y)
			time.Sleep(100 * time.Millisecond)
			setCursorPos(x, y)

			fmt.Printf("マウスを移動: (%d, %d) -> 1px右 -> 元の位置\n", x, y)

		case <-sigChan:
			fmt.Println("\nプログラムを終了します")
			return
		}
	}
}
