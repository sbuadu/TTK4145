package slave

import (
	"../driver"
	"../network/bcast"
	"../network/localip"
	"../orderManagement"
	"../util"
	"fmt"
	"time"
	"math/rand"
)

//TODO: make process pair functionality

func SendOrder(order util.Order, sendOrders chan util.Order, callback chan time.Time) {
		sendOrders <- order
		//fmt.Println("Slave Sent order", order)

		//TODO: callback functionality
}

func ListenRemoteOrders(listenForOrders chan util.Order, orderChan chan []util.Order) {
	//TODO: callback
	for {
		order := <- listenForOrders
		success := orderManagement.AddOrder(orderChan, order.FromButton.Floor, order.FromButton.TypeOfButton, order.ThisElevator, order.AtTime)
			if success == 1 {
				driver.SetButtonLamp(order.FromButton.Floor, order.FromButton.TypeOfButton, 1)

			}
		}

}

func ListenLocalOrders(callback chan time.Time, sendOrders chan util.Order, orderChan chan []util.Order) {

	//TODO: check if button is already on
	var buttons [4][3]int
	for {
		//fmt.Println(".localsuccess.")
		var recent [4][3]int
		buttons = driver.ListenForButtons()
		changed, floor, button := CompareMatrix(buttons, recent)

		if changed {
			if button == 0 || button == 1 {
				order := util.Order{thisElevator, util.Button{floor, button}, time.Now()} 
				go SendOrder(order, sendOrders, callback)
			} else {
				success := orderManagement.AddOrder(orderChan, floor, button, thisElevator, time.Now())
				if success == 1 {
					//TODO: Move this to after order is appended to orders
					driver.SetButtonLamp(floor, button, 1)
				}
			}
			changed = false
		}
		time.Sleep(10 * time.Millisecond)
	}
}
func goToFloor(order util.Order, currentFloor int) {
	orderFloor := order.FromButton.Floor
	higher := currentFloor < orderFloor
	var elevDir util.Direction
	if higher {
		elevDir = 0
	} else if !higher {
		elevDir = 1
	}
	driver.SetDoorLamp(0)
	driver.SteerElevator(elevDir)
	thisElevator.ElevDirection = elevDir
	for currentFloor != orderFloor {
		floor := driver.GetCurrentFloor()
		if floor != -1 {
			currentFloor = floor
			thisElevator.LastFloor = currentFloor
			driver.SetFloorIndicator(currentFloor)
		}
	}
	driver.SteerElevator(2)
	thisElevator.ElevDirection = 2
	driver.SetButtonLamp(orderFloor, order.FromButton.TypeOfButton, 0)
	driver.SetDoorLamp(1)
}

func ExecuteOrder(orderChan chan []util.Order) {
	currentFloor := driver.GetCurrentFloor()
	if currentFloor == -1 {
		currentFloor = 0
	}
	for {
		orderSlice := <-orderChan
			if len(orderSlice) > 0 {
			currentOrder := orderSlice[0]
			orderChan <- orderSlice
			floor := currentOrder.FromButton.Floor
			currentFloor = driver.GetCurrentFloor()
			if currentFloor == floor {
				driver.SetDoorLamp(1)
				orderSlice := <-orderChan
				orderSlice = orderManagement.RemoveOrder(currentOrder, orderSlice)
				orderChan <- orderSlice
				driver.SetButtonLamp(currentOrder.FromButton.Floor, currentOrder.FromButton.TypeOfButton, 0)
			} else {
				goToFloor(currentOrder, currentFloor)
				orderSlice := <-orderChan
				orderSlice = orderManagement.RemoveOrder(currentOrder, orderSlice)
				orderChan <- orderSlice
			}
		} else {
			orderChan <- orderSlice
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func CompareMatrix(newMatrix, oldMatrix [4][3]int) (changed bool, row, column int) {
	for i := 0; i < 4; i++ {
		for j := 0; j < 3; j++ {
			if newMatrix[i][j] != oldMatrix[i][j] {
				changed = true
				row = i
				column = j
				return changed, row, column
			}
		}
	}
	return false, 0, 0
}
var IP,_ = localip.LocalIP()
var thisElevator = util.Elevator{rand.Intn(100),IP,0,2}
func Slave() {
	
	var isBackup = false
	driver.InitElevator()
	orderChan := make(chan []util.Order, 100)
	orderChan <- []util.Order{}
	listenForOrders := make(chan util.Order)
	sendOrders := make(chan util.Order)
	callback := make(chan time.Time)
	if !isBackup {
		go bcast.Transmitter(20009, sendOrders)
		go bcast.Receiver(20010, listenForOrders)
		go bcast.Receiver(20009, callback)
		go ListenLocalOrders(callback, sendOrders, orderChan)
		go ExecuteOrder(orderChan)
		go ListenRemoteOrders(listenForOrders, orderChan)
	} else if isBackup {
		go ListenLocalOrders(callback, sendOrders, orderChan)
		go bcast.Receiver(20009, listenForOrders)
	}
}


func Test() {


	driver.InitElevator()
	//orderSlice := []util.Order{util.Order{util.Elevator{1, "IP", 0, 0}, util.Button{0, 0}, time.Now()}}
	orderChan := make(chan []util.Order, 100)
	orderChan <- []util.Order{}
	//listenForOrders := make(chan util.Order)
	sendOrders := make(chan util.Order)
	callback := make(chan time.Time)

	go ListenLocalOrders(callback, sendOrders, orderChan)
	orderChan <- []util.Order{}
	for {

		orderSlice := <-orderChan
		for i := 0; i < len(orderSlice); i++ {
			fmt.Println(orderSlice[i].FromButton.Floor)
			//fmt.Println(len(orderSlice))
		}
		orderChan <- orderSlice
	}

}
