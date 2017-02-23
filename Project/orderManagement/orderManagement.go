package orderManagement

import (

	"../util"
	"fmt"
	"math"
	"time"
)

var orderSlice = make([]util.Order, 0) //slice of orders

// 1 if success, 0 if duplicate order
func AddOrder(orderChan chan []util.Order, floor, button int, elevator util.Elevator, atTime time.Time) int {
	order := util.Order{elevator, util.Button{floor, button}, atTime}
	//TODO: check somehow if success

	orderSlice := <-orderChan

	if duplicateOrder(order, orderSlice) {
		orderChan <- orderSlice
		return 0
	} else {

		orderSlice = PrioritizeOrder(order, orderSlice)
		orderChan <- orderSlice
		return 1
	}
}

func RemoveOrder(order util.Order, orderSlice []util.Order) []util.Order {
	for i := 0; i < len(orderSlice); i++ {
		if orderSlice[i].FromButton.Floor == order.FromButton.Floor && orderSlice[i].FromButton.TypeOfButton == order.FromButton.TypeOfButton {

			orderSlice = append(orderSlice[:i], orderSlice[i+1:]...)
			fmt.Println("Removing order", order.FromButton.TypeOfButton, order.FromButton.Floor)
			if len(orderSlice) == 0 {
				return []util.Order{}
			}
			return orderSlice
		}
	}
	//panic("Order not found.. ")
	return orderSlice
}

//returns true if the order already exists in the slice

func duplicateOrder(order util.Order, orderSlice []util.Order) bool {
	for i := 0; i < len(orderSlice); i++ {
		if orderSlice[i].FromButton.Floor == order.FromButton.Floor && orderSlice[i].FromButton.TypeOfButton == order.FromButton.TypeOfButton {
			return true
		}
	}
	return false
}

func PrioritizeOrder(order util.Order, orderSlice []util.Order, elevator util.Elevator) []util.Order {
	
	if len(orderSlice) < 1 { //if no orders
		return append(orderSlice, []util.Order{order}...)

	}else if order.FromButton.TypeOfButton == 3{ //internal order 
		 index := -1
		for i := 0; i < len(orderSlice)-1; i++ {
			if  elevator.ElevDirection == 0 && order.FromButton.Floor > orderSlice[i].FromButton.Floor && orderSlice[i].FromButton.Floor > elevator.LastFloor{ //checking of the next orders could go first 
				index = i 
			}	else if elevator.ElevDirection == 1 && order.FromButton.Floor < orderSlice[i].FromButton.Floor && orderSlice[i].FromButton.Floor < elevator.LastFloor {
					index = i 
				}else 
				break; 
			}
		}

		if index == -1 {
			return append([]util.Order{order}, orderSlice, )
		}else{
			return append(orderSlice[:i+1], append([]util.Order{order}, orderSlice[i+1:]...)...)
		}
	}else{ //external order
		for i := 0; i < len(orderSlice)-1; i++ {
			// going up
			if elevator.ElevDirection == 0 && order.FromButton.TypeOfButton == 0 && elevator.LastFloor < order.FromButton.Floor && order.FromButton.Floor < orderSlice[i].FromButton.Floor {
				return append(orderSlice[:i], append([]util.Order{order}, orderSlice[i:]...)...)
				//elevator going down
			} else if  elevator.ElevDirection == 1 && order.FromButton.TypeOfButton == 1 && elevator.LastFloor > order.FromButton.Floor && order.FromButton.Floor > orderSlice[i].FromButton.Floor {
				return append(orderSlice[:i], append([]util.Order{order}, orderSlice[i:]...)...)
			}
		}
		return append(orderSlice, []util.Order{order}...)
	}
}

func findSuitableElevator(slaves [3]util.Elevator, order util.Order) util.Elevator {
	elevIndex := 0
	bestCost := 0
	for i := 0; i < len(slaves); i++ {
		cost := calculateCost(slaves[i], order)
		if cost > bestCost {
			elevIndex = i
			bestCost = cost
		}
	}

	return slaves[elevIndex]
}

/*
When calculating the cost there are three cases to be considered
1. The call is in the direction of travel and the elevator travelling in direction of the call
2. The call is in the oposite direction of travel, but the elevator is travelling in the direction of the call
3. The call is in the oposite direction of travel, and the elevator is travelling in the oposite direction of the call

the higher the cost the better
*/

func calculateCost(elevator util.Elevator, order util.Order) int {

	var distance = int(math.Abs(float64(elevator.LastFloor) - float64(order.FromButton.Floor)))

	if (elevator.ElevDirection == 3) || (elevator.ElevDirection == 0 && order.FromButton.TypeOfButton == 0 && (elevator.LastFloor < order.FromButton.Floor)) || (elevator.ElevDirection == 1 && order.FromButton.TypeOfButton == 1 && elevator.LastFloor > order.FromButton.Floor) {
		return 6 + distance*2
	} else if (elevator.ElevDirection == 0 && order.FromButton.TypeOfButton == 1 && elevator.LastFloor < order.FromButton.Floor) || (elevator.ElevDirection == 1 && order.FromButton.TypeOfButton == 0 && elevator.LastFloor > order.FromButton.Floor) {
		return 5 + distance*2
	} else {
		return 1
	}
}

func sendOrder(){
	//TO DO
}