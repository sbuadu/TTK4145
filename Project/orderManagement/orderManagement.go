package orderManagement

import(
	"time"
)
type Order struct {
	elevator Elevator
	fromButton Button
	atTime time.Time
}

type Elevator struct {
	ID int
	IP string
}
type Button struct {
	floor int
	button int
}

func AddOrder(orders chan Order, button Button, elevator Elevator, atTime time.Time) {
	order := Order{elevator,button, atTime}
	orders <- order
}

func removeOrder() {

}

func duplicateOrder() bool {
	return false
}

func prioritizeOrder() {

}

func findSuitableElevator() Elevator {
	return Elevator{0,"this"}
}

func calculateCost() int {
	return 0
}
