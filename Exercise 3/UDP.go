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

	n, addr, _ := conn.ReadFromUDP(buffer)
	// print info
	fmt.Println(string(buffer[:n]))
	fmt.Println(addr)
	conn.Close()

	sendserver, _ := net.ResolveUDPAddr("udp", "129.241.187.255:20009")
	fmt.Println("Sending to: ", sendserver)
	sendconn, _ := net.DialUDP("udp", server, sendserver)
	_,err := sendconn.Write([]byte("Hello"))
	fmt.Println(err)
	n1, addr1, _ := sendconn.ReadFromUDP(buffer2)
	fmt.Println(string(buffer2[:n1]))
	fmt.Println(addr1) 
}

