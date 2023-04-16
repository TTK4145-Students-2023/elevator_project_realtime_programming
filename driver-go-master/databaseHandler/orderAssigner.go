package databaseHandler

import (
	"Driver-go/elevatorHardware"
	"Driver-go/singleElevator"
	"math"
)

const (
	baseCost            = 10.0
	ratePerUnitDistance = 0.5
	buttonChangeCost    = 50.0
	directionChangeCost = 5.0
	waitingTime         = 10.0
	waitingTimeRate     = 0.1
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


func calculateCost(e singleElevator.Elevator, order elevatorHardware.ButtonEvent) float64 {

	currFloor := e.Floor
	currDir := e.Direction

	distance := math.Abs(float64(currFloor - order.Floor))

	cost := distance * ratePerUnitDistance
	if e.Behaviour == singleElevator.Idle {
		return cost

	} else {
		if currDir != elevatorHardware.MD_Stop && currDir != getDirection(currFloor, order.Floor) {
			cost += directionChangeCost
		}

		if (currDir == elevatorHardware.MD_Up && order.Button == elevatorHardware.BT_HallDown) ||
			(currDir == elevatorHardware.MD_Down && order.Button == elevatorHardware.BT_HallUp) {
			cost += buttonChangeCost
		}

		cost += waitingTimeCost(e)
	}

	return cost
}


func getDirection(fromFloor, toFloor int) elevatorHardware.MotorDirection {
	if fromFloor < toFloor {
		return elevatorHardware.MD_Up
	} else if fromFloor > toFloor {
		return elevatorHardware.MD_Down
	}
	return elevatorHardware.MD_Stop
}


func waitingTimeCost(e singleElevator.Elevator) float64 {
	if e.Behaviour == singleElevator.DoorOpen {
		return waitingTime * waitingTimeRate
	}
	return 0
}


func shouldITakeTheOrder(order elevatorHardware.ButtonEvent) bool {
	if order.Button == elevatorHardware.BT_Cab || singleElevator.GetIsAlone() || singleElevator.AvailableAtCurrFloor(order.Floor) {
		return true
	} else {
		return false
	}
}
