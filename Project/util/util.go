package util 
import(
"time"
)

var Nslaves = 3
type Direction int
const (
	Up Direction = iota
	Down
	Stop
)
type Order struct {
	ThisElevator Elevator
	FromButton Button
	AtTime time.Time
}

type Elevator struct {
	ID int
	IP string
	LastFloor int
	ElevDirection Direction
}
type Button struct {
	Floor, TypeOfButton int
}
