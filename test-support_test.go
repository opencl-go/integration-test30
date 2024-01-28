package cl30_test

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/opencl-go/cl30"
	"github.com/opencl-go/simulator-client/pkg/client"
)

type SimulatorClient struct {
	platform cl30.PlatformID
	client   *client.Client
}

var noSimulator sync.Once

func findAndConnectToSimulator(tb testing.TB) *SimulatorClient {
	tb.Helper()
	platformIDs, err := cl30.PlatformIDs()
	if (err != nil) && !errors.Is(err, cl30.StatusError(-1001)) {
		tb.Errorf("failed to query platform IDs: %v", err)
	}
	for _, platformID := range platformIDs {
		address, _ := cl30.PlatformInfoString(platformID, client.PlatformServerAddressInfo)
		if len(address) == 0 {
			continue
		}
		rawClient, err := client.NewClient(address)
		if err != nil {
			tb.Logf("Could not connect to '%s'", address)
			continue
		}
		simulator := &SimulatorClient{
			platform: platformID,
			client:   rawClient,
		}
		tb.Cleanup(func() {
			simulator.client.Disconnect()
		})
		return simulator
	}
	noSimulator.Do(func() {
		tb.Logf("No simulator platform found. Did you install the simulator-runtime library and register it?")
	})
	tb.SkipNow()
	return nil
}

func (sim *SimulatorClient) PrepareExtensionFunction(ctx context.Context, tb testing.TB) (string, uintptr) {
	tb.Helper()
	name, addr, err := sim.client.Platform().PrepareExtensionFunction(ctx)
	if err != nil {
		tb.Fatalf("failed to prepare function: %v", err)
	}
	tb.Cleanup(func() { _ = sim.client.Platform().ReleaseExtensionFunction(ctx, name) })
	return name, addr
}

func (sim *SimulatorClient) CreateDevice(ctx context.Context, tb testing.TB) cl30.DeviceID {
	tb.Helper()
	deviceID, err := sim.client.Devices().Create(ctx)
	if err != nil {
		tb.Fatalf("failed to create device: %v", err)
	}
	tb.Cleanup(func() { _ = sim.client.Devices().Delete(ctx, deviceID) })
	return cl30.DeviceID(deviceID)
}

func testContext(t *testing.T) context.Context {
	t.Helper()
	ctx := context.Background()
	if deadline, hasDeadline := t.Deadline(); hasDeadline {
		var cancel func()
		ctx, cancel = context.WithDeadline(ctx, deadline)
		t.Cleanup(cancel)
	}
	return ctx
}
