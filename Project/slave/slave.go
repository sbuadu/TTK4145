package slave

import(
	"time"
	"../orderManagement"
	"../driver"
	"../network/localip"
)

//TODO: make process pair functionality


func SendMessage() {
	//TODO:
}

func ListenRemoteOrders() {
	//TODO: Listen for new orders and add them to the orders
}

func ListenLocalOrders(orders chan orderManagement.Order) {
	//TODO: check if button is already on
	var buttons [4][3]int
	for {
		recent := buttons
		buttons = driver.ListenForButtons()
		changed,floor,button := compareMatrix(buttons,recent)

		if changed {
			IP,_ := localip.LocalIP()
			success := orderManagement.AddOrder(orders, floor, button, orderManagement.Elevator{1,IP},time.Now())
			if success == 1 {
				driver.SetButtonLamp(floor, button, 1)
			}
		}
		time.Sleep(10*time.Millisecond)
	}
}

func CompleteOrder() {
	//TODO: pick first order in list and perform
}

func CompareMatrix(new, old [4][3]int) (changed bool, floor, button int) {
	for i:=0;i<4;i++{
		for j:=0;j<3;j++{
			if new[i][j] != old[i][j] {
				changed = true
				floor = i
				button = j
			}
		}
	}
	return changed, floor, button
}
var orders = make([]orderManagement.Order,0)

func Slave() {
	driver.InitElevator()
	orderChan := make(chan orderManagement.Order)
	go ListenLocalOrders(orderChan)
	go CompleteOrder()
	for {
		orders = append(orders,<-orderChan)
	}
}
