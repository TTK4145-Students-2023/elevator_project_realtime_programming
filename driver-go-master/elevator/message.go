package elevator

import (
	"Driver-go/elevio"
	"time"
)

type OrderStruct struct {
	ElevatorID    string
	OrderedButton elevio.ButtonEvent
	PanelPair     OrderpanelPair
}

type StateUpdateStruct struct {
	ElevatorID string
	Elevator   Elevator
}

func MakeOrder(panelPair OrderpanelPair, button elevio.ButtonEvent) OrderStruct {
	order := OrderStruct{
		ElevatorID:    MyID,
		OrderedButton: button,
		PanelPair:     panelPair}

	return order
}

func SendStateUpdate(aliveTx chan StateUpdateStruct) {
	aliveMsg := StateUpdateStruct{
		ElevatorID: MyID,
		Elevator:   elevator}
	for {
		time.Sleep(200 * time.Millisecond)
		aliveMsg.Elevator = elevator //oppdaterer heismelding
		aliveTx <- aliveMsg
	}
}
