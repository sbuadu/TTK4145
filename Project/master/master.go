package master

import (
	"../network/bcast"
	"../orderManagement"
	"../util"
	"fmt"
	"time"
)

/*
Module: Master

date: 14.03.17


This module works as a supervisor of the elevators and their operation.
It is also resposible for assigning orders to individual elevators

*/

func sendOrder(order util.Order, sendOrdersChannel chan util.Order) {
	for i := 0; i < 3; i++ {
		sendOrdersChannel <- order
		time.Sleep(100 * time.Millisecond)
	}
}

func distributeOrder(order util.Order, sendOrdersChannel chan util.Order, orderChan chan [util.Nslaves][]util.Order, slaveAliveChan chan [util.Nslaves]bool, slavesChan chan [util.Nslaves]util.Elevator) {
	var workingSlaves = make([]util.Elevator, 0)

	for i := 0; i < util.Nslaves; i++ {
		slaveAlive := <-slaveAliveChan
		slaves := <-slavesChan
		slaveAliveChan <- slaveAlive
		slavesChan <- slaves
		if slaveAlive[i] {
			workingSlaves = append(workingSlaves, slaves[i])
		}
	}

	order.ThisElevator = orderManagement.FindSuitableElevator(order, workingSlaves)
	fmt.Print("This elevator: ", order.ThisElevator.IP)
	fmt.Println("  should go to floor:", order.FromButton.Floor)
	go sendOrder(order, sendOrdersChannel)

	//adding the order to the right backup slice
	for i := 0; i < util.Nslaves; i++ {
		if order.ThisElevator.IP == util.SlaveIPs[i] {
			orders := <-orderChan
			tempChan := make(chan []util.Order, 1)
			tempChan <- orders[i]
			dummyChan := make(chan []util.Order, 1)
			dummySlice := make([]util.Order, 0)
			dummyChan <- dummySlice

			orderManagement.AddOrder(tempChan, dummyChan, order.FromButton.Floor, order.FromButton.TypeOfButton, order.ThisElevator, order.AtTime)
			orders[i] = <-tempChan
			orderChan <- orders

		}

	}
}

func receiveOrder(orderChannel, sendOrdersChannel chan util.Order, orderChan chan [util.Nslaves][]util.Order, slaveAliveChan chan [util.Nslaves]bool, slavesChan chan [util.Nslaves]util.Elevator, callbackChannel chan time.Time) {
	for {
		order := <-orderChannel
		callbackChannel <- order.AtTime

		//removing the completed order from the backup slice
		if order.Completed {
			go sendOrder(order, sendOrdersChannel)
			for i := 0; i < util.Nslaves; i++ {
				if order.ThisElevator.IP == util.SlaveIPs[i] {
					orders := <-orderChan
					orders[i] = orderManagement.RemoveOrder(order, orders[i])
					orderChan <- orders
				}
			}

		} else if order.FromButton.TypeOfButton == 2 { //sending back order if it is an internal order
			go sendOrder(order, sendOrdersChannel)

		} else {

			distributeOrder(order, sendOrdersChannel, orderChan, slaveAliveChan, slavesChan)

		}
	}
}

