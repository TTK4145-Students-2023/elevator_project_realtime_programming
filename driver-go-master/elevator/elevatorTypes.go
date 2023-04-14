package elevator

import (
	"Driver-go/elevio"
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
	WS_Connected = iota
	WS_Unconnected
	WS_Immobile
)

type StateOfOrder int

const (
	SO_NoOrder = iota
	SO_NewOrder
	SO_Confirmed
)

type OrderpanelPair struct {
	OrderState StateOfOrder
	ElevatorID string
}

type Elevator struct {
	Floor          int
	ElevatorID     string
	Dirn           elevio.MotorDirection
	Requests       [NumFloors][NumButtons]OrderpanelPair
	Behaviour      ElevatorBehaviour
	DoorOpen       bool
	Operating      WorkingState
	SingleElevator bool
	OrderNumber    int
}
