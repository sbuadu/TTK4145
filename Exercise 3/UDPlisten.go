package main

import (
	"fmt"
	"net"
)


func main() {
	// udp ip addr
	server, _ := net.ResolveUDPAddr("udp", "129.241.187.255:30000")
//	server2, err := net.ResolveUDPAddr("udp", "129.241.187.255:20097")
//	fmt.Println(err)
	// create socket
	conn, _ := net.ListenUDP("udp", server)
	// create buffer to store msg
	buffer := make([]byte, 1024)
//	buffer2 := make([]byte, 1024)
	// read from server addr
//	conn2, err2 := net.DialUDP("udp", server, server2)
//	fmt.Println(err2)
//	conn2.WriteToUDP([]byte("Hello"),server)

		n, addr, _ := conn.ReadFromUDP(buffer)
//		n1, addr2, _ := conn2.ReadFromUDP(buffer2)
		// print info
		fmt.Println(string(buffer[:n]))
		fmt.Println(addr)
//		fmt.Println(string(buffer2[:n1]))
//		fmt.Println(addr2)

// conn2.Close()
conn.Close()
}

