package elevator

import (
	"Driver-go/elevio"
	"fmt"
	"time"
)

var elevator = Elevator_uninitialized()

func Fsm_init() Elevator {
	//elevator = Elevator_uninitialized()

	elevio.SetFloorIndicator(elevator.floor)
	SetAllLights(elevator)

	return elevator
}

func SetAllLights(es Elevator) {
	for floor := 0; floor < numFloors; floor++ {
		for btn := elevio.BT_HallUp; btn < numButtons; btn++ {
			elevio.SetButtonLamp(btn, floor, es.requests[floor][btn])
		}
	}
}

func Fsm_onInitBetweenFloors() {
	elevio.SetMotorDirection(elevio.MD_Down)
	elevator.dirn = elevio.MD_Down
	elevator.behaviour = EB_Moving
}

func Fsm_onRequestButtonPress(btnFloor int, btnType elevio.ButtonType) {
	//fmt.Printf("\n\n%s(%d, %s)\n", "fsm_onRequestButtonPress", btnFloor, btnType.ToString())
	elevatorPrint(elevator)

	switch elevator.behaviour {
	case EB_DoorOpen:
		if Requests_shouldClearImmediately(elevator, btnFloor, btnType) {
			elevator.Timer.Reset(time.Second * 3)
		} else {
			elevator.requests[btnFloor][btnType] = true
		}
	case EB_Moving:
		elevator.requests[btnFloor][btnType] = true
	case EB_Idle:
		elevator.requests[btnFloor][btnType] = true
		pair := Requests_chooseDirection(elevator)
		elevator.dirn = pair.dirn
		elevator.behaviour = pair.behaviour
		switch pair.behaviour {
		case EB_DoorOpen:
			elevio.SetDoorOpenLamp(true)
			elevator.Timer.Reset(time.Second * 3)
			elevator = Requests_clearAtCurrentFloor(elevator)
		case EB_Moving:
			elevio.SetMotorDirection(elevator.dirn)
		case EB_Idle:
		}
	}

	SetAllLights(elevator)

	fmt.Printf("\nNew state:\n")
	elevatorPrint(elevator)
}

func Fsm_onFloorArrival(newFloor int) {
	fmt.Printf("\n\n%s(%d)\n", "fsm_onFloorArrival", newFloor)
	elevatorPrint(elevator)

	elevator.floor = newFloor

	switch elevator.behaviour {
	case EB_Moving:
		if Requests_shouldStop(elevator) {
			elevio.SetMotorDirection(elevio.MD_Stop)
			elevio.SetDoorOpenLamp(true)
			elevator = Requests_clearAtCurrentFloor(elevator)
			//Timer_start(elevator.config.doorOpenDuration_s, &TimerActive)
			elevator.Timer.Reset(time.Second * 3)
			SetAllLights(elevator)
			elevator.behaviour = EB_DoorOpen
		}
	default:
	}

	fmt.Printf("\nNew state:\n")
	elevatorPrint(elevator)
}

func Fsm_onDoorTimeout() {
	//fmt.Printf("\n\n%s()\n", runtime.FuncForPC(reflect.ValueOf(fsm_onDoorTimeout).Pointer()).Name())
	elevatorPrint(elevator)

	pair := Requests_chooseDirection(elevator)
	elevator.dirn = pair.dirn
	elevator.behaviour = pair.behaviour

	switch elevator.behaviour {
	case EB_DoorOpen: //enten obstruction eller ikke
		//Timer_start(elevator.config.doorOpenDuration_s, &TimerActive)
		elevator.Timer.Reset(time.Second * 3)
		elevator = Requests_clearAtCurrentFloor(elevator)
		SetAllLights(elevator)
		fmt.Printf("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	case EB_Moving, EB_Idle:
		fmt.Printf("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
		elevio.SetDoorOpenLamp(false)
		elevio.SetMotorDirection(elevator.dirn)
	}

	fmt.Println("\nNew state:")
	elevatorPrint(elevator)
}
