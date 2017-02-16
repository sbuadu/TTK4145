package util 
import(
"time"
)

type Order struct {
	ThisElevator Elevator
	FromButton Button
	AtTime time.Time
}

type Elevator struct {
	ID int
	IP string
}
type Button struct {
	Floor, TypeOfButton int
}