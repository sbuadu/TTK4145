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

func main() {
	startMaster := flag.Bool("startMaster", false, "a bool")
	startMasterBackup := flag.Bool("startMasterBackup", false, "a bool")
	startSlave := flag.Bool("startSlave", false, "a bool")
	startSlaveBackup := flag.Bool("startSlaveBackup", false, "a bool")
	flag.Parse()
	//ifs and shit
	driver.SteerElevator(2)

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
}
