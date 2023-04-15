package elevator

import (
	"Driver-go/elevio"
	"time"
)

func Fsm_setLocalNewOrder(button elevio.ButtonEvent, chosenElevator string) Elevator {
	elevator.Requests[button.Floor][button.Button].OrderState = SO_NewOrder
	elevator.Requests[button.Floor][button.Button].ElevatorID = chosenElevator
	return elevator
}

func Fsm_setLocalConfirmedOrder(button elevio.ButtonEvent, chosenElevator string) Elevator {
	elevator.Requests[button.Floor][button.Button].OrderState = SO_Confirmed
	elevator.Requests[button.Floor][button.Button].ElevatorID = chosenElevator
	SetAllLights(elevator)
	return elevator
}

func Fsm_updateLocalRequests(updatedElevator Elevator) {
	elevator.Requests = updatedElevator.Requests
	SetAllLights(elevator)
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
