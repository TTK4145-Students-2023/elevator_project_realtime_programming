package databaseHandler

import (
	"Driver-go/singleElevator"
)

type ElevatorDatabase struct {
	ConnectedElevators int
	ElevatorList       []singleElevator.Elevator
}
