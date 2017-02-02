package main
import(
	"fmt"
	//"time"
	"log"
	"net"
	"strconv"
	 //"os/exec"
)

type Counter struct{
	State int
}

type Message struct{
	Data string
}


func Receiver(fromMaster chan Message){
for{
buf := make([]byte, 1024)
	
	server, _ := net.ResolveUDPAddr("udp",":20099")
conn, err := net.ListenUDP("udp", server)
if err != nil {
	log.Fatal(err)
}

	n, _, _ := conn.ReadFromUDP(buf)
	fromMaster <- Message{string(buf[:n])}
	
	conn.Close()
	}
 }


func main(){

fmt.Print("\n\n The backup is running \n\n")
fromMaster := make(chan Message)
go Receiver(fromMaster)

msg := <- fromMaster
var masterCounter Counter



masterCounter.State, _ = strconv.Atoi(msg.Data) 

fmt.Print("the received state is %d\n", masterCounter.State)



}