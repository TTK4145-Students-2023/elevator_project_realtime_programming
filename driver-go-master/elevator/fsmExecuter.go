package elevator

import (
	"Driver-go/elevio"
	"time"
)

func HandleNewOrder(chosenElevator string, button elevio.ButtonEvent, doorTimer *time.Timer, immobilityTimer *time.Timer) Elevator {
	var newElevatorUpdate Elevator

	if chosenElevator == MyID { // || orderBroadcast.OrderedButton.Button != elevio.BT_Cab {
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
