package manager

import (
	"Driver-go/elevator"
	"Driver-go/elevio"
)

type MessageStruct struct {
	SenderID    string
	MessageType int

	OrderCounter   int
	OrderedButton  elevio.ButtonEvent
	ChosenElevator string
	ArrivedFloor   int
	MyElevator     elevator.Elevator
}

type AckMessageStruct struct {
	IDOfAckReciever string
	MessageNumber   int
}

func MakeFloorMessage(floor int) MessageStruct {
	floorMsg := MessageStruct{
		SenderID:       elevator.MyID,
		MessageType:    17,
		OrderedButton:  elevio.ButtonEvent{Floor: 0, Button: 0},
		ChosenElevator: "",
		ArrivedFloor:   floor,
		MyElevator:     elevator.GetSingleElevatorStruct()}

	return floorMsg
}

func MakeOrderMessage(chosenElevator string, button elevio.ButtonEvent) MessageStruct {
	orderMsg := MessageStruct{
		SenderID:       elevator.MyID,
		MessageType:    69,
		OrderedButton:  button,
		ChosenElevator: chosenElevator,
		ArrivedFloor:   -1,
		MyElevator:     elevator.GetSingleElevatorStruct()}

	return orderMsg
}

func MakeAckMessage() AckMessageStruct {
	ackMsg := AckMessageStruct{
		IDOfAckReciever: elevator.MyID,
		MessageNumber:   0}

	return ackMsg
}

func MakeNewElevator() MessageStruct{
	newElevator := MessageStruct{
		SenderID:       elevator.MyID,
		MessageType:    666,
		OrderedButton:  elevio.ButtonEvent{Floor: 0, Button: 0},
		ChosenElevator: "",
		ArrivedFloor:   -1,
		MyElevator:     elevator.GetSingleElevatorStruct()}

	return newElevator
}