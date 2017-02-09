package main

import (
	"./driver"
	"fmt"
	"time"
	"./slave"
)

type Direction int
const (
        Up = Direction(iota)
        Down
        Stop
)

func main() {
	fmt.Println("Initializing Slave")

//testing that we are able to steer the elevator and return current floor
/*
	driver.SteerElevator(driver.Direction(0))
	time.Sleep(4*time.Second)
	driver.SteerElevator(driver.Direction(1))
	time.Sleep(4*time.Second)
	driver.SteerElevator(2)
	time.Sleep(1*time.Second)
	fmt.Println(driver.GetCurrentFloor())

*/
/*testing the btn signals
for{
	fmt.Println(driver.ListenForButtons())
time.Sleep(2*time.Millisecond)
}
*/


/* Testing that the lights work as they should
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
for {
	go Slave()
}
defer fmt.Println("All testing done")
}
}
