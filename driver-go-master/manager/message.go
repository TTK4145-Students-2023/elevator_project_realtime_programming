package manager

import (
	"Driver-go/elevator"
	"Driver-go/elevio"
)

type FloorArrivalMessageStruct struct {
	SenderID    string
	MessageType int

	OrderCounter int
	ArrivedFloor int
	MyElevator   elevator.Elevator
}

type SingleElevatorMessageStruct struct {
	SenderID   string
	MyElevator elevator.Elevator
}

type OrderMessageStruct struct {
	SenderID    string
	MessageType int

	OrderCounter   int
	OrderedButton  elevio.ButtonEvent
	ChosenElevator string
	MyElevator     elevator.Elevator
}

type AckMessageStruct struct {
	IDOfAckReciever string
	MessageNumber   int
}

func MakeFloorMessage(floor int) FloorArrivalMessageStruct {
	floorMsg := FloorArrivalMessageStruct{
		SenderID:     elevator.MyID,
		MessageType:  17,
		ArrivedFloor: floor,
		MyElevator:   elevator.GetSingleElevatorStruct()}

	return floorMsg
}

func MakeOrderMessage(chosenElevator string, button elevio.ButtonEvent) OrderMessageStruct {
	orderMsg := OrderMessageStruct{
		SenderID:       elevator.MyID,
		MessageType:    69,
		OrderedButton:  button,
		ChosenElevator: chosenElevator,
		MyElevator:     elevator.GetSingleElevatorStruct()}

	return orderMsg
}

func MakeAckMessage() AckMessageStruct {
	ackMsg := AckMessageStruct{
		IDOfAckReciever: elevator.MyID,
		MessageNumber:   0}

	return ackMsg
}

func MakeNewElevator() SingleElevatorMessageStruct { //Brukte bare floorarival nå for å teste
	newElevator := SingleElevatorMessageStruct{
		SenderID:   elevator.MyID,
		MyElevator: elevator.GetSingleElevatorStruct()}

	return newElevator
}
