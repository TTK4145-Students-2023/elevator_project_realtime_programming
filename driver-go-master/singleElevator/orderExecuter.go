package singleElevator

import (
	"Driver-go/elevio"
	"time"
)

type ElevatorStateUpdate struct {
	ElevatorID string
	Elevator   Elevator
}

func TransmittStateUpdate(stateUpdateTx chan ElevatorStateUpdate) {
	ElevatorUpdate := ElevatorStateUpdate{
		ElevatorID: elevatorObject.ElevatorID,
		Elevator:   elevatorObject}
	for {
		time.Sleep(200 * time.Millisecond)
		ElevatorUpdate.Elevator = elevatorObject
		stateUpdateTx <- ElevatorUpdate
	}
}




func HandleNewOrder(chosenElevator string, button elevio.ButtonEvent, doorTimer *time.Timer, immobilityTimer *time.Timer) Elevator {
	var newElevatorUpdate Elevator

	if chosenElevator == MyID {
		newElevatorUpdate = Fsm_onRequestButtonPress(button.Floor, button.Button, chosenElevator, doorTimer, immobilityTimer)
	} else {
		newElevatorUpdate = setLocalNewOrder(button, chosenElevator)
	}

	return newElevatorUpdate
}

func HandleConfirmedOrder(chosenElevator string, button elevio.ButtonEvent, doorTimer *time.Timer, immobilityTimer *time.Timer) Elevator {
	var newElevatorUpdate Elevator

	if chosenElevator == MyID {
		newElevatorUpdate = Fsm_onRequestButtonPress(button.Floor, button.Button, chosenElevator, doorTimer, immobilityTimer)
	} else {
		newElevatorUpdate = setLocalConfirmedOrder(button, chosenElevator)
	}

	return newElevatorUpdate
}

func setLocalNewOrder(button elevio.ButtonEvent, chosenElevator string) Elevator {
	elevatorObject.Requests[button.Floor][button.Button].OrderState = SO_NewOrder
	elevatorObject.Requests[button.Floor][button.Button].ElevatorID = chosenElevator
	return elevatorObject
}

func setLocalConfirmedOrder(button elevio.ButtonEvent, chosenElevator string) Elevator {
	elevatorObject.Requests[button.Floor][button.Button].OrderState = SO_Confirmed
	elevatorObject.Requests[button.Floor][button.Button].ElevatorID = chosenElevator
	SetAllLights(elevatorObject)
	return elevatorObject
}
