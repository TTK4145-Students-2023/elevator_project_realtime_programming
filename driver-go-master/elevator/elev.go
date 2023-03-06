package elevator

import (
	"Driver-go/elevio"
	"fmt"
)

const numFloors = 4
const numButtons = 3

type ElevatorBehaviour int

const (
	EB_Idle = iota
	EB_Moving
	EB_DoorOpen
)

type OrderpanelPair struct {
	order      bool
	elevatorID string
}

type Elevator struct {
	Floor      int
	ElevatorID string
	Dirn       elevio.MotorDirection
	requests   [numFloors][numButtons]OrderpanelPair
	Behaviour  ElevatorBehaviour
	DoorOpen   bool
}

func ebToString(eb ElevatorBehaviour) string {
	switch eb {
	case EB_Idle:
		return "EB_Idle"
	case EB_DoorOpen:
		return "EB_DoorOpen"
	case EB_Moving:
		return "EB_Moving"
	default:
		return "EB_UNDEFINED"
	}
}

func DirnToString(direction elevio.MotorDirection) string {
	switch direction {
	case elevio.MD_Up:
		return "MotorUp"
	case elevio.MD_Down:
		return "MotorDown"
	case elevio.MD_Stop:
		return "MotorStop"
	default:
		return "MotorUndefined"
	}
}

func elevatorPrint(es Elevator) {
	fmt.Println("  +--------------------+")
	fmt.Printf("  |floor = %-2d      	|\n", es.Floor)
	fmt.Printf("  |dirn  = %-12.12s|\n", DirnToString(es.Dirn))
	fmt.Printf("  |behav = %-12.12s|\n", ebToString(es.Behaviour))

	fmt.Println("  +--------------------+")
	fmt.Println("  |  | up  | dn  | cab |")
	for f := numFloors - 1; f >= 0; f-- {
		fmt.Printf("  | %d", f)
		for btn := 0; btn < numButtons; btn++ {
			if (f == numButtons-1 && btn == int(elevio.BT_HallUp)) ||
				(f == 0 && btn == int(elevio.BT_HallDown)) {
				fmt.Print("| 	")
			} else {
				fmt.Print(es.requests[f][btn])
			}
		}
		fmt.Print("|\n")
	}
	fmt.Println("  +--------------------+")
}

func Elevator_uninitialized(myID string) Elevator {
	elev := Elevator{Floor: -1}
	elev.Behaviour = EB_Idle
	elev.Dirn = elevio.MD_Stop
	elev.ElevatorID = myID

	return elev
}
