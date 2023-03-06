package elevator

import "Driver-go/elevio"

type OrderMessageStruct struct {
	SystemID   string
	MessageID  string
	ElevatorID string

	//orderCounter int

	OrderedButton elevio.ButtonEvent

	ChosenElevator string
}

type IAmAliveMessageStruct struct {
	systemID   string
	messageID  string
	elevatorID string

	elevator Elevator
}
