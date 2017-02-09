package slave

import(
	"time"
	"../orderManagement"
	"../driver"
	"../network/network/localip"
)
type Elevator struct {
	ID int
	IP string
}

type Order struct {
	elevator Elevator
	fromButton Button
	atTime time.Time
}

type Button struct{
	floor int
	button int
}
func SendMessage() {

}

func ListenRemoteOrders() {
	
}

func ListenLocalOrders(orders chan Order) {
	var buttons [4][3]int
	for {
		recent := buttons
		buttons = driver.ListenForButtons()
		changed,floor,button := compareMatrix(buttons,recent)

		if changed {
			IP,_ := localip.Localip()
			orderManagement.AddOrder(orders,Button{floor,button},Elevator{1,IP},time.Now())
		}
	}	
}

func CompleteOrder() {
}

func compareMatrix(new, old [4][3]int) (changed bool, floor, button int) {
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

func main() {
	driver.InitElevator()
	orderChan := make(chan orderManagement.Order)
	orders := make([]Order,0)
	go ListenLocalOrders(orderChan)
	for {
		orders = append(orders,<-orderChan)	
	}
} 