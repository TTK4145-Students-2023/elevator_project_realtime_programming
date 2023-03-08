package elevator

import (
	"Driver-go/elevio"
	"fmt"
	"time"
)

const MyID = "70310"

var elevator = Elevator_uninitialized(MyID)



func Fsm_init() {
	elevator = Elevator_uninitialized(MyID)

	elevio.SetFloorIndicator(elevator.Floor)
	SetAllLights(elevator)
}

func SetAllLights(es Elevator) {
	for floor := 0; floor < numFloors; floor++ {
		for btn := elevio.BT_HallUp; btn < numButtons; btn++ {
			elevio.SetButtonLamp(btn, floor, es.requests[floor][btn].order)
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
	ElevatorPrint(elevator)
	//fmt.Println(calculateCost(&elevator, btnFloor))
	if !elevator.requests[btnFloor][btnType].order { //La til denne for å sikre at man ikke omfordeler en ordre dersom en knapp blir trykket på flere ganger
		switch elevator.Behaviour {
		case EB_DoorOpen:
			if Requests_shouldClearImmediately(elevator, btnFloor, btnType) {
				timer.Reset(3 * time.Second)
			} else {
				elevator.requests[btnFloor][btnType].order = true
				elevator.requests[btnFloor][btnType].elevatorID = chosenElevator
			}
		case EB_Moving:
			elevator.requests[btnFloor][btnType].order = true
			elevator.requests[btnFloor][btnType].elevatorID = chosenElevator
		case EB_Idle:
			elevator.requests[btnFloor][btnType].order = true
			elevator.requests[btnFloor][btnType].elevatorID = chosenElevator
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
	}

	SetAllLights(elevator)

	fmt.Printf("\nNew state:\n")
	ElevatorPrint(elevator)
}

func Fsm_onFloorArrival(newFloor int, timer *time.Timer) {
	fmt.Printf("\n\n%s(%d)\n", "fsm_onFloorArrival", newFloor)
	ElevatorPrint(elevator)

	elevator.Floor = newFloor

	switch elevator.Behaviour {
	case EB_Moving:
		if Requests_shouldStop(elevator) {
			elevio.SetMotorDirection(elevio.MD_Stop)
			elevio.SetDoorOpenLamp(true)
			timer.Reset(3 * time.Second)
			elevator = Requests_clearAtCurrentFloor(elevator)
			SetAllLights(elevator)
			elevator.Behaviour = EB_DoorOpen
		}
	default:
	}

	fmt.Printf("\nNew state:\n")
	ElevatorPrint(elevator)
}

func Fsm_onDoorTimeout() {
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
			SetAllLights(elevator)
		case EB_Moving, EB_Idle:
			fmt.Printf("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
			elevio.SetDoorOpenLamp(false)
			elevio.SetMotorDirection(elevator.Dirn)
		}
	default:
		break
	}

	fmt.Println("\nNew state:")
	ElevatorPrint(elevator)
}
