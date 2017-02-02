package main

import (
	"fmt"
	"net"
)


func ListenForIp() {
	packet, _ := net.ListenPacket("udp",":30000")
	//addr, _ := hex.Dedoce(i, packet)
	fmt.Println("Hello from the network module")
	fmt.Printf("%p\n", packet)
	return packet
}

