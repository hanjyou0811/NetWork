package driver

import (
	"encoding/hex"
	"fmt"
	"os"

	network "github.com/hanjyou0811/NetWork"
)

const LOOPBACK_MTU = 65535

func loopback_output(device *network.NetDevice, typeid uint16, data []uint8) int {
	fmt.Fprintf(os.Stderr, "dev=%s, type=0x%04x, len=%d\n", device.Name, typeid, len(data))
	fmt.Print(hex.Dump(data))
	return network.Net_input(device, typeid, data)
}

var loopbackOps = network.NetDeviceOps{
	Output: loopback_output,
}

func Loopback_init() *network.NetDevice {
	dev := network.NewNetDevice(
		"loopback",
		network.NET_DEVICE_TYPE_LOOPBACK,
		LOOPBACK_MTU,
		0, // ヘッダなし
		0, // アドレスなし
	)
	dev.Flags |= network.NET_DEVICE_FLAG_LOOPBACK
	dev.Ops = &loopbackOps

	if network.Net_device_register(dev) == -1 {
		fmt.Fprintln(os.Stderr, "failed to register loopback device")
		return nil
	}

	fmt.Printf("success, dev=%s\n", dev.Name)
	return dev
}
