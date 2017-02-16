opackage orderManagement

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

var orderSlice = make([]Order,0) //slice of orders

// 1 if success, 0 if duplicate order
func AddOrder(orders chan Order, floor, button int, elevator Elevator, atTime time.Time) int {
	order := Order{elevator,Button{floor,button}, atTime}
	{
		orders <- order
		return 1
	} else {
		return 0
	}
}

func removeOrder(order Order, orderSlice []Order) []Order{
	for i = 0; i < len(orderSlice); ++i{
		if(orderSlice[i].FromButton.Floor == order.FromButton.Floor && orderSlice[i].FromButton.TypeOfButton == order.FromButton.TypeOfButton){
			
			orderSlice = append(orderSlice[:i], orderSlice[i+1:]...)
			return orderSlice
		}
	}
	panic("Order not found.. ")
	return orderSlice
}

//returns true if the order already exists in the slice 
func duplicateOrder(order Order, orderSlice []Order) bool {
	for i = 0; i < len(orderSlice); ++i{
		if(orderSlice[i].FromButton.Floor == order.FromButton.Floor && orderSlice[i].FromButton.TypeOfButton == order.FromButton.TypeOfButton){
			return true
		}
	}
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
