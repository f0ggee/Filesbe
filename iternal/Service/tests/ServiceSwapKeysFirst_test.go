package tests

import (
	"Kaban/iternal/Service/Handlers"
	"testing"
)

func TestSwapKeyFirst(t *testing.T) {
	var i = 0
	tests := []struct {
		name string
	}{
		{name: "TestSwapKeyFirst"},
		{name: "TestSwapKeyFirst"},
		{name: "TestSwapKeyFirst"},
		{name: "TestSwapKeyFirst"},
		{name: "TestSwapKeyFirst"},
		{name: "TestSwapKeyFirst"},
		{name: "TestSwapKeyFirst"},
		{name: "TestSwapKeyFirst"},
		{name: "TestSwapKeyFirst"},
		{name: "TestSwapKeyFirst"},
		{name: "TestSwapKeyFirst"},
		{name: "TestSwapKeyFirst"},
		{name: "TestSwapKeyFirst"},
	}
	for _, tt := range tests {
		t.Logf("Test: %s, i: %v", tt.name, i)
		i++
		t.Run(tt.name, func(t *testing.T) {
			Handlers.SwapKeyFirst()

		})
	}
}
