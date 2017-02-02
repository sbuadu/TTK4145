package main
import(
	"../network/network/bcast"
	"fmt"
	//"time"
	//"log"
	//"net"
	//"strconv"
	//"os/exec"
	//"../conn"
)

type Counter struct{
	State int
}

type Message struct{
	Data int
}

/*
func Receiver(port int, fromMaster chan Message){
	buf := make([]byte, 1024)
	server, err := net.ResolveUDPAddr("udp",fmt.Sprintf(":%d",port))
	//con := conn.DialBroadcastUDP(port)
	if err != nil {
		fmt.Println(err)
	}
	con, err1 := net.ListenUDP("udp",server)
	fmt.Println(err1)
	for {
		n, _, err2 := con.ReadFromUDP(buf)
		fmt.Println(err2)
		fromMaster <- Message{string(buf[:n])}
	}
	defer con.Close()
 }
*/

func main(){

	port := 20009
	fmt.Print("\n\nThe backup is running \n\n")
	fromMaster := make(chan Message)
	go bcast.Receiver(port, fromMaster)
	masterCounter := Counter{0}

	for {
		msg := <-fromMaster
		masterCounter.State = msg.Data

		fmt.Printf("the received state is %d\n", masterCounter.State)

	}

}
