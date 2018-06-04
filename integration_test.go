package main

import "testing"

func TestSum(t *testing.T) {
  var v float64
  v = 2 + 2
  if v != 4 {
    t.Error("Expected 4, got ", v)
  }
}
