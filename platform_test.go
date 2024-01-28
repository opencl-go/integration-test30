package cl30_test

import (
	"testing"

	"github.com/opencl-go/cl30"
)

func TestPlatformExtensionAddress(t *testing.T) {
	ctx := testContext(t)
	simulator := findAndConnectToSimulator(t)
	name, wantAddr := simulator.PrepareExtensionFunction(ctx, t)

	gotAddr := uintptr(cl30.ExtensionFunctionAddressForPlatform(simulator.platform, name))
	if gotAddr != wantAddr {
		t.Errorf("failed to retrieve address: got %v, want: %v", gotAddr, wantAddr)
	}
}
