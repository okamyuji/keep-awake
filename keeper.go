package main

import "fmt"

// Keeper はスリープ防止戦略を表すインターフェース。
// 各OS固有の実装がこのインターフェースを満たす。
type Keeper interface {
	Start() error
	Stop() error
	Name() string
}

// tryKeepers は渡された戦略リストを順番に試行し、最初に成功したものを返す。
// すべて失敗した場合はエラーを返す。
func tryKeepers(keepers []Keeper) (Keeper, error) {
	for _, k := range keepers {
		if err := k.Start(); err != nil {
			fmt.Printf("[%s] 利用不可: %v\n", k.Name(), err)
			continue
		}
		fmt.Printf("[%s] でスリープ防止を開始しました\n", k.Name())
		return k, nil
	}
	return nil, fmt.Errorf("利用可能なスリープ防止方法がありません")
}
