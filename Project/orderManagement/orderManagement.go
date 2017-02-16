package orderManagement

import(
"time"
"../util"
)


var orderSlice = make([]util.Order,0) //slice of orders

// 1 if success, 0 if duplicate order
func AddOrder(orderChan chan util.Order, floor, button int, elevator util.Elevator, atTime time.Time) int {
	order := util.Order{elevator,util.Button{floor,button}, atTime}
	//TODO: check somehow if success
	orderChan <- order
	return 1
}

func removeOrder(order util.Order, orderSlice []util.Order) []util.Order {
	for i := 0; i < len(orderSlice); i++{
		if(orderSlice[i].FromButton.Floor == order.FromButton.Floor && orderSlice[i].FromButton.TypeOfButton == order.FromButton.TypeOfButton){
			
			orderSlice = append(orderSlice[:i], orderSlice[i+1:]...)
			return orderSlice
		}
	}
	panic("Order not found.. ")
	return orderSlice
}

//returns true if the order already exists in the slice 
func duplicateOrder(order util.Order, orderSlice []util.Order) bool {
	for i := 0; i < len(orderSlice); i++{
		if(orderSlice[i].FromButton.Floor == order.FromButton.Floor && orderSlice[i].FromButton.TypeOfButton == order.FromButton.TypeOfButton){
			return true
		}
	}
	return false
}

func prioritizeOrder() {
	//TODO: walk through order slice and order them according to priority
}

func findSuitableElevator() util.Elevator {
	//TODO: Real functinonality
	return util.Elevator{0,"this"}
}

func calculateCost() int {
	//TODO: Functionality, maybe add parameters.
	return 0
}
