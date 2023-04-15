package singleElevator

import "Driver-go/elevio"

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
