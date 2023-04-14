package elevator

import (
	"Driver-go/elevio"
	"time"
)

const MyID = "15657"

var elevator = Elevator_uninitialized(MyID)

func SetAllLights(es Elevator) {
	for floor := 0; floor < NumFloors; floor++ {
		for btn := elevio.BT_HallUp; btn < NumButtons; btn++ {
			if es.Requests[floor][btn].OrderState == SO_Confirmed {
				elevio.SetButtonLamp(btn, floor, true)
			} else {
				elevio.SetButtonLamp(btn, floor, false)
			}
		}
	}
}

func Fsm_onRequestButtonPress(btnFloor int, btnType elevio.ButtonType, chosenElevator string, doorTimer *time.Timer, immobilityTimer *time.Timer) Elevator {

	switch elevator.Behaviour {
	case EB_DoorOpen:
		if Requests_shouldClearImmediately(elevator, btnFloor, btnType) {
			elevator = Requests_clearAtCurrentFloor(elevator)
			elevator = Requests_clearOnFloor(elevator.ElevatorID, elevator.Floor)
			doorTimer.Reset(3 * time.Second)
		} else {
			elevator.Requests[btnFloor][btnType].OrderState = SO_Confirmed
			elevator.Requests[btnFloor][btnType].ElevatorID = chosenElevator
		}
	case EB_Moving:
		immobilityTimer.Reset(3 * time.Second)

		elevator.Requests[btnFloor][btnType].OrderState = SO_Confirmed
		elevator.Requests[btnFloor][btnType].ElevatorID = chosenElevator

	case EB_Idle:
		elevator.Requests[btnFloor][btnType].OrderState = SO_Confirmed
		elevator.Requests[btnFloor][btnType].ElevatorID = chosenElevator

		directionBehaviourPair := Requests_chooseDirection(elevator)
		elevator.Dirn = directionBehaviourPair.dirn
		elevator.Behaviour = directionBehaviourPair.behaviour

		switch directionBehaviourPair.behaviour {
		case EB_DoorOpen:
			elevio.SetDoorOpenLamp(true)
			doorTimer.Reset(3 * time.Second)
			elevator = Requests_clearAtCurrentFloor(elevator)

		case EB_Moving:
			immobilityTimer.Reset(3 * time.Second)
			elevio.SetMotorDirection(elevator.Dirn)

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
		elevator.Dirn = pair.dirn
		elevator.Behaviour = pair.behaviour

		switch elevator.Behaviour {
		case EB_DoorOpen:
			elevator = Requests_clearAtCurrentFloor(elevator)
			timer.Reset(3 * time.Second)
			SetAllLights(elevator)
		case EB_Moving, EB_Idle:
			elevio.SetDoorOpenLamp(false)
			elevio.SetMotorDirection(elevator.Dirn)
		}
	default:
		break
	}
}
