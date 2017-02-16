package slave

import(
	"time"
	"../orderManagement"
	"../driver"
	"../network/localip"
	"../util"
)

//TODO: make process pair functionality


func SendMessage() {
	//TODO:
}

func ListenRemoteOrders() {
	//TODO: Listen for new orders and add them to the orders
}

func ListenLocalOrders(orderChan chan util.Order) {

	//TODO: check if button is already on
	var buttons [4][3]int
	for {
		recent := buttons
		buttons = driver.ListenForButtons()
		changed,floor,button := CompareMatrix(buttons,recent)

		if changed {
			IP,_ := localip.LocalIP()

			success := orderManagement.AddOrder(orderChan, floor, button,util.Elevator{1,IP},time.Now())
			if success == 1 {
				//TODO: Move this to after order is appended to orders
				driver.SetButtonLamp(floor, button, 1)
			}
		}
		time.Sleep(10*time.Millisecond)
	}
}

func ExecuteOrder() {
	for {
		if len(orders) > 0{
			currentOrder := orders[0]
			currentFloor := driver.GetCurrentFloor()
			if currentFloor == currentOrder.FromButton.Floor {
				driver.SetDoorLamp(1)
			} else if currentFloor < currentOrder.FromButton.Floor {
				driver.SteerElevator(0)
				for ;currentFloor < currentOrder.FromButton.Floor; {
					if driver.GetCurrentFloor() != -1 {
						currentFloor = driver.GetCurrentFloor()
					}
				}
			} else if currentFloor > currentOrder.FromButton.Floor {
				driver.SteerElevator(1)
				for ;currentFloor > currentOrder.FromButton.Floor; {
					if driver.GetCurrentFloor() != -1 {
						currentFloor = driver.GetCurrentFloor()
					}
				}
			}
		}
	}
}


func CompareMatrix(newMatrix, oldMatrix [4][3]int) (changed bool, row, column int) {
	for i:=0;i<4;i++{
		for j:=0;j<3;j++{
			if newMatrix[i][j] != oldMatrix[i][j] {
				changed = true
				row = i
				column = j
				return changed, row, column
			}
		}
	}
	return false,0,0
}
var orders = make([]util.Order,0)

func Slave() {
	//var isBackup bool
	driver.InitElevator()
	orderChan := make(chan util.Order)
	go ListenLocalOrders(orderChan)
	go ExecuteOrder()
	for {
		orders = append(orders,<-orderChan)
	}
}
