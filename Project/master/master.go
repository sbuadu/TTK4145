package master

import (
"../network/bcast"
"../network/localip"
"../util"
"os/exec"
"time"
"fmt"
"../orderManagement"
"bytes"
)

var slaveIPs = [util.Nslaves]string{"129.241.187.153"}
// var slaveIPs = [util.Nslaves]string{"129.241.187.153"}

//tested: works
func InitSlave(IP string) {
	spawnSlave := exec.Command("bash","./startSlave.sh",IP,"-startSlave")
	//spawnSlave := exec.Command("sshpass -p Sanntid15 ssh student@", IP, " go run /home/student/Documents/Group55/TTK4145/Project/main.go -startSlave")
	//spawnSlave := exec.Command("sshpass", "-p", "Sanntid15","ssh","student@",IP,"go run /home/student/Documents/Group55/TTK4145/Project/main.go -startSlave")
	 var out bytes.Buffer
    spawnSlave.Stdout = &out
    err := spawnSlave.Start()
    if err != nil {
        fmt.Println(err)
    }
}

//tested: works
func sendOrder(order util.Order, sendOrdersChannel chan util.Order) {
	for i :=0; i<3; i++{
		sendOrdersChannel <- order
		time.Sleep(1*time.Second)
	}
}

//tested: works
func DistributeIncompleteOrder(order util.Order, sendOrdersChannel chan util.Order, orderChan chan [util.Nslaves][]util.Order, slaveAliveChan chan [util.Nslaves]bool, slavesChan chan [util.Nslaves]util.Elevator) {
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

	order.ThisElevator = orderManagement.FindSuitableElevator(order, workingSlaves)
	fmt.Print("This elevator: ",order.ThisElevator.IP)
	fmt.Println("  should go to floor:" ,order.FromButton.Floor)
	go sendOrder(order, sendOrdersChannel)
	
			for i :=0;i<util.Nslaves;i++{ //adding the order to the right backup slice
				if order.ThisElevator.IP == slaveIPs[i]{
					orders := <-orderChan
					tempChan := make(chan []util.Order,1)
					tempChan <- orders[i]
					dummyChan := make(chan []util.Order,1)
					dummySlice := make([]util.Order,0)
					dummyChan <- dummySlice

					orderManagement.AddOrder(tempChan, dummyChan, order.FromButton.Floor, order.FromButton.TypeOfButton, order.ThisElevator, order.AtTime) 
					orders[i] = <- tempChan
					orderChan <- orders

				}

			}
		}
