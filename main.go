package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	os.Exit(run())
}

func run() int {
	interval := flag.Int("interval", 180, "スリープ防止の間隔(秒)")
	maxMove := flag.Int("maxmove", 5, "最大移動ピクセル数")
	flag.Parse()

	if *interval <= 0 {
		fmt.Fprintln(os.Stderr, "エラー: interval は正の整数を指定してください")
		return 1
	}

	logger, cleanup := setupLogger()
	defer cleanup()

	logger.Printf("設定: 間隔=%d秒, 最大移動距離=%dピクセル", *interval, *maxMove)
	logger.Println("Ctrl+Cで終了します")

	keepers := platformKeepers(*interval, *maxMove, logger)
	activeKeeper, err := tryKeepers(keepers, logger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		return 1
	}
	defer func() {
		if err := activeKeeper.Stop(); err != nil {
			logger.Printf("停止時にエラーが発生: %v\n", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Println("プログラムを終了します")
	return 0
}
