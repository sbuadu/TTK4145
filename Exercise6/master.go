package main

import (
	"fmt"
	//"net"
	//"log"
	"os/exec"
	//"strconv"
	"time"
	"../network/network/bcast"
)



type Counter struct{
	State int
}

type Message struct{
	Data int
}


func main(){

	fmt.Print("Let's count!\n\n")
	counter := Counter{0}

	spawnBackup := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run backup.go")

	//exec.Command("gnome-terminal", "-x", "go run ~/Documents/TTK4145/Exercise6/backup.go")
	spawnBackup.Start()


	toBackup := make(chan Message, 1)
	port := 20009
	go bcast.Transmitter(port,toBackup)


	for {

		fmt.Printf("Current state: %d \n", counter.State)
		msg := Message{counter.State}
		toBackup <- msg
		counter.State++
		time.Sleep(1*time.Second)

	}



}
