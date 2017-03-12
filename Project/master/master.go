package master

import (
"../network/bcast"
"../network/localip"
"../util"
"os/exec"
"time"
"fmt"
"../orderManagement"
)

var slaves [util.Nslaves]util.Elevator
var slaveIPs = [util.Nslaves]string{"129.241.187.161", "129.241.187.156","255.255.255.255"}
var slaveAlive [util.Nslaves]bool

func InitSlave(IP string) {
	spawnSlave := exec.Command("gnome-terminal", "-x", "sh", "-c", "sshpass -p 'Sanntid15' ssh student@", IP, "go run /home/student/Documents/Group55/TTK4145/Project/main.go -startSlave")
	spawnSlave.Start()
}

func sendOrder(order util.Order, sendOrders chan util.Order) {
	for i :=0; i<3; i++{
		sendOrders <- order
		time.Sleep(1*time.Second)
	}
}

func distributeOrder(listenForOrders, sendOrders chan util.Order, orderChan chan [util.Nslaves][]util.Order) {
	for {
		order := <-listenForOrders
		fmt.Println("Master Received order", order)
		if order.Completed {
			go sendOrder(order, sendOrders)
			for i :=0;i<util.Nslaves;i++{
				if order.ThisElevator.IP == slaveIPs[i]{
					orders := <-orderChan
					orders[i] = orderManagement.RemoveOrder(order,orders[i])
					orderChan <- orders
				}
			}
		} else {
			var workingSlaves = make([]util.Elevator,0)
			for i := 0; i<util.Nslaves; i++{
				if slaveAlive[i]{
					workingSlaves = append(workingSlaves, slaves[i])
				}
			}	
			order.ThisElevator = orderManagement.FindSuitableElevator(workingSlaves, order)
			go sendOrder(order, sendOrders)
		}
	}
}

func Master(isBackup bool) {
	var orderChannel = make(chan util.Order) //orders coming from elevators
	var statusChannel = make(chan util.Elevator) //Status from elevators
	var orderSliceChannel = make(chan []util.Order) //Orderslice of each elevator
	var sendOrder = make(chan util.Order)
	var orders [util.Nslaves][]util.Order
	//local process channels
	slavesChan := make(chan [util.Nslaves]util.Elevator)
	orderChan := make(chan [util.Nslaves][]util.Order)
	orderBackupChan := make(chan [util.Nslaves][]util.Order)
	slavesBackupChan := make(chan [util.Nslaves]util.Elevator)
	slaveAliveChan := make(chan [util.Nslaves]bool)

	slavesChan <- slaves
	orderChan <- orders
	slaveAliveChan <- slaveAlive
	for j := 0; j < util.Nslaves; j++ {
		orders[j] = make([]util.Order, 0)
	}
	firstTry := true

	if isBackup && firstTry {
		firstTry = false
		tmr := time.NewTimer(5 * time.Second)
		//listening to master	
		ordersFromMaster := make(chan [util.Nslaves][]util.Order)
		statusFromMaster := make(chan [util.Nslaves]util.Elevator)
		slaveAliveFromMaster := make(chan [util.Nslaves]bool)
		go bcast.Receiver(20011, ordersFromMaster, statusFromMaster,slaveAliveFromMaster)
		go func() {
			for {
				orders =<-ordersFromMaster
				slaves =<-statusFromMaster
				slaveAlive =<-slaveAliveFromMaster
				tmr.Reset(5 * time.Second)
			}
			}()
			go func () {
				<-tmr.C
				isBackup = false
				firstTry = true
				select{
				case <- slavesChan:
				default:
				}
				select {
				case <- orderChan:
				default:
				}
				select {
				case <- slaveAliveChan:
				default:
				}
				slavesChan <-slaves
				orderChan <- orders
				slaveAliveChan <- slaveAlive
				myIP, _ := localip.LocalIP()
				for i:=0;i<len(slaveIPs);i++{
					if slaveIPs[i] != myIP && slaveAlive[i] {
						spawnMasterBackup := exec.Command("gnome-terminal", "-x", "sh", "-c", "sshpass -p 'Sanntid15' ssh student@", slaveIPs[i], "go run /home/student/Documents/Group55/TTK4145/Project/main.go -startMasterBackup")
						spawnMasterBackup.Start()
						break
					}
				}
				}()
			}

			if !isBackup && firstTry {
				firstTry = false
				myIP, _ := localip.LocalIP()


			if  slaves[0].IP == ""{
				for i := 0; i < util.Nslaves; i++ {
					if myIP != slaveIPs[i] {
						InitSlave(slaveIPs[i])
						slaves = <-slavesChan
						slaves[i] = util.Elevator{slaveIPs[i],0,2}
						slavesChan <- slaves
						slaveAlive =<-slaveAliveChan
						slaveAlive[i] = true
						slaveAliveChan <- slaveAlive
					}else{

						spawnBackup := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run /home/student/Documents/Group55/TTK4145/Project/main.go -startSlave")
						spawnBackup.Start()
						slaves = <-slavesChan
						slaves[i] = util.Elevator{slaveIPs[i],0,2}
						slavesChan <- slaves
						slaveAlive =<-slaveAliveChan
						slaveAlive[i] = true
						slaveAliveChan <- slaveAlive
					}
				}
			}

		//start backup master on remote pc, take first in list that is not itself
		for i:=0;i<len(slaveIPs);i++{
			slaveAlive =<-slaveAliveChan
			slaveAliveChan <- slaveAlive
			if slaveIPs[i] != myIP && slaveAlive[i]{
				spawnMasterBackup := exec.Command("gnome-terminal", "-x", "sh", "-c", "nohup ssh student@", slaveIPs[i], "go run /home/student/Documents/Group55/TTK4145/Project/main.go -startMasterBackup")
				spawnMasterBackup.Start()
			}
		}
		//Set up communication
		go bcast.Transmitter(20009, sendOrder)
		go bcast.Receiver(20009, orderChannel, orderSliceChannel, statusChannel)
			
		go distributeOrder(orderChannel, sendOrder, orderChan)
		go bcast.Transmitter(20011,orderBackupChan, slavesBackupChan)
		go func() {
			for {
				orders =<-orderChan
				orderChan <- orders
				orderBackupChan <- orders
				slaves =<-slavesChan
				slavesChan <- slaves
				slavesBackupChan <- slaves
				time.Sleep(1*time.Second)
			}
		}()

		//updating info from slave
		go func() {
			for {
				orderSlice := <-orderSliceChannel
				if len(orderSlice) != 0{
					for i:=0;i<util.Nslaves;i++{
						if orderSlice[0].ThisElevator.IP == slaveIPs[i] {
							orders =<-orderChan
							orders[i] = orderSlice
							orderChan <- orders

						}
					}
				}
				status :=<-statusChannel

				for i:=0;i<util.Nslaves;i++{
					if status.IP == slaveIPs[i] {
						slaves=<-slavesChan
						slaves[i] = status
						slavesChan <- slaves

					}
				}
			}
			}()

		}

	}
