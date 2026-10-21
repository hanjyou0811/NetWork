package network

import (
	"fmt"
	"os"
)

func ip_init() int {
	if Net_protocol_register("IP", NET_PROTOCOL_TYPE_IP, ip_input) == -1 {
		fmt.Fprintln(os.Stderr, "failed to register IP protocol")
		return -1
	}
	return 0
}

func ip_input(data []byte, device *NetDevice) {
	fmt.Fprintf(os.Stderr, "IP input, dev=%s, len=%d\n", device.Name, len(data))
	Debugdump(data)
}
