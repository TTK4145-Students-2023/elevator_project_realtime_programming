package singleElevator

import (
	"Driver-go/elevio"
	"fmt"
)

func Elevator_uninitialized(myID string) Elevator {
	elevator := Elevator{Floor: -10}
	elevator.Behaviour = EB_Idle
	elevator.Direction = elevio.MD_Stop
	elevator.ElevatorID = myID
	elevator.Operating = WS_Unconnected
	elevator.OrderNumber = 0
	elevator.IsAlone = true

	return elevator
}

func SetAllLights(es Elevator) {
	for floor := 0; floor < NumFloors; floor++ {
		for btn := elevio.BT_HallUp; btn < NumButtons; btn++ {
			if es.Requests[floor][btn].OrderState == SO_Confirmed {
				elevio.SetButtonLamp(btn, floor, true)
			} else {
				elevio.SetButtonLamp(btn, floor, false)
			}
		}
	}
}

func Fsm_init() {
	elevatorObject = Elevator_uninitialized(MyID)

	elevio.SetFloorIndicator(elevatorObject.Floor)
	SetAllLights(elevatorObject)
}

func Fsm_onInitBetweenFloors() {
	elevio.SetMotorDirection(elevio.MD_Down)
	elevatorObject.Direction = elevio.MD_Down
	elevatorObject.Behaviour = EB_Moving
}

func GetSingleEleavtorObject() Elevator {
	return elevatorObject
}

func IsDoorOpen() bool {
	var doorOpen = false
	if elevatorObject.Behaviour == EB_DoorOpen {
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
	return (elevatorObject.Floor == floor) && (elevatorObject.Behaviour == EB_Idle)
}

func checkNoOrder(elevator Elevator, btn elevio.ButtonType) bool {
	return elevator.Requests[elevator.Floor][btn].OrderState == SO_NoOrder
}

func setNoOrder(e Elevator, floor int, buttonType elevio.ButtonType) Elevator {
	temp := e
	temp.Requests[floor][buttonType].OrderState = SO_NoOrder
	temp.Requests[floor][buttonType].ElevatorID = ""
	return temp
}

func setConfirmedOrder(e Elevator, floor int, buttonType elevio.ButtonType, chosenElevator string) Elevator {
	temp := e
	temp.Requests[floor][buttonType].OrderState = SO_Confirmed
	temp.Requests[floor][buttonType].ElevatorID = chosenElevator
	return temp
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
