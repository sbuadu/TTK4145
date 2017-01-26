package main

import (
"fmt"
"net"
"bufio"
"time"

)

const (

CONN_PORT = "33546"
CONN_TYPE = "tcp"

)

func main(){

conn, _:= net.Dial(CONN_TYPE, CONN_PORT)

for {

fmt.Fprintf(conn, "Hello everybody!\0 ")

reply := bufio.NewReader(conn).ReadString('\n')
fmt.Print("message from server: " + reply)
time.Sleep(10*time.Millisecond)

}
	
}
