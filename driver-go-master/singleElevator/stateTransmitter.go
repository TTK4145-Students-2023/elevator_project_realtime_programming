package singleElevator

import "time"

//Transmits the elevator's state to the network of elevators, for updating the database. 

type ElevatorStateUpdate struct {
	ElevatorID string
	Elevator   Elevator
}


func TransmitStateUpdate(stateUpdateChannelTx chan ElevatorStateUpdate) {
	ElevatorUpdate := ElevatorStateUpdate{
		ElevatorID: elevatorObject.ElevatorID,
		Elevator:   elevatorObject}
	for {
		time.Sleep(200 * time.Millisecond)
		ElevatorUpdate.Elevator = elevatorObject
		stateUpdateChannelTx <- ElevatorUpdate
	}
}
