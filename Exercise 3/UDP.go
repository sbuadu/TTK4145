package main

import (
	"fmt"
	"net"
)


func main() {
	// udp ip addr
	server, _ := net.ResolveUDPAddr("udp", "129.241.187.255:30000")
	// create socket
	conn, _ := net.ListenUDP("udp", server)
	// create buffer to store msg
	buffer := make([]byte, 1024)
	buffer2 := make([]byte, 1024)
	// read from server addr
	conn2, _ := net.DialUDP("udp", server, "129.241.187.255:20009")
	conn2.WriteToUDP("Hello",server)

	for {
		n, addr, _ := conn.ReadFromUDP(buffer)
		n1, addr2, _ := conn2.ReadFromUDP(buffer2)
		// print info
		fmt.Println(string(buffer[:n]))
		fmt.Println(addr)
}
}

