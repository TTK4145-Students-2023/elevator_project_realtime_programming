package elevator

import (
	"Driver-go/elevio"
	"fmt"
	"time"
)

//var elevator = Elevator_uninitialized()

func Fsm_init(elevator *Elevator) *Elevator {
	elevator = Elevator_uninitialized()

	elevio.SetFloorIndicator(elevator.Floor)
	SetAllLights(elevator)

	return elevator
}

func SetAllLights(es *Elevator) {
	for floor := 0; floor < numFloors; floor++ {
		for btn := elevio.BT_HallUp; btn < numButtons; btn++ {
			elevio.SetButtonLamp(btn, floor, es.requests[floor][btn])
		}
	}
}

func Fsm_onInitBetweenFloors(elevator *Elevator) {
	elevio.SetMotorDirection(elevio.MD_Down)
	elevator.dirn = elevio.MD_Down
	elevator.behaviour = EB_Moving
}

func Fsm_onRequestButtonPress(btnFloor int, btnType elevio.ButtonType, elevator *Elevator) {
	//fmt.Printf("\n\n%s(%d, %s)\n", "fsm_onRequestButtonPress", btnFloor, btnType.ToString())
	elevatorPrint(*elevator)

	switch elevator.behaviour {
	case EB_DoorOpen:
		if Requests_shouldClearImmediately(*elevator, btnFloor, btnType) {
			elevator.Timer.Reset(time.Second * 3)
		} else {
			elevator.requests[btnFloor][btnType] = true
			fmt.Printf("INNE I DOOR OPEN REQ BUTTON")
			//Fsm_onDoorTimeout(elevator)
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
			*elevator = Requests_clearAtCurrentFloor(elevator)
		case EB_Moving:
			elevio.SetMotorDirection(elevator.dirn)
		case EB_Idle:
		}
	}

	SetAllLights(elevator)

	fmt.Printf("\nNew state:\n")
	elevatorPrint(*elevator)
}

func Fsm_onFloorArrival(newFloor int, elevator *Elevator) {
	fmt.Printf("\n\n%s(%d)\n", "fsm_onFloorArrival", newFloor)
	elevatorPrint(*elevator)

	elevator.Floor = newFloor

	switch elevator.behaviour {
	case EB_Moving:
		if Requests_shouldStop(elevator) {
			elevio.SetMotorDirection(elevio.MD_Stop)
			elevio.SetDoorOpenLamp(true)
			*elevator = Requests_clearAtCurrentFloor(elevator)
			fmt.Printf("INNE I SHOULD STOP")
			//Timer_start(elevator.config.doorOpenDuration_s, &TimerActive)
			//Timer_doorOpen(&elevator.Timer)
			elevator.Timer.Reset(time.Second * 3)
			<-elevator.Timer.C
			fmt.Printf("etter timer/on floor arrival")
			SetAllLights(elevator)
			elevator.behaviour = EB_DoorOpen
			fmt.Print("etter satt door open")
		}
	default:
	}

	fmt.Printf("\nNew state:\n")
	elevatorPrint(*elevator)
}

func Fsm_onDoorTimeout(elevator *Elevator) {
	//fmt.Printf("\n\n%s()\n", runtime.FuncForPC(reflect.ValueOf(fsm_onDoorTimeout).Pointer()).Name())
	elevatorPrint(*elevator)
	fmt.Printf("INNE I DOOR OPEN")
	switch elevator.behaviour {
	case EB_DoorOpen:
		pair := Requests_chooseDirection(elevator)
		elevator.dirn = pair.dirn
		elevator.behaviour = pair.behaviour

		switch elevator.behaviour {
		case EB_DoorOpen:
			elevator.Timer.Reset(time.Second * 3)
			*elevator = Requests_clearAtCurrentFloor(elevator)
			SetAllLights(elevator)
		case EB_Moving, EB_Idle:
			elevio.SetDoorOpenLamp(false)
			elevio.SetMotorDirection(elevator.dirn)
		}
	default:
		break
	}

	fmt.Println("\nNew state:")
	elevatorPrint(*elevator)
}
