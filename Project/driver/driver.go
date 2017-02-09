package driver

import (

	//"io"
	//"channels"
	//"math"
"fmt"
	//"iota"
)


const(
	N_FLOORS = 4
	N_BUTTONS = 3

	BUTTON_CALL_UP = 0
	BUTTON_CALL_DOWN = 1
	BUTTON_COMMAND = 2

	MOTOR_SPEED = 2800
)


type Direction int
const (
	Up Direction = iota
	Down
	Stop
)


var	lamp_channel_matrix = [4][3]int{
	{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
	{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
	{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
	{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
}

var	button_channel_matrix = [4][3]int{
	{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
	{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
	{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
	{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},

}



func InitElevator() {
	initSuccess := ioInit()
	if initSuccess == 0{
		fmt.Println("Could not initialize hardware")
	}
	SetFloorIndicator(0)
	SetDoorLamp(0)

}

func ListenForButtons() [4][3]int {
	var pushedBtnMatrix [4][3]int
	for i := 0; i < N_FLOORS; i++ {
		for j := 0; j < N_BUTTONS; j++ {
			pushedBtnMatrix[i][j] = ioReadBit(button_channel_matrix[i][j])
			
		}
	}
	return pushedBtnMatrix

}

func GetCurrentFloor() int {
	if ioReadBit(SENSOR_FLOOR1)==1 {
		return 0
	} else if ioReadBit(SENSOR_FLOOR2)==1 {
		return 1
	} else if ioReadBit(SENSOR_FLOOR3)==1 {
		return 2
	} else if ioReadBit(SENSOR_FLOOR4)==1 {
		return 3
	} else {
		return -1
	}
}

func SetButtonLamp(floor int, button int , value int) {

	if floor < 0 || floor >= N_FLOORS || button < 0 || button >= N_BUTTONS{
		panic("Floor or button command out of range")
	} else if value == 1 {
		ioSetBit(lamp_channel_matrix[floor][button])
	} else {
		ioClearBit(lamp_channel_matrix[floor][button])
	}
}


func SetDoorLamp(value int) {
	if value == 1 {
		ioSetBit(LIGHT_DOOR_OPEN)
	}else{
		ioClearBit(LIGHT_DOOR_OPEN)
	}

}


func SetFloorIndicator(floor int) {
	if floor < 0 || floor >= N_FLOORS {
		panic("Floor not existing")
	} else if floor == 0 {
		ioClearBit(LIGHT_FLOOR_IND1)
		ioClearBit(LIGHT_FLOOR_IND2)
	} else if floor == 1 {
		ioClearBit(LIGHT_FLOOR_IND1)
		ioSetBit(LIGHT_FLOOR_IND2)
	} else if floor == 2 {
		ioSetBit(LIGHT_FLOOR_IND1)
		ioClearBit(LIGHT_FLOOR_IND2)
	} else if floor == 3 {
		ioSetBit(LIGHT_FLOOR_IND1)
		ioSetBit(LIGHT_FLOOR_IND2)
	}
}

func SteerElevator(dir Direction) {
	switch dir {
	case Stop:
		ioWriteAnalog(MOTOR, 0)
		fmt.Println("Stopping elevator")
	case Up:
		ioClearBit(MOTORDIR)
		ioWriteAnalog(MOTOR,MOTOR_SPEED)
		fmt.Println("Going up")
	case Down:
		ioSetBit(MOTORDIR)
		ioWriteAnalog(MOTOR,MOTOR_SPEED)
		fmt.Println("Going down")
	}
}
