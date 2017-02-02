package main

import (
	"fmt"
	"net"
	"log"
	"os/exec"
	"strconv"
	"time"

)



type Counter struct{
	State int
}

type Message struct{
	Data string
}

func Transmitter(toBackup chan Message){
//	conn := net.DialBroadcastUDP("20009")
	addr,err := net.ResolveUDPAddr("udp4", ":20099")
if err != nil {
	log.Fatal(err)
}

conn, err := net.ListenUDP("udp", addr)
	if err != nil{
	log.Fatal(err)
	}

	for {
		state := <- toBackup
		_, err := conn.WriteToUDP([]byte(state.Data), addr)
if err != nil {
	log.Fatal(err)
}

}

}


func main(){

	fmt.Print("Let's count!\n\n ")
	counter := Counter{0}

spawnBackup := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run backup.go")

//exec.Command("gnome-terminal", "-x", "go run ~/Documents/TTK4145/Exercise6/backup.go")
spawnBackup.Start()


toBackup := make(chan Message)
go Transmitter(toBackup)


for{

fmt.Printf("Current state: %d \n", counter.State)

toBackup <- Message{strconv.Itoa(counter.State)}
counter.State++ 
time.Sleep(1*time.Second)

}


}