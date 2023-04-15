package singleElevator

import (
	"Driver-go/elevatorHardware"
	"fmt"
)

func Elevator_uninitialized(myID string) Elevator {
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
			if es.Requests[floor][btn].OrderState == Confirmed {
				elevatorHardware.SetButtonLamp(btn, floor, true)
			} else {
				elevatorHardware.SetButtonLamp(btn, floor, false)
			}
		}
	}
}

func Fsm_init() {
	elevatorObject = Elevator_uninitialized(MyID)

	elevatorHardware.SetFloorIndicator(elevatorObject.Floor)
	SetAllLights(elevatorObject)
}

func Fsm_onInitBetweenFloors() {
	elevatorHardware.SetMotorDirection(elevatorHardware.MD_Down)
	elevatorObject.Direction = elevatorHardware.MD_Down
	elevatorObject.Behaviour = Moving
}

func GetSingleEleavtorObject() Elevator {
	return elevatorObject
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

func checkNoOrder(elevator Elevator, btn elevatorHardware.ButtonType) bool {
	return elevator.Requests[elevator.Floor][btn].OrderState == NoOrder
}

func setNoOrder(e Elevator, floor int, buttonType elevatorHardware.ButtonType) Elevator {
	temp := e
	temp.Requests[floor][buttonType].OrderState = NoOrder
	temp.Requests[floor][buttonType].ElevatorID = ""
	return temp
}

func setConfirmedOrder(e Elevator, floor int, buttonType elevatorHardware.ButtonType, chosenElevator string) Elevator {
	temp := e
	temp.Requests[floor][buttonType].OrderState = Confirmed
	temp.Requests[floor][buttonType].ElevatorID = chosenElevator
	return temp
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
	fmt.Printf("  |ID = %-2d         |\n", es.ElevatorID)
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

func setLocalNewOrder(button elevatorHardware.ButtonEvent, chosenElevator string) Elevator {
	elevatorObject.Requests[button.Floor][button.Button].OrderState = NewOrder
	elevatorObject.Requests[button.Floor][button.Button].ElevatorID = chosenElevator
	return elevatorObject
}

func setLocalConfirmedOrder(button elevatorHardware.ButtonEvent, chosenElevator string) Elevator {
	elevatorObject.Requests[button.Floor][button.Button].OrderState = Confirmed
	elevatorObject.Requests[button.Floor][button.Button].ElevatorID = chosenElevator
	SetAllLights(elevatorObject)
	return elevatorObject
}
