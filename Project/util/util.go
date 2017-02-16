package util 
import(
"time"

"../driver"
)

type Order struct {
	ThisElevator Elevator
	FromButton Button
	AtTime time.Time
}

type Elevator struct {
	ID int
	IP string
	lastFloor int
	direction driver.Direction
}
type Button struct {
	Floor, TypeOfButton int
}