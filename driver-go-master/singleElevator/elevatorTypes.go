package singleElevator

import (
	"Driver-go/elevatorHardware"
)

const NumFloors = 4
const NumButtons = 3

type ElevatorBehaviour int

const (
	Idle = iota
	Moving
	DoorOpen
	Undefined
)

type WorkingState int

const (
	Connected = iota
	Unconnected
	Immobile
)

type StateOfOrder int

const (
	NoOrder = iota
	NewOrder
	Confirmed
)

type DirectionBehaviourPair struct {
	direction elevatorHardware.MotorDirection
	behaviour ElevatorBehaviour
}

type OrderpanelPair struct {
	OrderState StateOfOrder
	ElevatorID string
}

type Elevator struct {
	Floor      int
	ElevatorID string
	Direction  elevatorHardware.MotorDirection
	Requests   [NumFloors][NumButtons]OrderpanelPair
	Behaviour  ElevatorBehaviour
	Operating  WorkingState
	IsAlone    bool
}

type ElevatorStateUpdate struct {
	ElevatorID string
	Elevator   Elevator
}
