package orderManagement

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

// 1 if success, 0 if duplicate order
func AddOrder(orders chan Order, floor, button int, elevator Elevator, atTime time.Time) int {
	order := Order{elevator,Button{floor,button}, atTime}
	if !duplicateOrder(order) {
		orders <- order
		return 1
	} else {
		return 0
	}
}

func removeOrder() {
	//TODO: Remove order from slice
}

func duplicateOrder(order Order) bool {
	//TODO: 
	return false
}

func prioritizeOrder() {
	//TODO: walk through order slice and order them according to priority
}

func findSuitableElevator() Elevator {
	//TODO: Real functinonality
	return Elevator{0,"this"}
}

func calculateCost() int {
	//TODO: Functionality, maybe add parameters.
	return 0
}
