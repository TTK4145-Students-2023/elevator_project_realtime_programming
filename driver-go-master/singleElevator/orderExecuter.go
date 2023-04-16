package singleElevator

import (
	"Driver-go/elevatorHardware"
	"time"
)

func ExecuteAssignedOrder(button elevatorHardware.ButtonEvent, chosenElevator string, doorTimer *time.Timer, immobilityTimer *time.Timer) Elevator {

	switch elevatorObject.Behaviour {
	case DoorOpen:
		if OrderShouldClearImmediately(elevatorObject, button.Floor, button.Button) {
			elevatorObject = ClearOrderAtCurrentFloor(elevatorObject)
			elevatorObject = ClearOrderAtThisFloor(elevatorObject.ElevatorID, elevatorObject.Floor)
			doorTimer.Reset(3 * time.Second)
		} else {
			elevatorObject = setLocalConfirmedOrder(button, chosenElevator)
		}
	case Moving:
		immobilityTimer.Reset(3 * time.Second)
		elevatorObject = setLocalConfirmedOrder(button, chosenElevator)

	case Idle:
		elevatorObject = setLocalConfirmedOrder(button, chosenElevator)

		directionBehaviourPair := ordersChooseDirection(elevatorObject)
		elevatorObject.Direction = directionBehaviourPair.direction
		elevatorObject.Behaviour = directionBehaviourPair.behaviour

		switch directionBehaviourPair.behaviour {
		case DoorOpen:
			elevatorHardware.SetDoorOpenLamp(true)
			doorTimer.Reset(3 * time.Second)
			elevatorObject = ClearOrderAtCurrentFloor(elevatorObject)

		case Moving:
			immobilityTimer.Reset(3 * time.Second)
			elevatorHardware.SetMotorDirection(elevatorObject.Direction)

		case Idle:
		}

	}
	SetAllLights(elevatorObject)
	return elevatorObject
}

func HandleNewOrder(chosenElevator string, button elevatorHardware.ButtonEvent, doorTimer *time.Timer, immobilityTimer *time.Timer) Elevator {
	var newElevatorUpdate Elevator

	if chosenElevator == MyID {
		newElevatorUpdate = ExecuteAssignedOrder(button, chosenElevator, doorTimer, immobilityTimer)
	} else {
		newElevatorUpdate = setLocalNewOrder(button, chosenElevator)
	}

	return newElevatorUpdate
}

func HandleConfirmedOrder(chosenElevator string, button elevatorHardware.ButtonEvent, doorTimer *time.Timer, immobilityTimer *time.Timer) Elevator {
	var newElevatorUpdate Elevator

	if chosenElevator == MyID {
		newElevatorUpdate = ExecuteAssignedOrder(button, chosenElevator, doorTimer, immobilityTimer)
	} else {
		newElevatorUpdate = setLocalConfirmedOrder(button, chosenElevator) //endre navn mer deskrriptivt
	}

	return newElevatorUpdate
}

func setLocalNewOrder(button elevatorHardware.ButtonEvent, chosenElevator string) Elevator {
	elevatorObject.Requests[button.Floor][button.Button].OrderState = NewOrder
	elevatorObject.Requests[button.Floor][button.Button].ElevatorID = chosenElevator
	return elevatorObject
}

func setLocalNoOrder(floor int, buttonType elevatorHardware.ButtonType) Elevator {
	elevatorObject.Requests[floor][buttonType].OrderState = NoOrder
	elevatorObject.Requests[floor][buttonType].ElevatorID = ""
	return elevatorObject
}


func setLocalConfirmedOrder(button elevatorHardware.ButtonEvent, chosenElevator string) Elevator {
	elevatorObject.Requests[button.Floor][button.Button].OrderState = ConfirmedOrder
	elevatorObject.Requests[button.Floor][button.Button].ElevatorID = chosenElevator
	SetAllLights(elevatorObject)
	return elevatorObject
}
