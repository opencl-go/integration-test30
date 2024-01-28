package cl30_test

import (
	"testing"
	"unsafe"

	"github.com/opencl-go/cl30"
)

func TestDeviceID(t *testing.T) {
	ctx := testContext(t)
	simulator := findAndConnectToSimulator(t)
	wantID := simulator.CreateDevice(ctx, t)

	gotIDs, err := cl30.DeviceIDs(simulator.platform, cl30.DeviceTypeAll)
	if err != nil {
		t.Errorf("failed to retrieve device IDs: %v", err)
	}
	found := false
	for _, gotID := range gotIDs {
		if gotID == wantID {
			found = true
		}
	}
	if !found {
		t.Errorf("failed to retrieve device: got: %+v, want: %v", gotIDs, wantID)
	}
}

func TestDevicePlatform(t *testing.T) {
	ctx := testContext(t)
	simulator := findAndConnectToSimulator(t)
	deviceID := simulator.CreateDevice(ctx, t)
	var platform cl30.PlatformID
	infoSize := unsafe.Sizeof(platform)

	sizeReturned, err := cl30.DeviceInfo(deviceID, cl30.DevicePlatformInfo, infoSize, unsafe.Pointer(&platform))
	if err != nil {
		t.Errorf("failed to retrieve info: %v", err)
	}
	if sizeReturned != infoSize {
		t.Errorf("failed to match info size: got: %v, want: %v", sizeReturned, infoSize)
	}
	if platform != simulator.platform {
		t.Errorf("failed to match platform: got: %v, want: %v", platform, simulator.platform)
	}
}

func TestDeviceVendorInfo(t *testing.T) {
	ctx := testContext(t)
	simulator := findAndConnectToSimulator(t)
	deviceID := simulator.CreateDevice(ctx, t)

	got, err := cl30.DeviceInfoString(deviceID, cl30.DeviceVendorInfo)
	if err != nil {
		t.Errorf("failed to retrieve info string: %v", err)
	}
	want := "github.com/opencl-go"
	if got != want {
		t.Errorf("failed vendor info: got: '%s', want: '%s'", got, want)
	}
}
