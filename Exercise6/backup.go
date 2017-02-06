package main
import(
	"../network/network/bcast"
	"fmt"
	"time"
	//"log"
	//"net"
	//"strconv"
	"os/exec"
	//"../conn"
)

type Counter struct{
	State int
}

type Message struct{
	Data int
}


 func takeOverAsMaster(counter Counter){
 	toBackup := make(chan Message, 1)

 	//staring a new backup 
	spawnBackup := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run backup.go")
	spawnBackup.Start()

	
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


var isBackup bool = true
func main(){

	
	tmr := time.NewTimer(3*time.Second)
	port := 20009
	fmt.Print("\n\nThe backup is running \n\n")
	fromMaster := make(chan Message)
	
	masterCounter := Counter{0}

	go bcast.Receiver(port, fromMaster)
	

	// if the timer runs out...
	go func(){
		<- tmr.C
		isBackup = false
		fmt.Print("Master seems to be dead, I'm taking over.. ")

				takeOverAsMaster(masterCounter)
	}()
	for {
	 	// receiving the current state of the master and printing it 
		if isBackup {
			msg := <-fromMaster
			masterCounter.State = msg.Data
			tmr.Reset(3*time.Second)
			fmt.Printf("the received state is %d\n", masterCounter.State)

	} 

}
}
