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

//var slaveIPs = [util.Nslaves]string{"129.241.187.153","129.241.187.157"}
var slaveIPs = [util.Nslaves]string{"129.241.187.153"}
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

func distributeOrder(orderChannel, sendOrders chan util.Order, orderChan chan [util.Nslaves][]util.Order, slaveAliveChan chan [util.Nslaves]bool, slavesChan chan [util.Nslaves]util.Elevator) {
	for {
	//fmt.Println("starting the distribution")
		order <- orderChannel
		if order.Completed {
			go sendOrder(order, sendOrders)
			for i :=0;i<util.Nslaves;i++{
				if order.ThisElevator.IP == slaveIPs[i]{
					orders := <-orderChan
					orders[i] = orderManagement.RemoveOrder(order,orders[i])
					orderChan <- orders

					fmt.Println("These are the orders after one is completed: ", orders[i])
				}
			}
		} else {
			fmt.Println("Order not complete")
			var workingSlaves = make([]util.Elevator,0)
			for i := 0; i<util.Nslaves; i++{
				slaveAlive :=<-slaveAliveChan
				slaves :=<-slavesChan
				slaveAliveChan<-slaveAlive
				slavesChan<-slaves
				if slaveAlive[i]{
					workingSlaves = append(workingSlaves, slaves[i])
				}
			}
			order.ThisElevator = orderManagement.FindSuitableElevator(workingSlaves, order)
			go sendOrder(order, sendOrders)
			for i :=0;i<util.Nslaves;i++{
				if order.ThisElevator.IP == slaveIPs[i]{
					orders := <-orderChan
					orderChan <- orders
					tempChan := make(chan []util.Order)
					tempChan <- orders[i]
					orderManagement.AddOrder(tempChan, make(chan []util.Order), order.FromButton.Floor, order.FromButton.TypeOfButton, order.ThisElevator, order.AtTime) 
					orders[i] = <- tempChan
					fmt.Println("These are the orders after adding an order: ", orders[i])
				}

			}

		}
	}
}

func Master(isBackup bool) {
	fmt.Println("Setting up master")
	var slaves [util.Nslaves]util.Elevator
	var slaveAlive [util.Nslaves]bool
	var orderChannel = make(chan util.Order)	//orders coming from elevators
	var statusChannel = make(chan util.Elevator,1)	//Status from elevators
	var orderSliceChannel = make(chan []util.Order,1) //Orderslice of each elevator
	var sendOrders = make(chan util.Order,1)		//Broadcast orders to slaves
	var orders  [util.Nslaves][]util.Order	//All orders at slaves

	//local process channels
	slavesChan := make(chan [util.Nslaves]util.Elevator,1)
	orderChan := make(chan [util.Nslaves][]util.Order,1)
	orderBackupChan := make(chan [util.Nslaves][]util.Order)
	slavesBackupChan := make(chan [util.Nslaves]util.Elevator)
	slaveAliveChan := make(chan [util.Nslaves]bool,1)

	slavesChan <- slaves
	orderChan <- orders
	slaveAliveChan <- slaveAlive

	for j := 0; j < util.Nslaves; j++ {
		orders[j] = make([]util.Order, 0)
	}
	firstTry := true
	
	fmt.Println("Master initial setup")

	if isBackup && firstTry {
		fmt.Println("I am a master backup")
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
				fmt.Println("I am the master")
				firstTry = false
				myIP, _ := localip.LocalIP()
				var backupIP string

			if  slaves[0].IP == ""{
				for i := 0; i < util.Nslaves; i++ {
					if myIP != slaveIPs[i] {
						InitSlave(slaveIPs[i])
						
					}else{

						spawnBackup := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run /home/student/Documents/Group55/TTK4145/Project/main.go -startSlave")
						spawnBackup.Start()
						
					}
					slaves = <-slavesChan
					slaves[i] = util.Elevator{slaveIPs[i],0,2}
					slavesChan <- slaves
					slaveAlive =<-slaveAliveChan
					slaveAlive[i] = true
					slaveAliveChan <- slaveAlive
					fmt.Println("Started slave on ", slaves[i].IP)
				}
			}
		//Settig up timers for slaves
		var timers [util.Nslaves] *time.Timer
		for i := 0; i <util.Nslaves; i++ {
			timers[i] = time.NewTimer(20*time.Second)
		}
		//start backup master on remote pc, take first in list that is not itself
		for i:=0;i<len(slaveIPs);i++{
			slaveAlive =<-slaveAliveChan
			slaveAliveChan <- slaveAlive
			if slaveIPs[i] != myIP && slaveAlive[i]{
				backupIP = slaveIPs[i]
				spawnMasterBackup := exec.Command("gnome-terminal", "-x", "sh", "-c", "sshpass -p 'Sanntid15' ssh student@", slaveIPs[i], "go run /home/student/Documents/Group55/TTK4145/Project/main.go -startMasterBackup")
				spawnMasterBackup.Start()
				break
			}
		}
		//Set up communication
		go bcast.Transmitter("255.255.255.255",20009, sendOrders)
		go bcast.Receiver(20009, orderChannel, orderSliceChannel, statusChannel)

		go bcast.Transmitter(backupIP,20011,orderBackupChan, slavesBackupChan)

		go distributeOrder(orderChannel, sendOrders, orderChan,slaveAliveChan,slavesChan)


			//sending updates to backup
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
				fmt.Println("Listening")

				status :=<-statusChannel
				fmt.Println("Received a status update")
				for i:=0;i<util.Nslaves;i++{
					if status.IP == slaveIPs[i] {
						slaves=<-slavesChan
						slaves[i] = status
						fmt.Println("Elevator at floor: ",status.LastFloor)
						slavesChan <- slaves
						timers[i].Reset(10*time.Second)
					}
				}
			}
		}()
		// checking for non-responsive slaves and working accordingly
		go func(){
			for {
			for j:= 0;j<util.Nslaves;j++{
				select{
				case <-timers[j].C:
					fmt.Println("Slave is dead")
					slaveAlive =<-slaveAliveChan
					slaveAlive[j] = false
					slaveAliveChan <- slaveAlive
					orders =<-orderChan
					for i:=0;i<len(orders[j]);i++{
						if !(orders[j][i].FromButton.TypeOfButton == 2) {
							orders[j][i].Completed = true
							go sendOrder(orders[j][i],sendOrders)
							orders[j][i].Completed = false
						}
					}
				}
			}
		}
		}()
		
	}

}

