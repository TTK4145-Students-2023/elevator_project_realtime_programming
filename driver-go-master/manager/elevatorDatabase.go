package manager

import (
	"Driver-go/elevator"
	"Driver-go/elevio"
)

type ElevatorManager struct {
	numElevators int
	ar           [3]elevator.Elevator
}

func AssignOrderToElevator(elevators ElevatorManager, order elevio.ButtonEvent) string {
	elevatorID := ""

	for i := 0; i < elevators.numElevators; i++ {
		calculateCost(&elevators.ar[i], order)
	}

	return elevatorID
}
