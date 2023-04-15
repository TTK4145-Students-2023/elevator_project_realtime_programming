package databaseHandler

import (
	"Driver-go/elevatorHardware"
	"Driver-go/singleElevator"
	"math"
)

// Constants
const (
	baseCost            = 10.0
	ratePerUnitDistance = 0.5
	buttonChangeCost    = 50.0
	directionChangeCost = 5.0
	waitingTime         = 10.0
	waitingTimeRate     = 0.1
)

func calculateCost(e singleElevator.Elevator, order elevatorHardware.ButtonEvent) float64 {
	// Determine current location of elevator and direction
	currFloor := e.Floor
	currDir := e.Direction

	// Calculate distance to requested floor
	distance := math.Abs(float64(currFloor - order.Floor))

	// Calculate cost based on distance
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

		// Add any additional costs
		cost += waitingTimeCost(e)
	}

	return cost
}

// Helper function to calculate direction to travel
func getDirection(fromFloor, toFloor int) elevatorHardware.MotorDirection {
	if fromFloor < toFloor {
		return elevatorHardware.MD_Up
	} else if fromFloor > toFloor {
		return elevatorHardware.MD_Down
	}
	return elevatorHardware.MD_Stop
}

// Helper function to calculate waiting time cost
func waitingTimeCost(e singleElevator.Elevator) float64 {
	if e.Behaviour == singleElevator.DoorOpen {
		return waitingTime * waitingTimeRate
	}
	return 0
}
