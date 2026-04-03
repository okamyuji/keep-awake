package main

import (
	"errors"
	"testing"
)

type mockKeeper struct {
	name     string
	startErr error
	stopped  bool
}

func (m *mockKeeper) Start() error { return m.startErr }
func (m *mockKeeper) Stop() error  { m.stopped = true; return nil }
func (m *mockKeeper) Name() string { return m.name }

func TestTryKeepers_FirstSucceeds(t *testing.T) {
	k1 := &mockKeeper{name: "first"}
	k2 := &mockKeeper{name: "second"}

	got, err := tryKeepers([]Keeper{k1, k2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name() != "first" {
		t.Errorf("expected 'first', got '%s'", got.Name())
	}
}

func TestTryKeepers_FallbackToSecond(t *testing.T) {
	k1 := &mockKeeper{name: "first", startErr: errors.New("unavailable")}
	k2 := &mockKeeper{name: "second"}

	got, err := tryKeepers([]Keeper{k1, k2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name() != "second" {
		t.Errorf("expected 'second', got '%s'", got.Name())
	}
}

func TestTryKeepers_AllFail(t *testing.T) {
	k1 := &mockKeeper{name: "first", startErr: errors.New("fail1")}
	k2 := &mockKeeper{name: "second", startErr: errors.New("fail2")}

	_, err := tryKeepers([]Keeper{k1, k2})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestTryKeepers_Empty(t *testing.T) {
	_, err := tryKeepers([]Keeper{})
	if err == nil {
		t.Fatal("expected error for empty keepers, got nil")
	}
}
