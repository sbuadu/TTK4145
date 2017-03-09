package main

import (
	"fmt"
	//"net"
	//"log"
	//"os/exec"
	//"strconv"
	"../network/network/bcast"
	"time"
)

type Counter struct {
	State int
}

type Message struct {
	Data int
}

func main() {
	var otherIP = "129.241.187.161"
	fmt.Print("Let's count!\n\n")
	counter := Counter{0}

	//spawnBackup := exec.Command("gnome-terminal", "-x", "sh", "-c", "ssh student@129.241.187.156 && echo "Sanntid15" && run /home/student/Documents/TTK4145/Exercise6/backup.go")

	//exec.Command("gnome-terminal", "-x", "go run ~/Documents/TTK4145/Exercise6/backup.go")
	//spawnBackup.Start()

	toBackup := make(chan Message, 1)
	port := 20009
	go bcast.Transmitter(otherIP, port, toBackup)

	for {

		fmt.Printf("Current state: %d \n", counter.State)
		msg := Message{counter.State}
		toBackup <- msg
		counter.State++
		time.Sleep(1 * time.Second)

	}

}
