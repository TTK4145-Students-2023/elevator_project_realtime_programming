package singleElevator

import (
	"Driver-go/elevio"
	"time"
)

const MyID = "16000"

var elevatorObject = Elevator_uninitialized(MyID)

func Fsm_onRequestButtonPress(btnFloor int, btnType elevio.ButtonType, chosenElevator string, doorTimer *time.Timer, immobilityTimer *time.Timer) Elevator {

	switch elevatorObject.Behaviour {
	case EB_DoorOpen:
		if Requests_shouldClearImmediately(elevatorObject, btnFloor, btnType) {
			elevatorObject = Requests_clearAtCurrentFloor(elevatorObject)
			elevatorObject = Requests_clearOnFloor(elevatorObject.ElevatorID, elevatorObject.Floor)
			doorTimer.Reset(3 * time.Second)
		} else {
			elevatorObject = setConfirmedOrder(elevatorObject, btnFloor, btnType, chosenElevator)
		}
	case EB_Moving:
		immobilityTimer.Reset(3 * time.Second)
		elevatorObject = setConfirmedOrder(elevatorObject, btnFloor, btnType, chosenElevator)

	case EB_Idle:
		elevatorObject = setConfirmedOrder(elevatorObject, btnFloor, btnType, chosenElevator)

		directionBehaviourPair := Requests_chooseDirection(elevatorObject)
		elevatorObject.Direction = directionBehaviourPair.direction
		elevatorObject.Behaviour = directionBehaviourPair.behaviour

		switch directionBehaviourPair.behaviour {
		case EB_DoorOpen:
			elevio.SetDoorOpenLamp(true)
			doorTimer.Reset(3 * time.Second)
			elevatorObject = Requests_clearAtCurrentFloor(elevatorObject)

		case EB_Moving:
			immobilityTimer.Reset(3 * time.Second)
			elevio.SetMotorDirection(elevatorObject.Direction)

		case EB_Idle:
		}

	}
	SetAllLights(elevatorObject)
	return elevatorObject
}

func Fsm_onFloorArrival(newFloor int, doorTimer *time.Timer, immobilityTimer *time.Timer) Elevator {
	elevatorObject.Floor = newFloor
	elevio.SetFloorIndicator(newFloor)
	SetWorkingState(WS_Connected)

	switch elevatorObject.Behaviour {
	case EB_Moving:
		immobilityTimer.Reset(3 * time.Second)

		if Requests_shouldStop(elevatorObject) {
			immobilityTimer.Stop()
			elevio.SetMotorDirection(elevio.MD_Stop)
			elevio.SetDoorOpenLamp(true)

			if !elevio.GetObstruction() {
				doorTimer.Reset(3 * time.Second)
			}

			elevatorObject = Requests_clearAtCurrentFloor(elevatorObject)
			SetAllLights(elevatorObject)
			elevatorObject.Behaviour = EB_DoorOpen
		}
	default:
	}
	return elevatorObject
}

func Fsm_onDoorTimeout(timer *time.Timer) {
	switch elevatorObject.Behaviour {
	case EB_DoorOpen:
		pair := Requests_chooseDirection(elevatorObject)
		elevatorObject.Direction = pair.direction
		elevatorObject.Behaviour = pair.behaviour

		switch elevatorObject.Behaviour {
		case EB_DoorOpen:
			elevatorObject = Requests_clearAtCurrentFloor(elevatorObject)
			timer.Reset(3 * time.Second)
			SetAllLights(elevatorObject)
		case EB_Moving, EB_Idle:
			elevio.SetDoorOpenLamp(false)
			elevio.SetMotorDirection(elevatorObject.Direction)
		}
	default:
		break
	}
}
