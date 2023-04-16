package singleElevator

import (
	"Driver-go/elevatorHardware"
)


//Structs used for holding information about the elevator's physical position and working state. 
//Also contains the struct for holding information about the orders assigned elevator and order state.
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
	ConfirmedOrder
)

type DirectionBehaviourPair struct {
	direction elevatorHardware.MotorDirection
	behaviour ElevatorBehaviour
}


type StateAndChosenElevator struct {
	OrderState           StateOfOrder
	AssingedElevatorID   string
}

type Elevator struct {
	Floor      int
	ElevatorID string
	Direction  elevatorHardware.MotorDirection
	Requests   [NumFloors][NumButtons]StateAndChosenElevator
	Behaviour  ElevatorBehaviour
	Operating  WorkingState
	IsAlone    bool
}

