package tests

import (
	"Kaban/iternal/Service/Handlers"
	"testing"
)

func TestSwapKeys(t *testing.T) {
	var i = 1
	tests := []struct {
		name string
	}{
		{name: "TestSwapKeys"},
		{name: "TestSwapKeys"},
		{name: "TestSwapKeys"},
		{name: "TestSwapKeys"},
		{name: "TestSwapKeys"},
		{name: "TestSwapKeys"},
		{name: "TestSwapKeys"},
		{name: "TestSwapKeys"},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Logf("Test: %s id : %v", tt.name, i)
		i++
		t.Run(tt.name, func(t *testing.T) {
			ok := Handlers.SwapKeys()
			if !ok {
				t.Error("Handlers.SwapKeys() failed")
			}
		})
	}
}
