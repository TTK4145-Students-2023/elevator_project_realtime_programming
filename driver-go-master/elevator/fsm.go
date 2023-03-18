package elevator

import (
	"Driver-go/elevio"
	"fmt"
	"time"
)

const MyID = "15657"

var elevator = Elevator_uninitialized(MyID)

func Fsm_init() {
	elevator = Elevator_uninitialized(MyID)

	elevio.SetFloorIndicator(elevator.Floor)
	SetAllLights(elevator)
}

func SetAllLights(es Elevator) {
	for floor := 0; floor < NumFloors; floor++ {
		for btn := elevio.BT_HallUp; btn < NumButtons; btn++ {
			elevio.SetButtonLamp(btn, floor, es.Requests[floor][btn].order)
			//Vurderte individuell sjekk på cab, men fordi caben kun er intern i arrayet, så må det være denne heisens cab uansett
		}
	}
}

func Fsm_onInitBetweenFloors() {
	elevio.SetMotorDirection(elevio.MD_Down)
	elevator.Dirn = elevio.MD_Down
	elevator.Behaviour = EB_Moving
}

func Fsm_onRequestButtonPress(btnFloor int, btnType elevio.ButtonType, chosenElevator string, timer *time.Timer) {
	//fmt.Printf("\n\n%s(%d, %s)\n", "fsm_onRequestButtonPress", btnFloor, btnType.ToString())
	//ElevatorPrint(elevator)
	//fmt.Println(calculateCost(&elevator, btnFloor))
	//La til denne for å sikre at man ikke omfordeler en ordre dersom en knapp blir trykket på flere ganger
	switch elevator.Behaviour {
	case EB_DoorOpen:
		if Requests_shouldClearImmediately(elevator, btnFloor, btnType) {
			timer.Reset(3 * time.Second)
			//fmt.Println("Her kan vi kjøre clearOnFloor()")
		} else {
			elevator.Requests[btnFloor][btnType].order = true
			elevator.Requests[btnFloor][btnType].ElevatorID = chosenElevator
		}
	case EB_Moving:
		elevator.Requests[btnFloor][btnType].order = true
		elevator.Requests[btnFloor][btnType].ElevatorID = chosenElevator
	case EB_Idle:
		elevator.Requests[btnFloor][btnType].order = true
		elevator.Requests[btnFloor][btnType].ElevatorID = chosenElevator
		pair := Requests_chooseDirection(elevator)
		elevator.Dirn = pair.dirn
		elevator.Behaviour = pair.behaviour
		switch pair.behaviour {
		case EB_DoorOpen:
			elevio.SetDoorOpenLamp(true)
			timer.Reset(3 * time.Second)
			elevator = Requests_clearAtCurrentFloor(elevator)
		case EB_Moving:
			elevio.SetMotorDirection(elevator.Dirn)
		case EB_Idle:
		}

	}

	SetAllLights(elevator)

	fmt.Printf("\nNew state:\n")
	//ElevatorPrint(elevator)
}

func Fsm_onFloorArrival(newFloor int, timer *time.Timer) {
	//fmt.Printf("\n\n%s(%d)\n", "fsm_onFloorArrival", newFloor)
	//ElevatorPrint(elevator)

	elevator.Floor = newFloor
	elevio.SetFloorIndicator(newFloor)

	switch elevator.Behaviour {
	case EB_Moving:
		if Requests_shouldStop(elevator) {
			elevio.SetMotorDirection(elevio.MD_Stop)
			if Requests_here(elevator) {
				elevio.SetDoorOpenLamp(true)
				if !elevio.GetObstruction() {
					timer.Reset(3 * time.Second)
				}
				elevator.Behaviour = EB_DoorOpen
			} else {
				elevator.Behaviour = EB_Idle
			}
			elevator = Requests_clearAtCurrentFloor(elevator)
			SetAllLights(elevator)
		}
	default:
	}

	fmt.Printf("\nNew state:\n")
	//ElevatorPrint(elevator)
}

func Fsm_onDoorTimeout(timer *time.Timer) {
	//fmt.Printf("\n\n%s()\n", runtime.FuncForPC(reflect.ValueOf(fsm_onDoorTimeout).Pointer()).Name())
	//ElevatorPrint(elevator)

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

	//fmt.Println("\nNew state:")
	//ElevatorPrint(elevator)
}
