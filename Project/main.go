package main

import (
	"./driver"
	"fmt"
	"time"
)

type Direction int
const (
        Up = Direction(iota)
        Down
        Stop
)

func main() {
	fmt.Println("Initializing elevator")
	driver.InitElevator()
	driver.SteerElevator(driver.Direction(0))
	time.Sleep(1500*time.Millisecond)
	driver.SteerElevator(1)
	time.Sleep(1*time.Second)
	driver.SteerElevator(2)
	time.Sleep(1*time.Second)
	fmt.Println(driver.GetCurrentFloor())
	for i:=0; i < 4; i++{
		for j := 0; j < 3; j++ {
			driver.SetButtonLamp(i,j,1)
			time.Sleep(500*time.Millisecond)
			driver.SetButtonLamp(i,j,0)
		}
	}
	driver.SetDoorLamp(1)
	time.Sleep(1*time.Second)
	driver.SetDoorLamp(0)
	for k := 0; k < 4; k++ {
		driver.SetFloorIndicator(k)
	}
	
	defer fmt.Println("All testing done")
}
