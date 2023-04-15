package singleElevator

import (
	"Driver-go/elevio"
	"fmt"
)

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

func DirectionToString(direction elevio.MotorDirection) string {
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
	fmt.Printf("  |Direction  = %-12.12s|\n", DirectionToString(es.Direction))
	fmt.Printf("  |behav = %-12.12s|\n", ebToString(es.Behaviour))
	fmt.Printf("  |door = %-2d          |\n", es.DoorOpen)
	fmt.Printf("  |operating = %-2d        |\n", es.Operating)
	fmt.Println("  +--------------------+")
	fmt.Println("  |  | up  | dn  | cab |")
	for f := NumFloors - 1; f >= 0; f-- {
		fmt.Printf("  | %d", f)
		for btn := 0; btn < NumButtons; btn++ {

			fmt.Print(es.Requests[f][btn])
		}
		fmt.Print("|\n")
	}
	fmt.Println("  +--------------------+")
}
