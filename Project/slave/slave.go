package slave

import (
"../driver"
"../network/bcast"
"../network/localip"
"../orderManagement"
"../util"
"fmt"
//"math/rand"
"os/exec"
"time"
)


//Tested: works
func SendOrder(order util.Order, sendOrders chan util.Order, orderChan , otherOrderChan chan []util.Order, callback chan time.Time) {
	sendOrders <- order
	fmt.Println("Sent order")
	sendSuccess := false

	for i := 0 ; i < 3; i++{
		time.Sleep(1000 * time.Millisecond)
		select{
			case 	timestamp := <-callback: 
			if timestamp == order.AtTime{
			sendSuccess = true
			

		}
		default: 
	}
	}
	fmt.Println("callback received: ", sendSuccess)
	if !sendSuccess && !order.Completed {
		success := orderManagement.AddOrder(orderChan, otherOrderChan, order.FromButton.Floor, order.FromButton.TypeOfButton, order.ThisElevator, order.AtTime)
		if success == 1 {
			driver.SetButtonLamp(order.FromButton.Floor, order.FromButton.TypeOfButton, 1)
			fmt.Println("doing the order myself")

		}
	}
}
	
//must check if light is lit if another elevator is taking the order
func ListenRemoteOrders(listenForOrders chan util.Order, orderChan, otherOrderChan chan []util.Order) {

	for {

		order := <-listenForOrders
		if order.ThisElevator.IP == thisElevator.IP { //the elevator should complete the order itself

			if !order.Completed{ 

				success := orderManagement.AddOrder(orderChan, otherOrderChan, order.FromButton.Floor, order.FromButton.TypeOfButton, order.ThisElevator, order.AtTime)
				if success == 1 {

					driver.SetButtonLamp(order.FromButton.Floor, order.FromButton.TypeOfButton, 1)
				}
			}

		}else{ // another elevator will complete the order
			otherOrders := <- otherOrderChan

			if !order.Completed{ 
				otherOrders = append(otherOrders, order)
				driver.SetButtonLamp(order.FromButton.Floor, order.FromButton.TypeOfButton, 1)

			}else{

				otherOrders = orderManagement.RemoveOrder(order, otherOrders)
				driver.SetButtonLamp(order.FromButton.Floor, order.FromButton.TypeOfButton, 0)

			}

			otherOrderChan <- otherOrders


		}

	}
}

func ListenLocalOrders(sendOrders chan util.Order, orderChan, otherOrderChan chan []util.Order, callback chan time.Time) {

	var buttons [util.Nfloors][util.Nbuttons]int
	for {

		var recent [util.Nfloors][util.Nbuttons]int
		buttons = driver.ListenForButtons()
		changed, floor, button := CompareMatrix(buttons, recent)

		if changed {
				order := util.Order{thisElevator, util.Button{floor, button}, time.Now(), false}
				fmt.Println("sending order")
				go SendOrder(order, sendOrders, orderChan, otherOrderChan, callback)
				time.Sleep(700*time.Millisecond)
			}
			changed = false
		}
		time.Sleep(100 * time.Millisecond)
	}

//tested: works
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


//tested: works
func ExecuteOrder(sendOrders chan util.Order, orderChan , otherOrderChan chan []util.Order, callback chan time.Time) {

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
				driver.SetButtonLamp(currentOrder.FromButton.Floor, currentOrder.FromButton.TypeOfButton, 0)
			} else {
				goToFloor(currentOrder, currentFloor)
			}
			orderSlice := <-orderChan
			orderSlice = orderManagement.RemoveOrder(currentOrder, orderSlice)
			orderChan <- orderSlice
			currentOrder.Completed = true
			SendOrder(currentOrder, sendOrders, orderChan , otherOrderChan, callback)
		} else {

			orderChan <- orderSlice
		}
		time.Sleep(util.DoorOpenTime)
	}
}

func CompareMatrix(newMatrix, oldMatrix [util.Nfloors][util.Nbuttons]int) (changed bool, row, column int) {
	for i := 0; i < util.Nfloors; i++ {
		for j := 0; j < util.Nbuttons; j++ {
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
var thisElevator = util.Elevator{IP, 0, 2}

func SlaveLoop(isBackup bool) {

	orderChan := make(chan []util.Order, 1)
	orderSlice := make([]util.Order, 0)
	otherOrderChan := make(chan []util.Order, 1)
	otherOrders := make([]util.Order, 0)
	orderChan <- orderSlice
	otherOrderChan <- otherOrders


	listenForOrders := make(chan util.Order)
	sendOrders := make(chan util.Order)
	callback := make(chan time.Time)

	firstTry := true

	for {
		if isBackup && firstTry {
			firstTry = false
			orderChanBackup := make(chan []util.Order, 1)
			stateChanBackup := make(chan util.Elevator, 1)

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
			}()

			//checking that the slave is alive
			go func() {
				for {
					if isBackup {

						thisElevator = <-stateChanBackup
						fmt.Println("Elevator at floor", thisElevator.LastFloor)
						tmr.Reset(3 * time.Second)
					} else {
						return
					}
				}
			}()

			go func() {
				for {
					if driver.GetCurrentFloor() == 3 && thisElevator.ElevDirection == 0 || driver.GetCurrentFloor() == 0 && thisElevator.ElevDirection == 1 {
						driver.SteerElevator(2)
					}
				}
			}()

			//updating the orderSlice backup
			go func() {
				for isBackup {
					if len(orderChanBackup) == cap(orderChanBackup) {
						orderSlice = <-orderChanBackup
					}
				}
			}()

		}

		if !isBackup && firstTry {
			myIP, _ := localip.LocalIP()
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
			stateChanMaster := make(chan util.Elevator,1) 

			go bcast.Transmitter("255.255.255.255",20008, sendOrders, stateChanMaster)
			go bcast.Receiver(20009, listenForOrders, callback)
			go bcast.Transmitter( myIP,20010, newOrderChanBackup, newStateChanBackup)

			go ListenLocalOrders(sendOrders, orderChan, otherOrderChan, callback)
			go ExecuteOrder(sendOrders, orderChan, otherOrderChan,  callback)
			go ListenRemoteOrders(listenForOrders, orderChan, otherOrderChan)

			//notifying I'm alive
			//updating orderSlice backups
			go func() {
				for {
					select {
					case <-newStateChanBackup:
					case <- stateChanMaster:
							default:
					}

					newStateChanBackup <- thisElevator
					newOrderChanBackup <- orderSlice
					stateChanMaster <- thisElevator
					orderSlice = <-orderChan
					orderChan <- orderSlice
				//	fmt.Println(orderSlice)

					time.Sleep(1000 * time.Millisecond)

				}
			}()
			spawnBackup := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run /home/student/Documents/Group55/TTK4145/Project/main.go -startSlaveBackup")
			spawnBackup.Start()

		}
		time.Sleep(1 * time.Second)
	}
}
