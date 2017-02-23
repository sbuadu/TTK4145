package master
import (

	"../util"
	//"fmt"
)

	var slaves [3]util.Elevator

func InitSlave() {

}

func sendOrder(order util.Order, sendOrders chan util.Order) {
	sendOrders <- order
	//TODO: callback functionality
}

func HandleOrder(listenForOrders chan util.Order, sendOrders chan util.Order) {
	for {
		order := <- listenForOrders
		//fmt.Println("Master Received order", order)
		//sendTo := orderManagement.findSuitableElevator(slaves, order)
		//TODO: populate slavcelist for this to work
		go sendOrder(order, sendOrders)
	}
}


