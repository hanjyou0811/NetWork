package network

import "testing"

func TestNetRunOpensRegisteredDevice(t *testing.T) {
	Devices = nil
	t.Cleanup(func() { Devices = nil })

	device := NewNetDevice("dummy", NET_DEVICE_TYPE_DUMMY, 128, 0, 0)
	if got := Net_device_register(device); got != 0 {
		t.Fatalf("Net_device_register() = %d, want 0", got)
	}
	if got := Net_run(); got != 0 {
		t.Fatalf("Net_run() = %d, want 0", got)
	}

	if got := Net_device_output(device, 0x0800, []byte{0x45}); got != 0 {
		t.Fatalf("Net_device_output() after Net_run() = %d, want 0", got)
	}
}
