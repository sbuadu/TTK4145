package orderManagement

import (
	"../util"
	"fmt"
	"math"
	"time"
)

/*
Module: Order Management

date: 14.03.17

This module contains functions for handeling orders and orderslices.
It is a module supporting both the master and slave module.
*/

func AddOrder(orderChan, otherOrderChan chan []util.Order, floor, button int, elevator util.Elevator, atTime time.Time) int {

	order := util.Order{elevator, util.Button{floor, button}, atTime, false}
	//fmt.Println("OM Adding order", order.FromButton.Floor)
	orderSlice := <-orderChan
	otherOrders := <-otherOrderChan

	if duplicateOrder(order, orderSlice) || duplicateOrder(order, otherOrders) {

		orderChan <- orderSlice
		otherOrderChan <- otherOrders
		fmt.Println("ignored duplicate order")
		return 0
	} else {
		orderSlice = prioritizeOrder(order, orderSlice, elevator)
		orderChan <- orderSlice
		otherOrderChan <- otherOrders
		fmt.Println("Ready to execute order: ", order.FromButton.Floor)
		return 1
	}
}

func RemoveOrder(order util.Order, orderSlice []util.Order) []util.Order {
	for i := 0; i < len(orderSlice); i++ {
		if orderSlice[i].FromButton.Floor == order.FromButton.Floor && orderSlice[i].FromButton.TypeOfButton == order.FromButton.TypeOfButton {
			orderSlice = append(orderSlice[:i], orderSlice[i+1:]...)
			i -= 1
			if len(orderSlice) == 0 {
				return []util.Order{}
			}

		}
	}
	return orderSlice
}

func duplicateOrder(order util.Order, orderSlice []util.Order) bool {
	for i := 0; i < len(orderSlice); i++ {
		if orderSlice[i].FromButton.Floor == order.FromButton.Floor && orderSlice[i].FromButton.TypeOfButton == order.FromButton.TypeOfButton {
			return true
		}
	}
	return false
}

func prioritizeOrder(order util.Order, orderSlice []util.Order, elevator util.Elevator) []util.Order {

	if len(orderSlice) < 1 { //if no other orders

		return append(orderSlice, order)

	} else if order.FromButton.TypeOfButton == 2 { //internal order
		index := -1
		for i := 0; i < len(orderSlice)-1; i++ {
			if elevator.ElevDirection == 0 && order.FromButton.Floor > orderSlice[i].FromButton.Floor && orderSlice[i].FromButton.Floor > elevator.LastFloor { //checking of the next orders could go first
				index = i
			} else if elevator.ElevDirection == 1 && order.FromButton.Floor < orderSlice[i].FromButton.Floor && orderSlice[i].FromButton.Floor < elevator.LastFloor {
				index = i
			}
		}

		new_item := append(make([]util.Order, 0), order)
		if index == -1 {
			return append(new_item, orderSlice[0:]...)
		} else {
			return append(orderSlice[:index+1], append(new_item, orderSlice[index+1:]...)...)
		}
	} else { //external order
		for i := 0; i < len(orderSlice)-1; i++ {
			// going up
			new_item := append(make([]util.Order, 0), order)
			if elevator.ElevDirection == 0 && order.FromButton.TypeOfButton == 0 && elevator.LastFloor < order.FromButton.Floor && order.FromButton.Floor < orderSlice[i].FromButton.Floor {
				return append(orderSlice[:i], append(new_item, orderSlice[i:]...)...)
				//elevator going down
			} else if elevator.ElevDirection == 1 && order.FromButton.TypeOfButton == 1 && elevator.LastFloor > order.FromButton.Floor && order.FromButton.Floor > orderSlice[i].FromButton.Floor {
				return append(orderSlice[:i], append(new_item, orderSlice[i:]...)...)
			}
		}
		return append(orderSlice, order)
	}
}

func FindSuitableElevator(order util.Order, slaves []util.Elevator) util.Elevator {
	elevIndex := 0
	bestCost := 0
	if order.FromButton.TypeOfButton == 2 {
		return order.ThisElevator
	}
	for i := 0; i < len(slaves); i++ {
		cost := calculateCost(order, slaves[i])
		if cost > bestCost {
			elevIndex = i
			bestCost = cost
		}
	}

	return slaves[elevIndex]
}

//the higher the cost the better the elevator is suited for the order
func calculateCost(order util.Order, elevator util.Elevator) int {

	var distance = int(math.Abs(float64(elevator.LastFloor) - float64(order.FromButton.Floor)))

	if (elevator.ElevDirection == 0 && order.FromButton.TypeOfButton == 0 && (elevator.LastFloor < order.FromButton.Floor)) || (elevator.ElevDirection == 1 && order.FromButton.TypeOfButton == 1 && elevator.LastFloor > order.FromButton.Floor) {
		return 6 + distance*2
	} else if (elevator.ElevDirection == 0 && order.FromButton.TypeOfButton == 1 && elevator.LastFloor < order.FromButton.Floor) || (elevator.ElevDirection == 1 && order.FromButton.TypeOfButton == 0 && elevator.LastFloor > order.FromButton.Floor) {
		return 5 + distance*2
	} else {
		return distance
	}
}
