package manager

import (
	"Driver-go/elevator"
	"Driver-go/elevio"
)

type ElevatorDatabase struct {
	NumElevators int
	ElevatorsInNetwork           [3]elevator.Elevator
}

func AssignOrderToElevator(elevators ElevatorDatabase, order elevio.ButtonEvent) string {
	lowCost := 100000.0
	elevatorID := ""

	for i := 0; i < elevators.NumElevators; i++ {
		c := calculateCost(&elevators.ElevatorsInNetwork[i], order) //OBS! Blanding av pekere og ikke pekere
		if c < lowCost {
			lowCost = c
			elevatorID = elevators.ElevatorsInNetwork[i].ElevatorID
		}
	}

	return elevatorID
}
