package util

import (
	"time"
)

const Nslaves = 3
const Nfloors = 4

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
	Completed 	 bool

}

type Elevator struct {
	ID            int
	IP            string
	LastFloor     int
	ElevDirection Direction
}
type Button struct {
	Floor, TypeOfButton int
}

const DoorOpenTime = 1000 * time.Millisecond
