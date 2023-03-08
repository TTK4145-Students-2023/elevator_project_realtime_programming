package manager

import (
	"Driver-go/elevator"
	"Driver-go/elevio"
)

type ElevatorDatabase struct {
	NumElevators int
	//ElevatorsInNetwork []elevator.Elevator

	Elevator13520 elevator.Elevator
	Elevator70310 elevator.Elevator
	Elevator54321 elevator.Elevator

	//ElevatorsInNetwork [2]elevator.Elevator
}

func AssignOrderToElevator(database ElevatorDatabase, order elevio.ButtonEvent) string {

	lowCost := 100000.0
	elevatorID := ""

	connectedElevators := [3]elevator.Elevator{database.Elevator13520, database.Elevator70310, database.Elevator54321}

	if order.Button == elevio.BT_Cab {
		elevatorID = elevator.MyID
	} else {
		for i := 0; i < database.NumElevators; i++ {
			c := calculateCost(&connectedElevators[i], order)                          //OBS! Blanding av pekere og ikke pekere
			if c < lowCost && connectedElevators[i].Operating == elevator.WS_Running { //Sjekker at calgt heis ikke er unconnected
				lowCost = c
				elevatorID = connectedElevators[i].ElevatorID
			}
		}
	}

	return elevatorID
}

/*

func ElevatorsInNetwork(database ElevatorDatabase) []elevator.Elevator{
	var elevatorsInNetwork []elevator.Elevator

	i := 0

	if database.Elevator13520.Operating == elevator.WS_Running {
		elevatorsInNetwork[i] = database.Elevator13520
		i++
	}

	if database.Elevator70310.Operating == elevator.WS_Running{
		elevatorsInNetwork[i] = database.Elevator70310
		i++
	}

	if database.Elevator54321.Operating == elevator.WS_Running {
		elevatorsInNetwork[i] = database.Elevator54321
		i++
	}

	return elevatorsInNetwork
}
*/
