package manager

import (
	"Driver-go/elevator"
	"Driver-go/elevio"
)

type ElevatorDatabase struct {
	/*NumElevators       int
	ElevatorsInNetwork [3]elevator.Elevator	*/

	Elevator1352 elevator.Elevator
	Elevator7031 elevator.Elevator

	//ElevatorsInNetwork [2]elevator.Elevator
}

func AssignOrderToElevator(elevators ElevatorDatabase, order elevio.ButtonEvent) string {
	lowCost := 100000.0
	elevatorID := ""

	elevatorsInNetwork := [2]elevator.Elevator{elevators.Elevator1352, elevators.Elevator7031}
 
	if order.Button == elevio.BT_Cab {
		elevatorID = elevator.MyID
	} else {
		for i := 0; i < 2; i++ {
			c := calculateCost(&elevatorsInNetwork[i], order) //OBS! Blanding av pekere og ikke pekere
			if c < lowCost {
				lowCost = c
				elevatorID = elevatorsInNetwork[i].ElevatorID
			}
		}
	}
 
	return elevatorID
 }
 
