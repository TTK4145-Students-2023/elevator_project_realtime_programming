package databaseHandler

import (
	"Driver-go/elevatorHardware"
	"Driver-go/singleElevator"
)

type OrderStruct struct {
	ElevatorID    string
	OrderedButton elevatorHardware.ButtonEvent
	PanelPair     singleElevator.OrderpanelPair
}

func MakeOrder(panelPair singleElevator.OrderpanelPair, button elevatorHardware.ButtonEvent) OrderStruct {
	order := OrderStruct{
		ElevatorID:    singleElevator.MyID,
		OrderedButton: button,
		PanelPair:     panelPair}

	return order
}

func AssignOrderToElevator(database ElevatorDatabase, order elevatorHardware.ButtonEvent) string {

	lowCost := 100000.0
	lowestCostElevator := ""

	elevatorList := database.ElevatorList

	if shouldITakeTheOrder(order) {
		lowestCostElevator = singleElevator.MyID
	} else {
		for i := 0; i < len(elevatorList); i++ {
			c := calculateCost(elevatorList[i], order)

			if c < lowCost && elevatorList[i].Operating == singleElevator.Connected {
				lowCost = c
				lowestCostElevator = elevatorList[i].ElevatorID
			} else if c == lowCost && elevatorList[i].Operating == singleElevator.Connected {

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

func shouldITakeTheOrder(order elevatorHardware.ButtonEvent) bool {
	if order.Button == elevatorHardware.BT_Cab || singleElevator.GetIsAlone() || singleElevator.AvailableAtCurrFloor(order.Floor) {
		return true
	} else {
		return false
	}
}
