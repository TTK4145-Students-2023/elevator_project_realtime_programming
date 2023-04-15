package singleElevator

import (
	"Driver-go/elevatorHardware"
	"time"
)

const MyID = "17000"

var elevatorObject = Elevator_uninitialized(MyID)

func FloorArrival(newFloor int, doorTimer *time.Timer, immobilityTimer *time.Timer) Elevator {
	elevatorObject.Floor = newFloor
	elevatorHardware.SetFloorIndicator(newFloor)
	SetWorkingState(Connected)

	switch elevatorObject.Behaviour {
	case Moving:
		immobilityTimer.Reset(3 * time.Second)

		if Requests_shouldStop(elevatorObject) {
			immobilityTimer.Stop()
			elevatorHardware.SetMotorDirection(elevatorHardware.MD_Stop)
			elevatorHardware.SetDoorOpenLamp(true)

			if !elevatorHardware.GetObstruction() {
				doorTimer.Reset(3 * time.Second)
			}

			elevatorObject = Requests_clearAtCurrentFloor(elevatorObject)
			SetAllLights(elevatorObject)
			elevatorObject.Behaviour = DoorOpen
		}
	default:
	}
	return elevatorObject
}

func DoorTimeout(doorTimer *time.Timer) {
	switch elevatorObject.Behaviour {
	case DoorOpen:
		pair := Requests_chooseDirection(elevatorObject)
		elevatorObject.Direction = pair.direction
		elevatorObject.Behaviour = pair.behaviour

		switch elevatorObject.Behaviour {
		case DoorOpen:
			elevatorObject = Requests_clearAtCurrentFloor(elevatorObject)
			doorTimer.Reset(3 * time.Second)
			SetAllLights(elevatorObject)
		case Moving, Idle:
			elevatorHardware.SetDoorOpenLamp(false)
			elevatorHardware.SetMotorDirection(elevatorObject.Direction)
		}
	default:
		break
	}
}
