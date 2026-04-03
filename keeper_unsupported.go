//go:build !darwin && !windows

package main

func platformKeepers(interval, maxMove int) []Keeper {
	return []Keeper{}
}
