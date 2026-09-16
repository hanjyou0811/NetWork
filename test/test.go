package main

import (
	"fmt"
	network "github.com/hanjyou0811/NetWork"
	"os"
	"time"
)

var terminate = 0
var test_data = []byte{
	0x45, 0x00, 0x00, 0x30,
	0x00, 0x80, 0x00, 0x00,
	0xff, 0x01, 0xbd, 0x4a,
	0x7f, 0x00, 0x00, 0x01,
	0x7f, 0x00, 0x00, 0x01,
	0x08, 0x00, 0x35, 0x64,
	0x00, 0x80, 0x00, 0x01,
	0x31, 0x32, 0x33, 0x34,
	0x35, 0x36, 0x37, 0x38,
	0x39, 0x30, 0x21, 0x40,
	0x23, 0x24, 0x25, 0x5e,
	0x26, 0x2a, 0x28, 0x29,
}

func dummy_init() *network.NetDevice {
	device := network.NewNetDevice("dummy", network.NET_DEVICE_TYPE_DUMMY, 128, 0, 0)
	if network.Net_device_register(device) == -1 {
		fmt.Fprintf(os.Stderr, "failed to register device\n")
		os.Exit(-1)
	}
	fmt.Printf("success dev=%s\n", device.Name)
	return device
}

func app_main() int {
	fmt.Fprintf(os.Stderr, "press Ctrl+C to terminate\n")
	dev := dummy_init()
	if network.Net_run() == -1 {
		fmt.Fprintf(os.Stderr, "failed to start network\n")
		os.Exit(-1)
	}
	for terminate == 0 {
		if network.Net_device_output(dev, 0x0800, test_data) == -1 {
			fmt.Fprintf(os.Stderr, "failed to send data\n")
			break
		}
		time.Sleep(time.Second)
	}
	return 0
}

func main() {
	app_main()
}
