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
	MyElevator Elevator
}

type FloorArrivalMessageStruct struct {
	SystemID   string
	MessageID  string
	ElevatorID string

	//orderCounter int
	ArrivedFloor int
	MyElevator Elevator

}

type IAmAliveMessageStruct struct {
	SystemID   string
	MessageID  string
	ElevatorID string

	Elevator Elevator
}

func MakeFloorMessage(floor int) FloorArrivalMessageStruct {
	floorMsg := FloorArrivalMessageStruct{SystemID: "Gruppe10",
	MessageID:    "Floor",
	ElevatorID:   MyID,
	ArrivedFloor: floor,
	MyElevator: elevator}

	return floorMsg
}

func MakeOrderMessage(chosenElevator string, button elevio.ButtonEvent) OrderMessageStruct{
	orderMsg := OrderMessageStruct{SystemID: "Gruppe10",
		MessageID:      "Order",
		ElevatorID:     MyID,
		OrderedButton:  button,
		ChosenElevator: chosenElevator,
		MyElevator: elevator}

	return orderMsg
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
		time.Sleep(100 * time.Millisecond)
	}
}

func WaitForAck()
