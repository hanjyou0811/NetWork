package network

import (
	"encoding/hex"
	"fmt"
	"os"
)

const NET_DEVICE_TYPE_DUMMY = 0x0000
const NET_DEVICE_TYPE_LOOPBACK = 0x0001
const NET_DEVICE_TYPE_ETHERNET = 0x0002

const NET_DEVICE_FLAG_UP = 0x0001
const NET_DEVICE_FLAG_LOOPBACK = 0x0010
const NET_DEVICE_FLAG_BROADCAST = 0x0020
const NET_DEVICE_FLAG_P2P = 0x0040
const NET_DEVICE_FLAG_NEED_ARP = 0x0100

const NET_PROTOCOL_TYPE_IP = 0x0800
const NET_PROTOCOL_TYPE_ARP = 0x0806
const NET_PROTOCOL_TYPE_IPV6 = 0x86DD

type NetDevice struct {
	index     uint
	Name      string
	typeid    uint16
	mtu       uint
	Flags     uint16
	hlen      uint16
	alen      uint16
	addr      []uint8
	broadcast []uint8
	Ops       *NetDeviceOps
	Private   interface{}
}

type NetDeviceOps struct {
	Open   func(device *NetDevice) int
	Close  func(device *NetDevice) int
	Output func(device *NetDevice, typeid uint16, data []uint8) int
}

var Devices []*NetDevice

type NetProtocol struct {
	Name    string
	Type    uint16
	Handler NetProtocolHandler
}

var Protocols []*NetProtocol

type NetProtocolHandler func(data []byte, dev *NetDevice)

func NewNetDevice(name string, typeid uint16, mtu uint, hlen uint16, alen uint16) *NetDevice {
	return &NetDevice{
		Name:   name,
		typeid: typeid,
		mtu:    mtu,
		hlen:   hlen,
		alen:   alen,
	}
}

func net_device_is_up(device *NetDevice) bool {
	return device.Flags&NET_DEVICE_FLAG_UP == 1
}

func Net_device_register(device *NetDevice) int {
	Devices = append(Devices, device)
	return 0
}

func net_device_open(device *NetDevice) int {
	fmt.Printf("dev=%s is opening\n", device.Name)
	if net_device_is_up(device) {
		fmt.Fprintf(os.Stderr, "already opened, dev=%s\n", device.Name)
		return -1
	}
	device.Flags |= NET_DEVICE_FLAG_UP
	return 0
}

func net_device_close(device *NetDevice) int {
	fmt.Printf("closing device: %s\n", device.Name)
	if net_device_is_up(device) == false {
		fmt.Fprintf(os.Stderr, "not opened,%s", device.Name)
		return -1
	}
	device.Flags &= ^uint16(NET_DEVICE_FLAG_UP)
	return 0
}

func Net_run() int {
	fmt.Println("start up ...")

	for i := 0; i < len(Devices); i++ {
		net_device_open(Devices[i])
	}
	fmt.Println("success")
	return 0
}

func net_shutdown() int {
	fmt.Println("shutting down ...")

	for i := 0; i < len(Devices); i++ {
		net_device_close(Devices[i])
	}
	fmt.Println("success")
	return 0
}

func Debugdump(data []uint8) {
	fmt.Print(hex.Dump(data))
}

func Net_device_input(device *NetDevice) int {
	fmt.Fprintf(os.Stderr, "dev=%s\n", device.Name)
	if !net_device_is_up(device) {
		fmt.Fprintf(os.Stderr, "not opened, dev=%s\n", device.Name)
		return -1
	}
	if device.Ops != nil && device.Ops.Open != nil {
		if device.Ops.Open(device) != 0 {
			fmt.Fprintf(os.Stderr, "failed to open device, dev=%s\n", device.Name)
			return -1
		}
	}
	device.Flags &^= uint16(NET_DEVICE_FLAG_UP)
	return 0
}

func Net_device_close(device *NetDevice) int {
	fmt.Fprintf(os.Stderr, "dev=%s\n", device.Name)
	if !net_device_is_up(device) {
		fmt.Fprintf(os.Stderr, "not opened, dev=%s\n", device.Name)
		return -1
	}
	if device.Ops != nil && device.Ops.Close != nil {
		if device.Ops.Close(device) != 0 {
			fmt.Fprintf(os.Stderr, "failed to close device, dev=%s\n", device.Name)
			return -1
		}
	}
	device.Flags &^= uint16(NET_DEVICE_FLAG_UP)
	return 0
}

func Net_device_output(device *NetDevice, typeid uint16, data []uint8) int {
	length := len(data)
	fmt.Fprintf(os.Stderr, "dev=%s, type=0x%04x, len=%d\n", device.Name, typeid, length)
	Debugdump(data)
	if !net_device_is_up(device) {
		fmt.Fprintf(os.Stderr, "not opened, dev=%s\n", device.Name)
		return -1
	}
	if device.mtu < uint(length) {
		fmt.Fprintf(os.Stderr, "too long, dev=%s, mtu=%d, len=%d\n", device.Name, device.mtu, length)
		return -1
	}
	if device.Ops != nil && device.Ops.Output != nil {
		if device.Ops.Output(device, typeid, data) != 0 {
			fmt.Fprintf(os.Stderr, "failed to output data, dev=%s\n", device.Name)
			return -1
		}
	}
	return 0
}

func Net_input(device *NetDevice, typeid uint16, data []uint8) int {
	fmt.Fprintf(os.Stderr, "dev=%s, type=0x%04x, len=%d\n", device.Name, typeid, len(data))
	Debugdump(data)
	for _, proto := range Protocols {
		if typeid == proto.Type {
			proto.Handler(data, device)
			return 0
		}
	}
	return 0
}

func Net_protocol_register(name string, typeid uint16, handler NetProtocolHandler) int {
	for _, proto := range Protocols {
		if typeid == proto.Type {
			fmt.Fprintf(
				os.Stderr,
				"already registered, type=%s(0x%04x), exist=%s(0x%04x)\n",
				name, typeid, proto.Name, proto.Type,
			)
			return -1
		}
	}

	proto := &NetProtocol{
		Name:    name,
		Type:    typeid,
		Handler: handler,
	}
	Protocols = append(Protocols, proto)

	fmt.Printf("registered, type=%s(0x%04x)\n", proto.Name, typeid)
	return 0
}

func Net_init() int {
	fmt.Println("initialize network stack ...")
	if ip_init() == -1 {
		fmt.Fprintln(os.Stderr, "failed to initialize IP protocol")
		return -1
	}
	fmt.Println("success")
	return 0
}