func MasterLoop(isBackup bool) {
	var slaves [util.Nslaves]util.Elevator
	var slaveAlive [util.Nslaves]bool
	var orders [util.Nslaves][]util.Order

	//channels for communication with slaves
	orderChannel := make(chan util.Order, 1)
	statusChannel := make(chan util.Elevator, 1)
	sendOrdersChannel := make(chan util.Order, 1)
	callbackChannel := make(chan time.Time, 1)

	//local process channels
	slavesChan := make(chan [util.Nslaves]util.Elevator, 1)
	orderChan := make(chan [util.Nslaves][]util.Order, 1)
	slaveAliveChan := make(chan [util.Nslaves]bool, 1)

	//channels for communication with master backup
	orderBackupChan := make(chan [util.Nslaves][]util.Order, 1)
	slavesBackupChan := make(chan [util.Nslaves]util.Elevator, 1)
	slaveAliveBackupChan := make(chan [util.Nslaves]bool, 1)

	slavesChan <- slaves
	orderChan <- orders
	slaveAliveChan <- slaveAlive

	for j := 0; j < util.Nslaves; j++ {
		orders[j] = make([]util.Order, 0)
	}
	firstTry := true

	//starting the Master loop
	for {

		if isBackup && firstTry {
			fmt.Println("I am a master backup")
			firstTry = false
			tmr := time.NewTimer(5 * time.Second)

			//channels for communication with master master
			ordersFromMaster := make(chan [util.Nslaves][]util.Order, 1)
			statusFromMaster := make(chan [util.Nslaves]util.Elevator, 1)
			slaveAliveFromMaster := make(chan [util.Nslaves]bool, 1)

			go bcast.Receiver(20011, ordersFromMaster, statusFromMaster, slaveAliveFromMaster)

			//Receiving status update from Master
			go func() {
				for {
					//fmt.Println("receiving update from Master")
					orders = <-ordersFromMaster
					slaves = <-statusFromMaster
					tmr.Reset(5 * time.Second)
				}
			}()

			//Receiving changes in elevator status (dead or alive)
			go func() {
				slaveAlive = <-slaveAliveFromMaster
			}()

			//listening for timer laps and taking over operation
			go func() {
				<-tmr.C
				isBackup = false
				firstTry = true
				fmt.Println("Master is dead, dobby is a free elf!")
				select {
				case <-slavesChan:
				default:
				}
				select {
				case <-orderChan:
				default:
				}
				select {
				case <-slaveAliveChan:
				default:
				}
				slavesChan <- slaves
				orderChan <- orders
				slaveAliveChan <- slaveAlive

			}()
		}

		if !isBackup && firstTry {
			fmt.Println("I am the master")
			firstTry = false

			if slaves[0].IP == "" { //if slaves arent initialized, initialize
				for i := 0; i < util.Nslaves; i++ {
					slaves = <-slavesChan
					slaves[i] = util.Elevator{util.SlaveIPs[i], 0, 2}
					slavesChan <- slaves
					slaveAlive = <-slaveAliveChan
					slaveAlive[i] = true
					slaveAliveChan <- slaveAlive
					fmt.Println("I belive there are slaves on ", slaves[i].IP)
				}
				slaveAliveBackupChan <- slaveAlive
			}

			//Settig up timers for slaves
			var timers [util.Nslaves]*time.Timer
			for i := 0; i < util.Nslaves; i++ {
				timers[i] = time.NewTimer(20 * time.Second)
			}

			//Set up communication
			go bcast.Transmitter("255.255.255.255", 20009, sendOrdersChannel, callbackChannel)
			go bcast.Receiver(20008, orderChannel, statusChannel)
			go bcast.Transmitter("255.255.255.255", 20011, orderBackupChan, slavesBackupChan, slaveAliveBackupChan)

			go receiveOrder(orderChannel, sendOrdersChannel, orderChan, slaveAliveChan, slavesChan, callbackChannel)

			//sending updates to backup
			go func() {
				for {
					orders = <-orderChan
					orderChan <- orders
					slaves = <-slavesChan
					slavesChan <- slaves
					select {
					case <-orderBackupChan:
					case <-slavesBackupChan:
					default:
					}
					orderBackupChan <- orders
					slavesBackupChan <- slaves
					//fmt.Println("Sent update to backup")
					time.Sleep(1 * time.Second)
				}
			}()

			//Receiving status update from slave
			go func() {
				for {
					status := <-statusChannel
					for i := 0; i < util.Nslaves; i++ {
						if status.IP == util.SlaveIPs[i] {
							slaves = <-slavesChan
							slaves[i] = status
							slavesChan <- slaves
							//fmt.Println(status.IP, " Present")
							timers[i].Reset(5 * time.Second)
						}
					}
				}
			}()

			//must test if works nomatter which order the slave IPs are listed
			// checking for non-responsive slaves and working accordingly
			go func() {
				for {
					for j := 0; j < util.Nslaves; j++ {
						select {
						case <-timers[j].C:
							fmt.Println("Slave is dead. IP: ", util.SlaveIPs[j])
							slaveAlive = <-slaveAliveChan
							slaveAlive[j] = false
							slaveAliveChan <- slaveAlive
							slaveAliveBackupChan <- slaveAlive
							orders = <-orderChan
							orderChan <- orders
							fmt.Println("redistibuting dead slaves orders")
							for i := 0; i < len(orders[j]); i++ {

								if !(orders[j][i].FromButton.TypeOfButton == 2) {
									orders[j][i].Completed = true
									fmt.Println("Sending complete dead order")
									go sendOrder(orders[j][i], sendOrdersChannel)
									orders[j][i].Completed = false
									time.Sleep(2 * time.Second)
									distributeOrder(orders[j][i], sendOrdersChannel, orderChan, slaveAliveChan, slavesChan)

								}
							}
						default:
						}
					}
				}
			}()

		}
		time.Sleep(1 * time.Second)
	}

}
