package master

import (
	"../network/bcast"
	"../network/localip"
	"../util"
	"os/exec"
	//"fmt"
	//"../orderManagement"
)

var slaves [util.Nslaves]util.Elevator
var orders [util.Nslaves][]util.Order
var slaveIPs = [util.Nslaves]string{"129.241.187.161", "129.241.187.156","255.255.255.255"}
var slaveAlive [util.Nslaves]bool

func InitSlave(IP string) {
	for i := 0; i < len(slaveIPs); i++ {
		spawnSlave := exec.Command("gnome-terminal", "-x", "sh", "-c", "nohup ssh student@", IP, "go run /home/student/Documents/TTK4145/Exercise6/backup.go")
		spawnSlave.Start()
	}
}

func sendOrder(order util.Order, sendOrders chan util.Order) {
	sendOrders <- order
	//TODO: callback functionality
}

func distributeOrder(listenForOrders chan util.Order, sendOrders chan util.Order) {
	for {
		order := <-listenForOrders
		fmt.Println("Master Received order", order)
		sendTo := orderManagement.FindSuitableElevator(slaves, order)
		go sendOrder(order, sendOrders)
	}
}

func Master(isBackup bool) {
	var orderChannels [util.Nslaves]chan util.Order
	var statusChannels [util.Nslaves]chan util.Elevator
	var orderSliceChannel [util.Nslaves]chan []util.Order
	var sendOrderChannels [util.Nslaves]chan util.Order
	var slaveOrderSlices [util.Nslaves][]util.Order
	
	for j := 0; j < util.Nslaves; j++ {
		orderChannels[j] = make(chan util.Order)
		statusChannels[j] = make(chan util.Elevator)
		orderSliceChannel[j] = make(chan []util.Order)
		sendOrderChannels[j] = make(chan util.Order)
		slaveOrderSlices[j] = make([]util.Order, 0)
	}
	firstTry := true

	if isBackup && firstTry {
		firstTry = false
		for k := 0; k < util.Nslaves; k++ {
			go bcast.Receiver(20011, orderSliceChannel[k], statusChannels[k])
		}
		go func() {
			for c := 0; ; c++ {
				if slaveAlive[c] {
					slaves[c] = <-statusChannels[c]
					orders[c] = <-orderSliceChannel[c]
				}
				if c == util.Nslaves-1 {
					c = 0
				}
			}
		}()
	}

	if !isBackup && firstTry {
		myIP, _ := localip.LocalIP()

		//this should only be done once, right ? So first try should not be set to true when a backup becomes master..
		for i := 0; i < util.Nslaves; i++ {
			if myIP != slaveIPs[i] {
				InitSlave(slaveIPs[i])
			}
		}
		//start backup master on remote pc, take first in list that is not itself



		for c := 0; c < util.Nslaves; c++ {
			go bcast.Transmitter(20009, sendOrderChannels[c])
			go bcast.Receiver(20009, listenForOrders, listenForSlaves, listenForOrderSlice)
			//should all these channels have different names??
		}
		
		go distributeOrder(listenForOrders, sendOrders)

		//updating info from slave
		go func() {
			for {
				orderSliceSlave = <-listenForOrderSlice
			}
		}()

	}

}
