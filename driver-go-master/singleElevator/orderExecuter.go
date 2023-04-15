package singleElevator

import (
	"Driver-go/elevio"
	"time"
)

type ElevatorUpdateToDatabase struct {
	ElevatorID string
	Elevator   Elevator
}

func SendElevatorToDatabase(aliveTx chan ElevatorUpdateToDatabase) {
	ElevatorUpdate := ElevatorUpdateToDatabase{
		ElevatorID: elevatorObject.ElevatorID,
		Elevator:   elevatorObject}
	for {
		time.Sleep(200 * time.Millisecond)
		ElevatorUpdate.Elevator = elevatorObject
		aliveTx <- ElevatorUpdate
	}
}

func Fsm_setLocalNewOrder(button elevio.ButtonEvent, chosenElevator string) Elevator {
	elevatorObject.Requests[button.Floor][button.Button].OrderState = SO_NewOrder
	elevatorObject.Requests[button.Floor][button.Button].ElevatorID = chosenElevator
	return elevatorObject
}

func Fsm_setLocalConfirmedOrder(button elevio.ButtonEvent, chosenElevator string) Elevator {
	elevatorObject.Requests[button.Floor][button.Button].OrderState = SO_Confirmed
	elevatorObject.Requests[button.Floor][button.Button].ElevatorID = chosenElevator
	SetAllLights(elevatorObject)
	return elevatorObject
}

func Fsm_updateLocalRequests(updatedElevator Elevator) {
	elevatorObject.Requests = updatedElevator.Requests
	SetAllLights(elevatorObject)
}

func HandleNewOrder(chosenElevator string, button elevio.ButtonEvent, doorTimer *time.Timer, immobilityTimer *time.Timer) Elevator {
	var newElevatorUpdate Elevator

	if chosenElevator == MyID {
		newElevatorUpdate = Fsm_onRequestButtonPress(button.Floor, button.Button, chosenElevator, doorTimer, immobilityTimer)
	} else {
		newElevatorUpdate = Fsm_setLocalNewOrder(button, chosenElevator)
	}

	return newElevatorUpdate
}

func HandleConfirmedOrder(chosenElevator string, button elevio.ButtonEvent, doorTimer *time.Timer, immobilityTimer *time.Timer) Elevator {
	var newElevatorUpdate Elevator

	if chosenElevator == MyID {
		newElevatorUpdate = Fsm_onRequestButtonPress(button.Floor, button.Button, chosenElevator, doorTimer, immobilityTimer)
	} else {
		newElevatorUpdate = Fsm_setLocalConfirmedOrder(button, chosenElevator)
	}

	return newElevatorUpdate
}
