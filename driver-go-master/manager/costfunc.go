package manager

import (
	"Driver-go/elevio"
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

func calculateCost(e singleElevator.Elevator, order elevio.ButtonEvent) float64 {
	// Determine current location of elevator and direction
	currFloor := e.Floor
	currDir := e.Direction

	// Calculate distance to requested floor
	distance := math.Abs(float64(currFloor - order.Floor))

	// Calculate cost based on distance
	cost := distance * ratePerUnitDistance
	if e.Behaviour == singleElevator.EB_Idle {
		return cost

	} else {
		if currDir != elevio.MD_Stop && currDir != getDirection(currFloor, order.Floor) {
			cost += directionChangeCost
		}

		if (currDir == elevio.MD_Up && order.Button == elevio.BT_HallDown) ||
			(currDir == elevio.MD_Down && order.Button == elevio.BT_HallUp) {
			cost += buttonChangeCost
		}

		// Add any additional costs
		cost += waitingTimeCost(e)
	}

	return cost
}

// Helper function to calculate direction to travel
func getDirection(fromFloor, toFloor int) elevio.MotorDirection {
	if fromFloor < toFloor {
		return elevio.MD_Up
	} else if fromFloor > toFloor {
		return elevio.MD_Down
	}
	return elevio.MD_Stop
}

// Helper function to calculate waiting time cost
func waitingTimeCost(e singleElevator.Elevator) float64 {
	if e.Behaviour == singleElevator.EB_DoorOpen {
		return waitingTime * waitingTimeRate
	}
	return 0
}
