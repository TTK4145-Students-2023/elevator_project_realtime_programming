package elevator

import (
	"Driver-go/elevio"
	"fmt"
)

const NumFloors = 4
const NumButtons = 3

type ElevatorBehaviour int

const (
	EB_Idle = iota
	EB_Moving
	EB_DoorOpen
	EB_Undefined
)

type WorkingState int

const (
	WS_Running = iota
	WS_Unconnected
	WS_NoMotor
)

type OrderpanelPair struct {
	order      bool
	ElevatorID string
}

type Elevator struct {
	Floor       int
	ElevatorID  string
	Dirn        elevio.MotorDirection
	Requests    [NumFloors][NumButtons]OrderpanelPair
	Behaviour   ElevatorBehaviour
	DoorOpen    bool
	Operating   WorkingState
	OrderNumber int
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

func ElevatorPrint(es Elevator) {
	fmt.Println("  +--------------------+")
	fmt.Printf("  |ID = %-2d         |\n", es.ElevatorID)
	fmt.Printf("  |floor = %-2d         |\n", es.Floor)
	fmt.Printf("  |dirn  = %-12.12s|\n", DirnToString(es.Dirn))
	fmt.Printf("  |behav = %-12.12s|\n", ebToString(es.Behaviour))
	fmt.Printf("  |door = %-2d          |\n", es.DoorOpen)
	fmt.Printf("  |operating =         |\n", es.Operating)
	fmt.Println("  +--------------------+")
	fmt.Println("  |  | up  | dn  | cab |")
	for f := NumFloors - 1; f >= 0; f-- {
		fmt.Printf("  | %d", f)
		for btn := 0; btn < NumButtons; btn++ {
			if (f == NumButtons-1 && btn == int(elevio.BT_HallUp)) ||
				(f == 0 && btn == int(elevio.BT_HallDown)) {
				fmt.Print("|    ")
			} else {
				fmt.Print(es.Requests[f][btn])
			}
		}
		fmt.Print("|\n")
	}
	fmt.Println("  +--------------------+")
}

func Elevator_uninitialized(myID string) Elevator {
	elev := Elevator{Floor: -10}
	elev.Behaviour = EB_Idle
	elev.Dirn = elevio.MD_Stop
	elev.ElevatorID = myID
	elev.Operating = WS_Unconnected
	elev.OrderNumber = 0

	return elev
}

func Elevator_increaseOrderNumber() {
	elevator.OrderNumber++
}

func IsDoorOpen() bool {
	var doorOpen = false
	if elevator.Behaviour == EB_DoorOpen {
		doorOpen = true
	}
	return doorOpen
}
