package singleElevator

import (
	"Driver-go/elevatorHardware"
	"time"
)

// Holds the elevator object info.
const MyID = "16000"

var elevatorObject = MakeElevatorObject(MyID)

//Functions for performing elevator actions related to the input from the floor sensor, buttonpanel and order distributer.
//Also handles door operations and executing of assinged orders.

func InitializeIfElevatorBetweenFloors() {
	if elevatorHardware.GetFloor() == -1 {
		elevatorHardware.SetMotorDirection(elevatorHardware.MD_Down)
		elevatorObject.Direction = elevatorHardware.MD_Down
		elevatorObject.Behaviour = Moving
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

		if ElevatorShouldStop(elevatorObject) {
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

func ExecuteAssignedOrder(button elevatorHardware.ButtonEvent, chosenElevator string, doorTimer *time.Timer, immobilityTimer *time.Timer) Elevator {
	switch elevatorObject.Behaviour {
	case DoorOpen:
		if OrderShouldClearImmediately(elevatorObject, button.Floor, button.Button) {
			elevatorObject = ClearOrderAtCurrentFloor(elevatorObject)
			elevatorObject = ClearOrderAtThisFloor(elevatorObject.ElevatorID, elevatorObject.Floor)
			doorTimer.Reset(3 * time.Second)
		} else {
			elevatorObject = SaveLocalConfirmedOrder(button, chosenElevator)
		}
	case Moving:
		immobilityTimer.Reset(3 * time.Second)
		elevatorObject = SaveLocalConfirmedOrder(button, chosenElevator)

	case Idle:
		elevatorObject = SaveLocalConfirmedOrder(button, chosenElevator)

		directionBehaviourPair := OrdersChooseDirection(elevatorObject)
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

func DoorTimeout(doorTimer *time.Timer) {
	switch elevatorObject.Behaviour {
	case DoorOpen:
		pair := OrdersChooseDirection(elevatorObject)
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
	} else if IsDoorOpen() && !obstruction {
		immobilityTimer.Stop()
		SetWorkingState(Connected)
		doorTimer.Reset(3 * time.Second)
	}
}
