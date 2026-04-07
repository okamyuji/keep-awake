//go:build !darwin && !windows

package main

import "log"

func platformKeepers(_, _ int, _ *log.Logger) []Keeper {
	return []Keeper{}
}
