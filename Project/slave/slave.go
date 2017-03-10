package slave

import (
"../driver"
"../network/bcast"
"../network/localip"
"../orderManagement"
"../util"
"fmt"
"math/rand"
"os/exec"
"time"
)


func SendOrder(order util.Order, sendOrders chan util.Order, orderChan chan []util.Order) {
	if len(sendOrders) < cap(sendOrders) {
		sendOrders <- order
	} else {
		success := orderManagement.AddOrder(orderChan, order.FromButton.Floor, order.FromButton.TypeOfButton, order.ThisElevator, order.AtTime)
		if success == 1 {

			driver.SetButtonLamp(order.FromButton.Floor, order.FromButton.TypeOfButton, 1)
		}

	}
}

func ListenRemoteOrders(listenForOrders chan util.Order, orderChan chan []util.Order) {
	//TODO: callback

	//should have a boolean here to check if this is to be added inn orderSlice, or if button should just be lit as in another elevator will complete the order.. 

	for {
		order := <-listenForOrders
		success := orderManagement.AddOrder(orderChan, order.FromButton.Floor, order.FromButton.TypeOfButton, order.ThisElevator, order.AtTime)
		if success == 1 {

			driver.SetButtonLamp(order.FromButton.Floor, order.FromButton.TypeOfButton, 1)

		}
	}

}

func ListenLocalOrders(callback chan time.Time, sendOrders chan util.Order, orderChan chan []util.Order) {


	var buttons [util.Nfloors][util.Nslaves]int
	for {
		
		var recent [util.Nfloors][util.Nslaves]int
		buttons = driver.ListenForButtons()
		changed, floor, button := CompareMatrix(buttons, recent)
	
		if changed {
			if button == 0 || button == 1 { //external order
				order := util.Order{thisElevator, util.Button{floor, button}, time.Now()}
				go SendOrder(order, sendOrders, orderChan)
			} else { //internal order
				success := orderManagement.AddOrder(orderChan, floor, button, thisElevator, time.Now())
				if success == 1 {

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
		time.Sleep(util.DoorOpenTime)
	}
}

func CompareMatrix(newMatrix, oldMatrix [util.Nfloors][util.Nslaves]int) (changed bool, row, column int) {
	for i := 0; i < util.Nfloors; i++ {
		for j := 0; j < util.Nslaves; j++ {
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

var IP, _ = localip.LocalIP()
var thisElevator = util.Elevator{rand.Intn(100), IP, 0, 2}

func Slave(isBackup bool) {

	orderChan := make(chan []util.Order, 1)
	orderSlice := make([]util.Order, 0)
	orderChan <- orderSlice

	listenForOrders := make(chan util.Order)
	sendOrders := make(chan util.Order)
	callback := make(chan time.Time)

	orderChanMaster := make(chan []util.Order, 1)  //used to send updates on order slice to master
	stateChanMaster := make(chan util.Elevator, 1) // used to send updates on the elevators state to master
	firstTry := true

	for {
		if isBackup && firstTry {
			firstTry = false
			orderChanBackup := make(chan []util.Order, 1)  //used to send updates on order slice to backup
			stateChanBackup := make(chan util.Elevator, 1) // used to send updates on the elevators state to backup

			go bcast.Receiver(20010, orderChanBackup, stateChanBackup)
			tmr := time.NewTimer(5 * time.Second)

			go func() {
				<-tmr.C
				isBackup = false
				firstTry = true
				select {
				case <-orderChan:
				default:
				}
				orderChan <- orderSlice
				fmt.Println("Taking over as slave")
				spawnBackup := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run /home/student/Documents/Group55/TTK4145/Project/main.go -startSlaveBackup")
				spawnBackup.Start()

				}()

			//checking that the slave is alive
				go func() {
					for {
						if isBackup {

							thisElevator = <-stateChanBackup

							tmr.Reset(5 * time.Second)
						} else {
							return
						}
					}
					}()

					go func() {
						if driver.GetCurrentFloor() == 3 && thisElevator.ElevDirection == 0 || driver.GetCurrentFloor() == 0 && thisElevator.ElevDirection == 1 {
							driver.SteerElevator(2)
						}
						}()

			//updating the orderSlice backup
						go func() {
							for {
								if len(orderChanBackup) == cap(orderChanBackup) && isBackup {
									orderSlice = <-orderChanBackup
								}

							}
							}()

		}

		if !isBackup && firstTry {
			firstTry = false
			driver.InitElevator()
			orderSlice := <-orderChan
			orderChan <- orderSlice
			for i := 0; i < len(orderSlice); i++ {
				driver.SetButtonLamp(orderSlice[i].FromButton.Floor, orderSlice[i].FromButton.TypeOfButton, 1)
			}

			fmt.Println("I'm a slave now")
			newStateChanBackup := make(chan util.Elevator, 1)
			newOrderChanBackup := make(chan []util.Order, 1)

			go bcast.Transmitter(20009, sendOrders, orderChanMaster, stateChanMaster)
			go bcast.Receiver(20009, listenForOrders, callback)
			go bcast.Transmitter(20010, newOrderChanBackup, newStateChanBackup)

			go ListenLocalOrders(callback, sendOrders, orderChan)
			go ExecuteOrder(orderChan)
			go ListenRemoteOrders(listenForOrders, orderChan)

			//notifying I'm alive
			//updating orderSlice backups
			go func() {
				for {
					select {
						case <-newStateChanBackup:
							default:
					}
					newStateChanBackup <- thisElevator

					orderSlice = <-orderChan
					newOrderChanBackup <- orderSlice
					orderChan <- orderSlice
					//stateChanMaster <- thisElevator

					time.Sleep(100 * time.Millisecond)

				}
			}()

		}
		time.Sleep(1 * time.Second)
	}
}
