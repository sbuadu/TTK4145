package util

import (
	"time"
)

const Nslaves = 3
const Nfloors = 4
const DoorOpenTime = 1000 * time.Millisecond
const Nbuttons = 3

var SlaveIPs = [Nslaves]string{"129.241.187.46", "129.241.187.48","129.241.187.38"}

type Direction int

const (
	Up Direction = iota
	Down
	Stop
)

type Order struct {
	ThisElevator Elevator
	FromButton   Button
	AtTime       time.Time
	Completed    bool
}

type Elevator struct {
	IP            string
	LastFloor     int
	ElevDirection Direction
}
type Button struct {
	Floor, TypeOfButton int
}
