package manager

import (
	"Driver-go/elevio"
	"Driver-go/singleElevator"
)

type OrderStruct struct {
	ElevatorID    string
	OrderedButton elevio.ButtonEvent
	PanelPair     singleElevator.OrderpanelPair
}

func MakeOrder(panelPair singleElevator.OrderpanelPair, button elevio.ButtonEvent) OrderStruct {
	order := OrderStruct{
		ElevatorID:    singleElevator.MyID,
		OrderedButton: button,
		PanelPair:     panelPair}

	return order
}

func AssignOrderToElevator(database ElevatorDatabase, order elevio.ButtonEvent) string {

	lowCost := 100000.0
	lowestCostElevator := ""

	elevatorList := database.ElevatorList

	if shouldITakeTheOrder(order) {
		lowestCostElevator = singleElevator.MyID
	} else {
		for i := 0; i < len(elevatorList); i++ {
			c := calculateCost(elevatorList[i], order)

			if c < lowCost && elevatorList[i].Operating == singleElevator.WS_Connected {
				lowCost = c
				lowestCostElevator = elevatorList[i].ElevatorID
			} else if c == lowCost && elevatorList[i].Operating == singleElevator.WS_Connected {

				var temp = database.ElevatorList[i].ElevatorID
				if temp < lowestCostElevator {
					lowCost = c
					lowestCostElevator = elevatorList[i].ElevatorID
				}
			}

		}

	}

	return lowestCostElevator
}

func FindDeadOrders(database ElevatorDatabase, deadElevatorID string) []elevio.ButtonEvent {
	deadElevator := GetElevatorFromID(database, deadElevatorID)
	var deadOrders []elevio.ButtonEvent
	var order elevio.ButtonEvent

	for floor := 0; floor < singleElevator.NumFloors; floor++ {
		for button := elevio.BT_HallUp; button < elevio.BT_Cab; button++ {
			ownerOfOrder := deadElevator.Requests[floor][button].ElevatorID
			order.Button = elevio.ButtonType(button)
			order.Floor = floor

			if ownerOfOrder == deadElevatorID {
				deadOrders = append(deadOrders, order)
			}
		}

	}
	return deadOrders
}

/*func FinLowestCostElevator(elevatorList []singleElevator.Elevator, order elevio.ButtonEvent) singleElevator.Elevator{

}*/

func FindCabCallsForElevator(database ElevatorDatabase, newPeer string) []OrderStruct {
	var cabsToBeSent []OrderStruct
	for i := 0; i < len(database.ElevatorList); i++ {
		if database.ElevatorList[i].ElevatorID == newPeer && newPeer != singleElevator.MyID {
			for floor := 0; floor < singleElevator.NumFloors; floor++ {
				if database.ElevatorList[i].Requests[floor][elevio.BT_Cab].ElevatorID == newPeer {
					var button elevio.ButtonEvent
					button.Floor = floor
					button.Button = elevio.BT_Cab
					panelPair := singleElevator.OrderpanelPair{ElevatorID: newPeer, OrderState: singleElevator.SO_Confirmed}
					cabsToBeSent = append(cabsToBeSent, MakeOrder(panelPair, button))
				}
			}
		}
	}
	return cabsToBeSent
}
