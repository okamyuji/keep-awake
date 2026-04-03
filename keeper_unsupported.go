//go:build !darwin && !windows

package main

import "log"

func platformKeepers(interval, maxMove int, logger *log.Logger) []Keeper {
	return []Keeper{}
}
