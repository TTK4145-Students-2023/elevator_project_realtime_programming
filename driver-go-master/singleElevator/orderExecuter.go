package singleElevator

import (
	"Driver-go/elevatorHardware"
	"time"
)

func TransmitStateUpdate(stateUpdateTx chan ElevatorStateUpdate) {
	ElevatorUpdate := ElevatorStateUpdate{
		ElevatorID: elevatorObject.ElevatorID,
		Elevator:   elevatorObject}
	for {
		time.Sleep(200 * time.Millisecond)
		ElevatorUpdate.Elevator = elevatorObject
		stateUpdateTx <- ElevatorUpdate
	}
}

func ExecuteAssignedOrder(btnFloor int, btnType elevatorHardware.ButtonType, chosenElevator string, doorTimer *time.Timer, immobilityTimer *time.Timer) Elevator {
	switch elevatorObject.Behaviour {
	case DoorOpen:
		if Requests_shouldClearImmediately(elevatorObject, btnFloor, btnType) {
			elevatorObject = Requests_clearAtCurrentFloor(elevatorObject)
			elevatorObject = Requests_clearOnFloor(elevatorObject.ElevatorID, elevatorObject.Floor)
			doorTimer.Reset(3 * time.Second)
		} else {
			elevatorObject = setConfirmedOrder(elevatorObject, btnFloor, btnType, chosenElevator)
		}
	case Moving:
		immobilityTimer.Reset(3 * time.Second)
		elevatorObject = setConfirmedOrder(elevatorObject, btnFloor, btnType, chosenElevator)

	case Idle:
		elevatorObject = setConfirmedOrder(elevatorObject, btnFloor, btnType, chosenElevator)

		directionBehaviourPair := Requests_chooseDirection(elevatorObject)
		elevatorObject.Direction = directionBehaviourPair.direction
		elevatorObject.Behaviour = directionBehaviourPair.behaviour

		switch directionBehaviourPair.behaviour {
		case DoorOpen:
			elevatorHardware.SetDoorOpenLamp(true)
			doorTimer.Reset(3 * time.Second)
			elevatorObject = Requests_clearAtCurrentFloor(elevatorObject)

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
		newElevatorUpdate = ExecuteAssignedOrder(button.Floor, button.Button, chosenElevator, doorTimer, immobilityTimer)
	} else {
		newElevatorUpdate = setLocalNewOrder(button, chosenElevator)
	}

	return newElevatorUpdate
}

func HandleConfirmedOrder(chosenElevator string, button elevatorHardware.ButtonEvent, doorTimer *time.Timer, immobilityTimer *time.Timer) Elevator {
	var newElevatorUpdate Elevator

	if chosenElevator == MyID {
		newElevatorUpdate = ExecuteAssignedOrder(button.Floor, button.Button, chosenElevator, doorTimer, immobilityTimer)
	} else {
		newElevatorUpdate = setLocalConfirmedOrder(button, chosenElevator) //endre navn mer deskrriptivt
	}

	return newElevatorUpdate
}
