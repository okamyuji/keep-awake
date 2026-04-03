package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	interval := flag.Int("interval", 180, "スリープ防止の間隔(秒)")
	maxMove := flag.Int("maxmove", 5, "最大移動ピクセル数")
	flag.Parse()

	fmt.Printf("設定:\n - 間隔: %d秒\n - 最大移動距離: %dピクセル\n", *interval, *maxMove)
	fmt.Println("Ctrl+Cで終了します")

	keepers := platformKeepers(*interval, *maxMove)
	activeKeeper, err := tryKeepers(keepers)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}
	defer activeKeeper.Stop()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\nプログラムを終了します")
}
