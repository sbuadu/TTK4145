package master
import (

	"../util"
	//"fmt"
)

	var slaves [util.Nslaves]util.Elevator
	var orders [util.Nslaves][]util.Orders
	var slaveIPs [util.Nslaves]util.Elevator.IP
	var slaveAlive [util.Nslaves]bool
func InitSlave(IP util.Elevator.IP) {
	//TODO: This
	//start elevator from command line
	//add to list of elevators
	//update slaveAlive
}

func sendOrder(order util.Order, sendOrders chan util.Order) {
	sendOrders <- order
	//TODO: callback functionality
}

func distributeOrder(listenForOrders chan util.Order, sendOrders chan util.Order) {
	for {
		order := <- listenForOrders
		//fmt.Println("Master Received order", order)
		sendTo := orderManagement.findSuitableElevator(slaves, order)
		//TODO: populate slavcelist for this to work
		go sendOrder(order, sendOrders)
	}
}

func Master(isBackup bool) {

	listenForOrders := make(chan util.Order)
	sendOrders := make(chan util.Order)
	listenForSlaves := make(chan util.Elevator)
	listenForOrderSlice := make(chan []util.Order)
	i isBackup {
	}
	for i:=0;i<util.Nslaves;i++ {
		InitSlave(slaveIPs[i])
	}
	//start backup master on remote pc, take first in list that is not itself
        go bcast.Transmitter(20010, sendOrders)
        go bcast.Receiver(20009, listenForOrders,listenForSlaves)
	go distributeOrder(listenForOrders,sendOrders)
}
