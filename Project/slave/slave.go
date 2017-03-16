package slave

import (
	"../driver"
	"../network/bcast"
	"../network/localip"
	"../orderManagement"
	"../util"
	"fmt"
	"os/exec"
	"time"
)

/*
Module: slave

date: 14.03.17

This module handles the operation of the individual elevator.
*/

func sendOrder(order util.Order, sendOrders chan util.Order, orderChan, otherOrderChan chan []util.Order, callback chan time.Time) {
	sendOrders <- order
	sendSuccess := false

	for i := 0; i < 3; i++ {
		time.Sleep(500 * time.Millisecond)
		select {
		case timestamp := <-callback:
			if timestamp == order.AtTime {
				sendSuccess = true
				return
			}
		default:

		}
	}

	if !sendSuccess && !order.Completed {
		success := orderManagement.AddOrder(orderChan, otherOrderChan, order.FromButton.Floor, order.FromButton.TypeOfButton, order.ThisElevator, order.AtTime)
		if success == 1 {
			driver.SetButtonLamp(order.FromButton.Floor, order.FromButton.TypeOfButton, 1)

		} else {
		}
	}
}

//TODO: Must check if light is lit if another elevator is taking the order
//TODO: that an order received by one elevator is conducted by another if assigned to do so..
func listenRemoteOrders(listenForOrders chan util.Order, orderChan, otherOrderChan chan []util.Order) {

	for {
		select {

		case order := <-listenForOrders:

			if order.ThisElevator.IP == IP { //the elevator should complete the order itself

				if !order.Completed {
					success := orderManagement.AddOrder(orderChan, otherOrderChan, order.FromButton.Floor, order.FromButton.TypeOfButton, order.ThisElevator, order.AtTime)
					if success == 1 {
						driver.SetButtonLamp(order.FromButton.Floor, order.FromButton.TypeOfButton, 1)
					}
				}

			} else { // another elevator will complete the order

				if order.FromButton.TypeOfButton != 2 {
					otherOrders := <-otherOrderChan

					if !order.Completed {

						otherOrders = append(otherOrders, order)
						driver.SetButtonLamp(order.FromButton.Floor, order.FromButton.TypeOfButton, 1)

					} else {

						otherOrders = orderManagement.RemoveOrder(order, otherOrders)
						driver.SetButtonLamp(order.FromButton.Floor, order.FromButton.TypeOfButton, 0)

					}

					otherOrderChan <- otherOrders

				}

			}
		default:
		}
		time.Sleep(500 * time.Millisecond)

	}
}

func listenLocalOrders(sendOrders chan util.Order, orderChan, otherOrderChan chan []util.Order, callback chan time.Time, thisElevatorChan chan util.Elevator) {

	var buttons [util.Nfloors][util.Nbuttons]int
	for {

		var recent [util.Nfloors][util.Nbuttons]int
		buttons = driver.ListenForButtons()
		changed, floor, button := CompareMatrix(buttons, recent)

		if changed {
			thisElevatorTmp := <-thisElevatorChan
			thisElevatorChan <- thisElevatorTmp
			order := util.Order{thisElevatorTmp, util.Button{floor, button}, time.Now(), false}
			go sendOrder(order, sendOrders, orderChan, otherOrderChan, callback)
			time.Sleep(700 * time.Millisecond)
		}
		changed = false
		time.Sleep(100 * time.Millisecond)
	}

}

func goToFloor(order util.Order, currentFloor int, thisElevatorChan chan util.Elevator) {
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
	thisElevatorTmp := <-thisElevatorChan
	thisElevatorTmp.ElevDirection = elevDir
	thisElevatorChan <- thisElevatorTmp
	for currentFloor != orderFloor {
		floor := driver.GetCurrentFloor()
		if floor != -1 {
			currentFloor = floor
			thisElevatorTmp := <-thisElevatorChan
			thisElevatorTmp.LastFloor = currentFloor
			thisElevatorChan <- thisElevatorTmp
			driver.SetFloorIndicator(currentFloor)
		}
	}
	driver.SteerElevator(2)
	thisElevatorTmp = <-thisElevatorChan
	thisElevatorTmp.ElevDirection = 2
	thisElevatorChan <- thisElevatorTmp
	driver.SetButtonLamp(orderFloor, order.FromButton.TypeOfButton, 0)
	driver.SetDoorLamp(1)
}

