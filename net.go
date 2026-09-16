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

type NetDevice struct {
	index    uint
	Name     string
	typeid   uint16
	mtu      uint
	flags    uint16
	hlen     uint16
	alen     uint16
	addr     []uint8
	broadcat []uint8
}

var Devices []*NetDevice

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
	return device.flags&NET_DEVICE_FLAG_UP == 1
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
	device.flags |= NET_DEVICE_FLAG_UP
	return 0
}

func net_device_close(device *NetDevice) int {
	fmt.Printf("closing device: %s\n", device.Name)
	if net_device_is_up(device) == false {
		fmt.Fprintf(os.Stderr, "not opened,%s", device.Name)
		return -1
	}
	device.flags &= ^uint16(NET_DEVICE_FLAG_UP)
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

func debugdump(data []uint8) {
	fmt.Print(hex.Dump(data))
}

func Net_device_output(device *NetDevice, typeid uint16, data []uint8) int {
	length := len(data)
	fmt.Fprintf(os.Stderr, "dev=%s, type=0x%04x, len=%d\n", device.Name, typeid, length)
	debugdump(data)
	if !net_device_is_up(device) {
		fmt.Fprintf(os.Stderr, "not opened, dev=%s\n", device.Name)
		return -1
	}
	if device.mtu < uint(length) {
		fmt.Fprintf(os.Stderr, "too long, dev=%s, mtu=%d, len=%d\n", device.Name, device.mtu, length)
		return -1
	}
	return 0
}
