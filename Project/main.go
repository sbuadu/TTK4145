package main

import (
	"./driver"
	"fmt"
	//"os/exec"
	"./master"
	//"./network/bcast"
	"./slave"
	//"./util"
	//"time"
	//"sync"
	"flag"
)

type Direction int

const (
	Up = Direction(iota)
	Down
	Stop
)

func main() {
	/*
		fmt.Println("Testing the order prioritizing")
		go slave.Test()
		for{

		}


	*/

	startMaster := flag.Bool("startMaster", false, "a bool")
	startMasterBackup := flag.Bool("startMasterBackup", false, "a bool")
	startSlave := flag.Bool("startSlave", false, "a bool")
	startSlaveBackup := flag.Bool("startSlaveBackup", false, "a bool")
	flag.Parse()
	//ifs and shit
	driver.SteerElevator(2)

	//startSlave := exec.Command("gnome-terminal", "-x", "-c","go run /home/student/TTK4145/Project/slave/slave.go -isBackup=false")

	//exec.Command("gnome-terminal", "-x", "go run ~/Documents/TTK4145/Exercise6/backup.go")
	//startSlave.Start()
	if *startMaster {
		fmt.Println("Starting Master")
		go master.MasterLoop(false)
	}
	if *startMasterBackup {
		go master.MasterLoop(true)
	}
	if *startSlave {
		fmt.Println("Starting slave")
		go slave.SlaveLoop(false)
	}
	if *startSlaveBackup {
		fmt.Println("Starting slavebackup")
		go slave.SlaveLoop(true)
	}
	select {}
	//testing that we are able to steer the elevator and return current floor
	/*
	   	driver.SteerElevator(driver.Direction(0))
	   	time.Sleep(4*time.Second)
	   	driver.SteerElevator(driver.Direction(1))
	   	time.Sleep(4*time.Second)
	   =======
	   	time.Sleep(60 * time.Second)
	   >>>>>>> c6b4b8ae90f8a8f8dce133b8b4134b0f18408640
	   	driver.SteerElevator(2)

	   	//testing that we are able to steer the elevator and return current floor
	   	/*
	   	   	driver.SteerElevator(driver.Direction(0))
	   	   	time.Sleep(4*time.Second)
	   	   	driver.SteerElevator(driver.Direction(1))
	   	   	time.Sleep(4*time.Second)
	   	   	driver.SteerElevator(2)
	   	   	time.Sleep(1*time.Second)
	   	   	fmt.Println(driver.GetCurrentFloor())



	   	   /*
	   	   //testing the btn signals
	   	   	for{
	   	   		fmt.Println(driver.ListenForButtons())
	   	   		time.Sleep(2*time.Millisecond)
	   	   	}
	*/

	/*
	   //Testing that the lights work as they should
	   for i:=0; i < 4; i++{
	   	for j := 0; j < 3; j++ {
	   		fmt.Printf("Button lamp floor %d lamp %d\n",i,j)
	   			//fmt.Println(driver.lamp_channel_matrix[i][j])
	   		driver.SetButtonLamp(i,j,1)
	   		time.Sleep(1000*time.Millisecond)
	   		driver.SetButtonLamp(i,j,0)
	   	}
	   }
	   driver.SetDoorLamp(1)
	   time.Sleep(1*time.Second)
	   driver.SetDoorLamp(0)
	   for k := 0; k < 4; k++ {
	   	driver.SetFloorIndicator(k)
	   	fmt.Printf("Light on floor %d\n", k+1)
	   	time.Sleep(1*time.Second)
	   }
	*/

	/*
	   driver.SteerElevator(driver.Direction(0))

	   for {


	   if( driver.GetCurrentFloor() == 1){
	   driver.SteerElevator(driver.Direction(2))
	   return
	   }

	   }*/
}
