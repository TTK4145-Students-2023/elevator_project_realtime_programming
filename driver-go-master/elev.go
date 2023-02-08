package elev

import (
	"Driver-go/elevio"
)

const numFloors = 4
const numButtons = 4

type ElevatorBehaviour int

const (
	EB_Idle = iota
	EB_Moving
	EB_DoorOpen
)

type Elevator struct {
	floor     int
	dirn      elevio.MotorDirection
	requests  [numFloors][numButtons]int
	behaviour ElevatorBehaviour
}

func Elevator_uninitialized() Elevator {
	elev := Elevator{floor: -1}
	elev.behaviour = EB_Idle
	elev.dirn = elevio.MD_Stop

	return elev
}
