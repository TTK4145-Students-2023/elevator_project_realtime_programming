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
			if es.Requests[floor][btn].OrderState == SO_Confirmed {
				elevio.SetButtonLamp(btn, floor, true)
			}
			//elevio.SetButtonLamp(btn, floor, es.Requests[floor][btn].order)
			//Vurderte individuell sjekk på cab, men fordi caben kun er intern i arrayet, så må det være denne heisens cab uansett
		}
	}
}

func Fsm_onInitBetweenFloors() {
	elevio.SetMotorDirection(elevio.MD_Down)
	elevator.Dirn = elevio.MD_Down
	elevator.Behaviour = EB_Moving
}

func Fsm_onRequestButtonPress(btnFloor int, btnType elevio.ButtonType, chosenElevator string, doorTimer *time.Timer, immobilityTimer *time.Timer) {
	//fmt.Printf("\n\n%s(%d, %s)\n", "fsm_onRequestButtonPress", btnFloor, btnType.ToString())
	ElevatorPrint(elevator)
	//fmt.Println(calculateCost(&elevator, btnFloor))
	fmt.Println("----inne i on req buttonpress------------")
	//La til denne for å sikre at man ikke omfordeler en ordre dersom en knapp blir trykket på flere ganger
	switch elevator.Behaviour {
	case EB_DoorOpen:
		if Requests_shouldClearImmediately(elevator, btnFloor, btnType) {
			doorTimer.Reset(3 * time.Second)
			fmt.Println("Her kan vi kjøre clearOnFloor()")
		} else {
			elevator.Requests[btnFloor][btnType].OrderState = SO_NewOrder
			elevator.Requests[btnFloor][btnType].ElevatorID = chosenElevator
		}
	case EB_Moving:
		immobilityTimer.Reset(3 * time.Second)
		fmt.Println("Nå har jeg resetet immobilityTimer i Fsm_Req, case EB_Moving_1")
		elevator.Requests[btnFloor][btnType].OrderState = SO_NewOrder
		elevator.Requests[btnFloor][btnType].ElevatorID = chosenElevator
	case EB_Idle:
		elevator.Requests[btnFloor][btnType].OrderState = SO_NewOrder
		elevator.Requests[btnFloor][btnType].ElevatorID = chosenElevator
		pair := Requests_chooseDirection(elevator)
		elevator.Dirn = pair.dirn
		elevator.Behaviour = pair.behaviour
		switch pair.behaviour {
		case EB_DoorOpen:
			elevio.SetDoorOpenLamp(true)
			doorTimer.Reset(3 * time.Second)
			elevator = Requests_clearAtCurrentFloor(elevator)
		case EB_Moving:
			immobilityTimer.Reset(3 * time.Second)
			fmt.Println("Nå har jeg resetet immobilityTimer i Fsm_Req, case EB_Moving_2")
			elevio.SetMotorDirection(elevator.Dirn)
		case EB_Idle:
		}

	}

	SetAllLights(elevator)

	fmt.Printf("\nNew state:\n")
	ElevatorPrint(elevator)
}

func Fsm_onFloorArrival(newFloor int, doorTimer *time.Timer, immobilityTimer *time.Timer) {
	fmt.Printf("\n\n%s(%d)\n", "fsm_onFloorArrival", newFloor)
	ElevatorPrint(elevator)

	elevator.Floor = newFloor
	elevio.SetFloorIndicator(newFloor)
	SetWorkingState(WS_Connected)
	switch elevator.Behaviour {
	case EB_Moving:
		immobilityTimer.Reset(3 * time.Second)
		fmt.Println("Nå har jeg resetet immobilityTimer i Fsm_FloorA, case EB_Moving")
		if Requests_shouldStop(elevator) {
			immobilityTimer.Stop()
			fmt.Println("Stoppet immobilityTimer i fsm_floorA")
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

	fmt.Printf("\nNew state:\n")
	ElevatorPrint(elevator)
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

func Fsm_updateQueue(updatedElevator Elevator) {
	elevator.Requests = updatedElevator.Requests
}
