package gobacktest

import (
	"fmt"
	"testing"
)

// TestBack test backTest interface
func TestBack(t *testing.T) {
	fmt.Println("create new backTest")
	backConfig := make(map[string]string)
	back := NewBackTest(backConfig)

	back.Next()
}