//tested: works
		func distributeOrder(orderChannel, sendOrdersChannel chan util.Order, orderChan chan [util.Nslaves][]util.Order, slaveAliveChan chan [util.Nslaves]bool, slavesChan chan [util.Nslaves]util.Elevator, callbackChannel chan time.Time) {
			for {
				fmt.Println("starting the distribution")
				order :=<- orderChannel
				fmt.Println("Here")
				callbackChannel <- order.AtTime

				if order.Completed { //removing the completed order from the backup slice
					go sendOrder(order, sendOrdersChannel)
					for i :=0;i<util.Nslaves;i++{
						if order.ThisElevator.IP == slaveIPs[i]{
							orders := <-orderChan
							orders[i] = orderManagement.RemoveOrder(order,orders[i])
							orderChan <- orders
						}
					}
				} else {

					DistributeIncompleteOrder(order, sendOrdersChannel, orderChan, slaveAliveChan, slavesChan)


				}
			}
		}

		func MasterLoop(isBackup bool) {
			var slaves [util.Nslaves]util.Elevator
			var slaveAlive [util.Nslaves]bool
	var orders  [util.Nslaves][]util.Order	//A backup of all orders the slaves are to complete


	//channels for communication with slaves
	orderChannel := make(chan util.Order,1)	
	statusChannel := make(chan util.Elevator,1)	
	sendOrdersChannel := make(chan util.Order,1)	
	callbackChannel := make(chan time.Time,1)	


	//local process channels
	slavesChan := make(chan [util.Nslaves]util.Elevator,1)
	orderChan := make(chan [util.Nslaves][]util.Order,1)
	slaveAliveChan := make(chan [util.Nslaves]bool, 1)

	//channels for communication with backup
	orderBackupChan := make(chan [util.Nslaves][]util.Order,1)
	slavesBackupChan := make(chan [util.Nslaves]util.Elevator,1)
	slaveAliveBackupChan := make(chan [util.Nslaves]bool,1)

	slavesChan <- slaves
	orderChan <- orders
	slaveAliveChan <- slaveAlive


	for j := 0; j < util.Nslaves; j++ {
		orders[j] = make([]util.Order, 0)
	}
	firstTry := true


		//tested: 
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
				fmt.Println("receiving update from Master")
				orders =<-ordersFromMaster
				slaves =<-statusFromMaster
				//slaveAlive =<-slaveAliveFromMaster //this is not done as often as the two others, should be moved.. 
				tmr.Reset(5 * time.Second)
			}
			}()

			go func(){
				slaveAlive =<- slaveAliveFromMaster
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

					}()
				}
			//tested:
				if !isBackup && firstTry {
					fmt.Println("I am the master")
					firstTry = false
					myIP, _ := localip.LocalIP()
					var backupIP string

			if  slaves[0].IP == ""{ //if slaves arent initialized, initialize 
				for i := 0; i < util.Nslaves; i++ {
					InitSlave(slaveIPs[i])
						slaves = <-slavesChan
						slaves[i] = util.Elevator{slaveIPs[i],0,2}
						slavesChan <- slaves
						slaveAlive =<-slaveAliveChan
						slaveAlive[i] = true
						slaveAliveChan <- slaveAlive
						fmt.Println("Started slave on ", slaves[i].IP)
					}
					fmt.Println("here")
					slaveAliveBackupChan <- slaveAlive
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
						fmt.Println("Spawning a backup on IP", backupIP)
						spawnMasterBackup := exec.Command("gnome-terminal", "-x", "sh", "-c", "sshpass -p Sanntid15 ssh student@", slaveIPs[i], " go run /home/student/Documents/Group55/TTK4145/Project/main.go -startMasterBackup")
						spawnMasterBackup.Start()
						break
					}
				}

		//Set up communication
				go bcast.Transmitter("255.255.255.255",20009, sendOrdersChannel, callbackChannel)
				go bcast.Receiver(20008, orderChannel, statusChannel)
				go bcast.Transmitter(backupIP,20011,orderBackupChan, slavesBackupChan, slaveAliveBackupChan)


				go distributeOrder(orderChannel, sendOrdersChannel, orderChan,slaveAliveChan,slavesChan, callbackChannel)


		//tested: 
		//sending updates to backup
				go func() {
					for {
						orders =<-orderChan
						orderChan <- orders
						slaves =<-slavesChan
						slavesChan <- slaves

						orderBackupChan <- orders	
						slavesBackupChan <- slaves

						time.Sleep(1*time.Second)
					}
					}()

		//tested: works
		//updating info from slave
					go func() {
						for {
							status :=<-statusChannel
							for i:=0;i<util.Nslaves;i++{
								if status.IP == slaveIPs[i] {
									slaves=<-slavesChan
									slaves[i] = status
									slavesChan <- slaves
									timers[i].Reset(5*time.Second)
								}
							}
						}
						}()

		//must test if works nomatter which order the slave IPs are listed
		// checking for non-responsive slaves and working accordingly
						go func(){
							for {
								for j:= 0;j<util.Nslaves;j++{
									select{
									case <-timers[j].C:
										fmt.Println("Slave is dead. IP: ", slaveIPs[j])
										slaveAlive =<-slaveAliveChan
										slaveAlive[j] = false
										slaveAliveChan <- slaveAlive
										slaveAliveBackupChan <- slaveAlive
										orders =<-orderChan
										orderChan<- orders
										fmt.Println("redistibuting dead slaves orders")
										for i:=0;i<len(orders[j]);i++{

											if !(orders[j][i].FromButton.TypeOfButton == 2) {
												orders[j][i].Completed = true
												go sendOrder(orders[j][i],sendOrdersChannel)
												orders[j][i].Completed = false

												DistributeIncompleteOrder(orders[j][i], sendOrdersChannel , orderChan,  slaveAliveChan, slavesChan)


											}
										}

									}
								}
							}
							}()

						}

					}

