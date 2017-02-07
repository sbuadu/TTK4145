package driver

import (

"io"
"channels"
"math"
"fmt"
)


const(
N_FLOORS = 4
N_BUTTONS = 3

BUTTON_CALL_UP = 0 
BUTTON CALL_DOWN = 1
BUTTON_COMMAND = 2

MOTOR_SPEED = 2800
)


var	lamp_channel_matrix = [4][3]int{
	{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
	{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
	{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
	{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
}

var	button_channel_matrix = [2][3]int{
	{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
	{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},

}




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

func setButtonLamp(button int, floor int , value int) {
//must check that the button number and floor is valid..
//must find a way to handle this type of error 
//is this written to a log? 

if value {
	io.ioSetBit(lamp_channel_matrix[floor][button])
}else{
	io.ioClearBit(lamp_channel_matrix[floor][button])
}


}





func setDoorLamp() {}

// up, down or stop


func setFloorIndicator(floor int) {}
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
