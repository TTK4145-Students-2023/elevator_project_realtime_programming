package elevator

import (
	"Driver-go/elevio"
	"time"
)

const MyID = "18000"

var elevator = Elevator_uninitialized(MyID)

func Fsm_onRequestButtonPress(btnFloor int, btnType elevio.ButtonType, chosenElevator string, doorTimer *time.Timer, immobilityTimer *time.Timer) Elevator {

	switch elevator.Behaviour {
	case EB_DoorOpen:
		if Requests_shouldClearImmediately(elevator, btnFloor, btnType) {
			elevator = Requests_clearAtCurrentFloor(elevator)
			elevator = Requests_clearOnFloor(elevator.ElevatorID, elevator.Floor)
			doorTimer.Reset(3 * time.Second)
		} else {
			elevator = setConfirmedOrder(elevator, btnFloor, btnType, chosenElevator)
		}
	case EB_Moving:
		immobilityTimer.Reset(3 * time.Second)
		elevator = setConfirmedOrder(elevator, btnFloor, btnType, chosenElevator)

	case EB_Idle:
		elevator = setConfirmedOrder(elevator, btnFloor, btnType, chosenElevator)

		directionBehaviourPair := Requests_chooseDirection(elevator)
		elevator.Direction = directionBehaviourPair.direction
		elevator.Behaviour = directionBehaviourPair.behaviour

		switch directionBehaviourPair.behaviour {
		case EB_DoorOpen:
			elevio.SetDoorOpenLamp(true)
			doorTimer.Reset(3 * time.Second)
			elevator = Requests_clearAtCurrentFloor(elevator)

		case EB_Moving:
			immobilityTimer.Reset(3 * time.Second)
			elevio.SetMotorDirection(elevator.Direction)

		case EB_Idle:
		}

	}
	SetAllLights(elevator)
	return elevator
}

func Fsm_onFloorArrival(newFloor int, doorTimer *time.Timer, immobilityTimer *time.Timer) Elevator {
	elevator.Floor = newFloor
	elevio.SetFloorIndicator(newFloor)
	SetWorkingState(WS_Connected)

	switch elevator.Behaviour {
	case EB_Moving:
		immobilityTimer.Reset(3 * time.Second)

		if Requests_shouldStop(elevator) {
			immobilityTimer.Stop()
			elevio.SetMotorDirection(elevio.MD_Stop)
			elevio.SetDoorOpenLamp(true)

			if !elevio.GetObstruction() {
				doorTimer.Reset(3 * time.Second)
			}

			elevator = Requests_clearAtCurrentFloor(elevator)
			SetAllLights(elevator)
			elevator.Behaviour = EB_DoorOpen
		}
	default:
	}
	return elevator
}

func Fsm_onDoorTimeout(timer *time.Timer) {
	switch elevator.Behaviour {
	case EB_DoorOpen:
		pair := Requests_chooseDirection(elevator)
		elevator.Direction = pair.direction
		elevator.Behaviour = pair.behaviour

		switch elevator.Behaviour {
		case EB_DoorOpen:
			elevator = Requests_clearAtCurrentFloor(elevator)
			timer.Reset(3 * time.Second)
			SetAllLights(elevator)
		case EB_Moving, EB_Idle:
			elevio.SetDoorOpenLamp(false)
			elevio.SetMotorDirection(elevator.Direction)
		}
	default:
		break
	}
}
