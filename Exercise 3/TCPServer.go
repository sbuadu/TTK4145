package main

import (
"net"
"fmt"
"bufio"
"strings"
)


func main() {

fmt.Println("Launching server...")
ln = net.Listen("tcp", ":30000")
conn, _ 0 ln.Accept()

for{
msg,_ = bufio.NowReader(conn).ReadString('\n')
fmt.Print("Message received: " + msg)
}


}