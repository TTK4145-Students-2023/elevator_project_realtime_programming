package elevator

import (
	"Driver-go/elevio"
	"time"
)

type OrderMessageStruct struct {
	SystemID   string
	MessageID  string
	ElevatorID string

	//orderCounter int

	OrderedButton elevio.ButtonEvent

	ChosenElevator string
}

type FloorArrivalMessageStruct struct {
	SystemID   string
	MessageID  string
	ElevatorID string

	//orderCounter int

	ArrivedFloor int
}

type IAmAliveMessageStruct struct {
	SystemID   string
	MessageID  string
	ElevatorID string

	Elevator Elevator
}

// The example message. We just send one of these every second.
func SendIAmAlive(aliveTx chan IAmAliveMessageStruct) {
	aliveMsg := IAmAliveMessageStruct{SystemID: "Gruppe10",
		MessageID:  "Alive",
		ElevatorID: MyID,
		Elevator:   elevator}
	for {
		aliveMsg.Elevator = elevator //oppdaterer heismelding
		aliveTx <- aliveMsg
		time.Sleep(1 * time.Second)
	}
}