func executeOrder(sendOrders chan util.Order, orderChan, otherOrderChan chan []util.Order, callback chan time.Time, thisElevatorChan chan util.Elevator) {

	currentFloor := driver.GetCurrentFloor()
	for currentFloor == -1 {
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
				goToFloor(currentOrder, currentFloor, thisElevatorChan)
			}
			orderSlice := <-orderChan
			orderSlice = orderManagement.RemoveOrder(currentOrder, orderSlice)
			orderChan <- orderSlice
			currentOrder.Completed = true
			sendOrder(currentOrder, sendOrders, orderChan, otherOrderChan, callback)
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

func SlaveLoop(isBackup bool) {

	orderSlice := make([]util.Order, 0)
	otherOrders := make([]util.Order, 0)

	//local process channels
	orderChan := make(chan []util.Order, 1)
	otherOrderChan := make(chan []util.Order, 1)
	thisElevatorChan := make(chan util.Elevator, 1)

	//channels for communication with master
	listenForOrders := make(chan util.Order, 1)
	sendOrders := make(chan util.Order, 1)
	callback := make(chan time.Time, 1)

	var IP, _ = localip.LocalIP()
	var thisElevator = util.Elevator{IP, 0, 2}
	orderChan <- orderSlice
	otherOrderChan <- otherOrders
	thisElevatorChan <- thisElevator
	firstTry := true

	for {
		if isBackup && firstTry {
			firstTry = false

			//channels for communicating with slave master
			orderChanBackup := make(chan []util.Order, 1)
			stateChanBackup := make(chan util.Elevator, 1)

			go bcast.Receiver(20010, orderChanBackup, stateChanBackup)

			tmr := time.NewTimer(5 * time.Second)

			//listening for timer laps and taking over operation
			go func() {

				<-tmr.C
				fmt.Println("timer lapsed")
				isBackup = false
				firstTry = true
				select {
				case <-orderChan:
				default:
				}
				select {
				case <-otherOrderChan:
				default:
				}
				select {
				case <-thisElevatorChan:
				default:
				}
				thisElevatorChan <- thisElevator
				orderChan <- orderSlice
				otherOrderChan <- otherOrders
				fmt.Println("Taking over as slave")
			}()

			//listening for updates from slave
			go func() {
				for {
					if isBackup {

						thisElevator = <-stateChanBackup
						tmr.Reset(3 * time.Second)
						tmpOrderSlice := <-orderChanBackup
						if len(tmpOrderSlice) != 0 {
							if tmpOrderSlice[0].ThisElevator.IP == thisElevator.IP {
								orderSlice = tmpOrderSlice
							} else {
								otherOrders = tmpOrderSlice
							}
						}
					} else {
						return
					}
				}
			}()

			//Guard
			go func() {
				for isBackup {
					thisElevatorTMP :=<-thisElevatorChan
					if driver.GetCurrentFloor() == 3 && thisElevatorTMP.ElevDirection == 0 || driver.GetCurrentFloor() == 0 && thisElevatorTMP.ElevDirection == 1 {
						driver.SteerElevator(2)
					}
					thisElevatorChan <- thisElevatorTMP
					time.Sleep(200 * time.Second)
				}
			}()
		}

		if !isBackup && firstTry {

			myIP, _ := localip.LocalIP()
			firstTry = false
			driver.InitElevator()

			orderSlice := <-orderChan

			for i := 0; i < len(orderSlice); i++ {
				driver.SetButtonLamp(orderSlice[i].FromButton.Floor, orderSlice[i].FromButton.TypeOfButton, 1)
			}
			otherOrders = <-otherOrderChan
			for i := 0; i < len(otherOrders); i++ {
				driver.SetButtonLamp(otherOrders[i].FromButton.Floor, otherOrders[i].FromButton.TypeOfButton, 1)
			}
			orderChan <- orderSlice
			otherOrderChan <- otherOrders

			fmt.Println("I'm a slave now")

			//channels for communicating with slave backup
			newStateChanBackup := make(chan util.Elevator, 1)
			newOrderChanBackup := make(chan []util.Order, 1)
			stateChanMaster := make(chan util.Elevator, 1)

			go bcast.Transmitter("255.255.255.255", 20008, sendOrders, stateChanMaster)
			go bcast.Receiver(20009, listenForOrders, callback)
			go bcast.Transmitter(myIP, 20010, newOrderChanBackup, newStateChanBackup)

			go listenLocalOrders(sendOrders, orderChan, otherOrderChan, callback, thisElevatorChan)
			go executeOrder(sendOrders, orderChan, otherOrderChan, callback, thisElevatorChan)
			go listenRemoteOrders(listenForOrders, orderChan, otherOrderChan)

			//notifying I'm alive
			//updating orderSlice backups
			go func() {
				for {
					select {
					case <-newStateChanBackup:
					default:
					}
					select {
					case <-stateChanMaster:
					default:
					}
					thisElevatorTmp := <-thisElevatorChan
					thisElevatorChan <- thisElevatorTmp
					newStateChanBackup <- thisElevatorTmp
					newOrderChanBackup <- orderSlice
					stateChanMaster <- thisElevatorTmp
					orderSlice = <-orderChan
					orderChan <- orderSlice

					time.Sleep(1000 * time.Millisecond)

				}
			}()

			//Checking for long-time incompleted orders
			go func() {
				for {

					select {
					case otherOrders = <-otherOrderChan:

						if len(otherOrders) > 0 {
							for i := 0; i < len(otherOrders); i++ {
								otherOrders = <-otherOrderChan
								if time.Since(otherOrders[i].AtTime) > time.Second*20 {
									otherOrders[i].ThisElevator = thisElevator
									otherOrders = orderManagement.RemoveOrder(otherOrders[i], otherOrders)
									otherOrderChan <- otherOrders

									orderManagement.AddOrder(orderChan, otherOrderChan, otherOrders[i].FromButton.Floor, otherOrders[i].FromButton.TypeOfButton, otherOrders[i].ThisElevator, otherOrders[i].AtTime)
									i -= 1
								} else {
									otherOrderChan <- otherOrders

								}

							}

						} else {
							otherOrderChan <- otherOrders

						}

					default:
					}

					time.Sleep(10 * time.Second)
				}
			}()

			spawnBackup := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run /home/student/Documents/Group55/TTK4145/Project/main.go -startSlaveBackup")
			spawnBackup.Start()

		}
		time.Sleep(1 * time.Second)
	}
}
