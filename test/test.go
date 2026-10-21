package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	network "github.com/hanjyou0811/NetWork"
	"github.com/hanjyou0811/NetWork/driver"
)

var (
	dev      *network.NetDevice
	signalCh = make(chan os.Signal, 1)
)
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
		return nil
	}
	fmt.Printf("success dev=%s\n", device.Name)
	return device
}

func setup() int {
	signal.Notify(signalCh, os.Interrupt)
	fmt.Fprintln(os.Stderr, "setup protocol stack...")
	// Net_init は未実装のため、ここではデバイスの初期化から行う。
	dev = driver.Loopback_init()
	if dev == nil {
		fmt.Fprintln(os.Stderr, "loopback_init() failure")
		return -1
	}
	if network.Net_run() == -1 {
		fmt.Fprintf(os.Stderr, "failed to start network\n")
		return -1
	}
	return 0
}

func app_main() int {
	defer signal.Stop(signalCh)
	if setup() == -1 {
		return 1
	}
	fmt.Fprintln(os.Stderr, "press Ctrl+C to terminate")
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-signalCh:
			return 0
		default:
		}
		if network.Net_device_output(dev, 0x0800, test_data) == -1 {
			fmt.Fprintf(os.Stderr, "failed to send data\n")
			return 1
		}
		select {
		case <-signalCh:
			return 0
		case <-ticker.C:
		}
	}
}

func main() {
	os.Exit(app_main())
}
