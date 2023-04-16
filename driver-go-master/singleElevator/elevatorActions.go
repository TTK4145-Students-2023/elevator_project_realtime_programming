package singleElevator

import (
	"Driver-go/elevatorHardware"
	"time"
)

const MyID = "16000"

var elevatorObject = MakeUnitintializedElevator(MyID)

func InitializeElevatorBetweenFloors() {
	elevatorHardware.SetMotorDirection(elevatorHardware.MD_Down)
	elevatorObject.Direction = elevatorHardware.MD_Down
	elevatorObject.Behaviour = Moving
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

func FloorArrival(newFloor int, doorTimer *time.Timer, immobilityTimer *time.Timer) Elevator {
	elevatorObject.Floor = newFloor
	elevatorHardware.SetFloorIndicator(newFloor)
	SetWorkingState(Connected)
	SetAllLights(elevatorObject)

	switch elevatorObject.Behaviour {
	case Moving:
		immobilityTimer.Reset(3 * time.Second)

		if elevatorShouldStop(elevatorObject) {
			immobilityTimer.Stop()
			elevatorHardware.SetMotorDirection(elevatorHardware.MD_Stop)
			elevatorHardware.SetDoorOpenLamp(true)

			if !elevatorHardware.GetObstruction() {
				doorTimer.Reset(3 * time.Second)
			}

			elevatorObject = ClearOrderAtCurrentFloor(elevatorObject)
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
		pair := ordersChooseDirection(elevatorObject)
		elevatorObject.Direction = pair.direction
		elevatorObject.Behaviour = pair.behaviour

		switch elevatorObject.Behaviour {
		case DoorOpen:
			elevatorObject = ClearOrderAtCurrentFloor(elevatorObject)
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

func HandleObstruction(obstruction bool, doorTimer *time.Timer, immobilityTimer *time.Timer) {
	if IsDoorOpen() && obstruction {
		doorTimer.Stop()
		immobilityTimer.Reset(3 * time.Second)
	} else if !obstruction && IsDoorOpen() {
		immobilityTimer.Stop()
		SetWorkingState(Connected)
		doorTimer.Reset(3 * time.Second)
	}
}
