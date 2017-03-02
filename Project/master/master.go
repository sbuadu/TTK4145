package master

import (
	"../network/bcast"
	"../util"
	//"fmt"
	//"../orderManagement"
)

var slaves [util.Nslaves]util.Elevator
var orders [util.Nslaves][]util.Order
var slaveIPs [util.Nslaves]string
var slaveAlive [util.Nslaves]bool

func InitSlave(IP string) {
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
		order := <-listenForOrders
		//fmt.Println("Master Received order", order)
		//sendTo := orderManagement.FindSuitableElevator(slaves, order)
		//TODO: populate slavcelist for this to work
		go sendOrder(order, sendOrders)
	}
}

func Master(isBackup bool) {
	//receiving from slave
	listenForOrders := make(chan util.Order)
	listenForSlaves := make(chan util.Elevator)
	listenForOrderSlice := make(chan []util.Order, 1)

	//sending to slave
	sendOrders := make(chan util.Order)

	//variable saving states
	orderSliceSlave := []util.Order{}

	if isBackup {

	}

	if !isBackup {
		for i := 0; i < util.Nslaves; i++ {
			InitSlave(slaveIPs[i])
		}
		//start backup master on remote pc, take first in list that is not itself

		go bcast.Transmitter(20009, sendOrders)
		go bcast.Receiver(20009, listenForOrders, listenForSlaves, listenForOrderSlice)

		go distributeOrder(listenForOrders, sendOrders)

		//updating info from slave
		go func() {
			for {
				orderSliceSlave = <-listenForOrderSlice
			}
		}()

	}

}
