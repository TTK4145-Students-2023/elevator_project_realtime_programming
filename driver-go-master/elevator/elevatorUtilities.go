package elevator

import "Driver-go/elevio"

func Elevator_uninitialized(myID string) Elevator {
	elev := Elevator{Floor: -10}
	elev.Behaviour = EB_Idle
	elev.Dirn = elevio.MD_Stop
	elev.ElevatorID = myID
	elev.Operating = WS_Unconnected
	elev.OrderNumber = 0
	elev.SingleElevator = true

	return elev
}

func Fsm_init() {
	elevator = Elevator_uninitialized(MyID)

	elevio.SetFloorIndicator(elevator.Floor)
	SetAllLights(elevator)
}

func Fsm_onInitBetweenFloors() {
	elevio.SetMotorDirection(elevio.MD_Down)
	elevator.Dirn = elevio.MD_Down
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
