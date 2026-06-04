package main

import "testing"

func TestPlaceholder(t *testing.T) {
    if got := 1 + 1; got != 2 {
        t.Fatalf("expected 2, got %d", got)
    }
}

func TestAlwaysPass(t *testing.T) {
    got := true
    if !got {
        t.Fatal("this should never happen")
    }
}

func TestStringCheck(t *testing.T) {
    got := "go"
    if got != "go" {
        t.Fatal("string comparison failed")
    }
}
