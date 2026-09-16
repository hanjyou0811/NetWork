package main

import (
	"fmt"
	"os"

	network "github.com/hanjyou0811/NetWork"
)

var terminate = 0

func dummy_init() {
	device := network.NewNetDevice("dummy", network.NET_DEVICE_TYPE_DUMMY, 128, 0, 0)
	if network.Net_device_register(device) == -1 {
		fmt.Fprintf(os.Stderr, "success, dev=%s", device.Name)
		return
	}
	fmt.Println("success dev=%s", device.Name)
}

func app_main() int {
	fmt.Fprintf(os.Stderr, "press Ctrl+C to terminate")
    for !terminate {
        if (network.Net_device_output(dev, 0x0800, )))
    }
    
    return 0
}
