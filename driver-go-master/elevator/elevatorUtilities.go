package elevator

import "Driver-go/elevio"

func Elevator_uninitialized(myID string) Elevator {
	elev := Elevator{Floor: -10}
	elev.Behaviour = EB_Idle
	elev.Direction = elevio.MD_Stop
	elev.ElevatorID = myID
	elev.Operating = WS_Unconnected
	elev.OrderNumber = 0
	elev.SingleElevator = true

	return elev
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
	elevator = Elevator_uninitialized(MyID)

	elevio.SetFloorIndicator(elevator.Floor)
	SetAllLights(elevator)
}

func Fsm_onInitBetweenFloors() {
	elevio.SetMotorDirection(elevio.MD_Down)
	elevator.Direction = elevio.MD_Down
	elevator.Behaviour = EB_Moving
}

func GetSingleEleavtorStruct() Elevator {
	return elevator
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

func GetIAmAlone() bool {
	return elevator.SingleElevator
}
func SetIAmAlone(alone bool) {
	elevator.SingleElevator = alone
}

func SetWorkingState(state WorkingState) {
	elevator.Operating = state
}

func AvailableAtCurrFloor(floor int) bool {
	return (elevator.Floor == floor) && (elevator.Behaviour == EB_Idle)
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
