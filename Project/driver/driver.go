package driver

import (
	"io"
	"channels"
	"fmt"
)
const motorSpeed int = 2800
type Direction struct{
	Up int = 1
	Down int = -1
	Stop int = 0
}
func initElevator() {
	init_success := io.ioInit()
	if initSuccess == 0{
		fmt.Println("Could not initialize hardware")
	}
	setFloorIndicator(0)
	setDoorLamp()

}

func ListenForButtons() {}

func getCurrentFloor() int {
	if io.ioReadBit(channel.SENSOR_FLOOR1)==1 {
		return 0
	} else if io.ioReadBit(channel.SENSOR_FLOOR1)==1 {
		return 1
	} else if io.ioReadBit(channel.SENSOR_FLOOR1)==1 {
		return 2
	} else if io.ioReadBit(channel.SENSOR_FLOOR1)==1 {
		return 3
	} else {
		return -1
	}
}

func setButtonLamp() {}

func setDoorLamp() {}

// up, down or stop

func steerElevator(dir Direction) {
	if dir==Up {
		io.ioWriteAnalog(channels.MOTOR, 0)
	}
	else if dir > 0 {
		io.ioClearBit(channels.MOTORDIR)
		ioWriteAnalog(channels.MOTOR,motorSpeed)
	}
	else if dir < 0{
		io.ioSetBit(channels.MOTORDIR)
		ioWriteAnalog(channels.MOTOR,motorSpeed)
	}
}

func setFloorIndicator(floor int) {
}
