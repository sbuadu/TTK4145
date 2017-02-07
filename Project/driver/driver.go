package driver

import (
	"io"
	"channels"
)

struct 
func initElevator() {
	init_success := io.ioInit()
	if initSuccess == 0{
		fmt.Println("Could not initialize hardware")
	}

}

func ListenForButtons() {}

func getCurrentFloor() {}

func setButtonLamp() {}

func setDoorLamp() {}

// up, down or stop
func steerElevator(dir Direction) {}

func setFloorIndicator(floor int) {}
