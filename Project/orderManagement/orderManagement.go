package orderManagement

import(
"time"
"../util"
//"math"
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

func findSuitableElevator(slaves [3]util.Elevator, order util.Order) util.Elevator {


	//TODO: Real functinonality


	return util.Elevator{0,"this"}
}

/*not ready work in progress
func calculateCost(elevator util.Elevator, order util.Order) int {
var suitableDir = 0;
//var distance = Abs(elevator.LastFloor - order.Button.Floor)  
if (elevator.Direction == 0 && order.Button.TypeOfButton == 0 && elevator.LastFloor < order.Button.Floor) || (elevator.Direction == 1 && order.Button.TypeOfButton == 1 && elevator.LastFloor > order.Button.Floor) {
suitableDir = 1; 
}


	return 0
}
*/