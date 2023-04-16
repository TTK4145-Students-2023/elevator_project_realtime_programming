package singleElevator

import (
	"Driver-go/elevatorHardware"
	"fmt"
)

//Set functions and get functions for the elevator. Also print utilities.

func MakeElevatorObject(myID string) Elevator {
	elevator := Elevator{Floor: -10}
	elevator.Behaviour = Idle
	elevator.Direction = elevatorHardware.MD_Stop
	elevator.ElevatorID = myID
	elevator.Operating = Unconnected
	elevator.IsAlone = true

	return elevator
}

func SetAllLights(es Elevator) {
	for floor := 0; floor < NumFloors; floor++ {
		for btn := elevatorHardware.BT_HallUp; btn < NumButtons; btn++ {
			if es.Requests[floor][btn].OrderState == ConfirmedOrder {
				elevatorHardware.SetButtonLamp(btn, floor, true)
			} else {
				elevatorHardware.SetButtonLamp(btn, floor, false)
			}
		}
	}
}



func IsDoorOpen() bool {
	var doorOpen = false
	if elevatorObject.Behaviour == DoorOpen {
		doorOpen = true
	}
	return doorOpen
}

func GetIsAlone() bool {
	return elevatorObject.IsAlone
}
func SetIsAlone(alone bool) {
	elevatorObject.IsAlone = alone
}

func SetWorkingState(state WorkingState) {
	elevatorObject.Operating = state
}

func AvailableAtCurrFloor(floor int) bool {
	return (elevatorObject.Floor == floor) && (elevatorObject.Behaviour == Idle)
}

func ebToString(eb ElevatorBehaviour) string {
	switch eb {
	case Idle:
		return "Idle"
	case DoorOpen:
		return "DoorOpen"
	case Moving:
		return "Moving"
	default:
		return "UNDEFINED"
	}
}

func DirectionToString(direction elevatorHardware.MotorDirection) string {
	switch direction {
	case elevatorHardware.MD_Up:
		return "MotorUp"
	case elevatorHardware.MD_Down:
		return "MotorDown"
	case elevatorHardware.MD_Stop:
		return "MotorStop"
	default:
		return "MotorUndefined"
	}
}

func ElevatorPrint(es Elevator) {
	fmt.Println("  +--------------------+")
	fmt.Printf("  |ID = %-12.12s         |\n", es.ElevatorID)
	fmt.Printf("  |floor = %-2d         |\n", es.Floor)
	fmt.Printf("  |Direction  = %-12.12s|\n", DirectionToString(es.Direction))
	fmt.Printf("  |behav = %-12.12s|\n", ebToString(es.Behaviour))
	//fmt.Printf("  |door = %-2d          |\n", es.DoorOpen)
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
