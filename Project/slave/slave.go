package slave

import(
	"time"
	"../orderManagement"
	"../driver"
	"../network/localip"
	"../network/bcast"
	"../util"
	"fmt"
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
		var recent [4][3]int
		buttons = driver.ListenForButtons()
		changed,floor,button := CompareMatrix(buttons,recent)

		if changed {
			IP,_ := localip.LocalIP()
			//TODO: check if order is duplicate
			success := orderManagement.AddOrder(orderChan, floor, button,util.Elevator{1,IP,0,1},time.Now())
			if success == 1 {
				//TODO: Move this to after order is appended to orders
				driver.SetButtonLamp(floor, button, 1)
			}
			changed = false
		}
		time.Sleep(10*time.Millisecond)
	}
}
func goToFloor(order util.Order, currentFloor int) {
	orderFloor := order.FromButton.Floor
	higher := currentFloor < orderFloor
	var elevDir util.Direction
	if higher {
		elevDir = 0
	} else if !higher {
		elevDir = 1
	}
	driver.SetDoorLamp(0)
	driver.SteerElevator(elevDir)
	for ;currentFloor != orderFloor; {
		floor := driver.GetCurrentFloor()
		if floor != -1 {
			currentFloor = floor
			driver.SetFloorIndicator(currentFloor)
		}
	}
	driver.SteerElevator(2)
	driver.SetButtonLamp(orderFloor, order.FromButton.TypeOfButton,0)
}

func ExecuteOrder(orderChan chan util.Order) {
	currentFloor := driver.GetCurrentFloor()
	if currentFloor == -1 {
		currentFloor = 0
	}
	for {
		currentOrder := <-orderChan
		floor := currentOrder.FromButton.Floor
		currentFloor = driver.GetCurrentFloor()
		if currentFloor == floor {
			driver.SetDoorLamp(1)
		} else {
			goToFloor(currentOrder,currentFloor)
		}
		time.Sleep(10*time.Millisecond)
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

var orders = make([]util.Order,10)

func main() {
	var isBackup bool
	driver.InitElevator()
	orderChan := make(chan util.Order,100)
	backup := make(chan util.Order)
	if !isBackup {
		go ListenLocalOrders(orderChan)
		go ExecuteOrder(orderChan)
		go bcast.Transmitter(20009,backup)
	}  else if isBackup {
		go ListenLocalOrders(orderChan)
		go bcast.Receiver(20009,backup)
	}
}


func test(){

	driver.InitElevator()
	orderChan := make(chan util.Order,100)

	go ListenLocalOrders(orderChan)
orderSlice := []util.Order{}
	for{
		orderSlice = orderManagement.PrioritizeOrder(<- orderChan, orderSlice)
		fmt.Println(orderSlice)	
	}

}