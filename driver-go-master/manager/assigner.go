package manager

import (
	"Driver-go/elevator"
	"Driver-go/elevio"
)

func AssignOrderToElevator(database ElevatorDatabase, order elevio.ButtonEvent) string {

	lowCost := 100000.0
	lowestCostElevator := ""

	elevatorList := database.ElevatorList

	if shouldITakeTheOrder(order) {
		lowestCostElevator = elevator.MyID
	} else {
		for i := 0; i < len(elevatorList); i++ {
			c := calculateCost(elevatorList[i], order)

			if c < lowCost && elevatorList[i].Operating == elevator.WS_Connected {
				lowCost = c
				lowestCostElevator = elevatorList[i].ElevatorID
			} else if c == lowCost && elevatorList[i].Operating == elevator.WS_Connected {

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

	for floor := 0; floor < elevator.NumFloors; floor++ {
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

/*func FinLowestCostElevator(elevatorList []elevator.Elevator, order elevio.ButtonEvent) elevator.Elevator{
	
}*/

func FindCabCallsForElevator(database ElevatorDatabase, newPeer string) []elevator.OrderStruct {
	var cabsToBeSent []elevator.OrderStruct
	for i := 0; i < len(database.ElevatorList); i++ {
		if database.ElevatorList[i].ElevatorID == newPeer && newPeer != elevator.MyID {
			for floor := 0; floor < elevator.NumFloors; floor++ {
				if database.ElevatorList[i].Requests[floor][elevio.BT_Cab].ElevatorID == newPeer {
					var button elevio.ButtonEvent
					button.Floor = floor
					button.Button = elevio.BT_Cab
					panelPair := elevator.OrderpanelPair{ElevatorID: newPeer, OrderState: elevator.SO_Confirmed}
					cabsToBeSent = append(cabsToBeSent, elevator.MakeOrder(panelPair, button))
				}
			}
		}
	}
	return cabsToBeSent
}
