package manager

import (
	"Driver-go/elevator"
	"Driver-go/elevio"
	"Driver-go/network/peers"
	"fmt"
)

type ElevatorDatabase struct {
	NumElevators       int
	ElevatorsInNetwork []elevator.Elevator
}

func AssignOrderToElevator(database ElevatorDatabase, order elevio.ButtonEvent) string {

	lowCost := 100000.0
	elevatorID := ""

	connectedElevators := database.ElevatorsInNetwork

	if order.Button == elevio.BT_Cab {
		elevatorID = elevator.MyID
	} else {
		for i := 0; i < database.NumElevators; i++ {
			c := calculateCost(&connectedElevators[i], order)                          //OBS! Blanding av pekere og ikke pekere
			if c < lowCost && connectedElevators[i].Operating == elevator.WS_Running { //Sjekker at calgt heis ikke er unconnected
				lowCost = c
				elevatorID = connectedElevators[i].ElevatorID
			}
		}
	}

	return elevatorID
}

func ReassignDeadOrders(orderTx chan elevator.OrderMessageStruct, database ElevatorDatabase, deadElevatorID string) {
	deadElev := GetElevatorFromID(database, deadElevatorID)
	fmt.Println(" -----dead elevator id -----")
	fmt.Println(deadElev.ElevatorID)
	fmt.Println(("here are the orders"))
	for floor := 0; floor < elevator.NumFloors; floor++ {
		for button := 0; button < elevator.NumButtons; button++ {
			var order elevio.ButtonEvent
			order.Button = elevio.ButtonType(button)
			order.Floor = floor
			fmt.Println(deadElev.Requests[floor][button])

			if deadElev.Requests[floor][button].ElevatorID == deadElevatorID {
				fmt.Println("--------------FOUND DEADORDER--------------------------")
				SendOrderMessage(orderTx, order, database)
			}
		}

	}
	fmt.Println("-----------------REASSIGNED-----------------")
	elevator.ElevatorPrint(GetElevatorFromID(database, elevator.MyID))
}

func IsElevatorInDatabase(elevatorID string, database ElevatorDatabase) bool {
	for i := 0; i < database.NumElevators; i++ {
		if database.ElevatorsInNetwork[i].ElevatorID == elevatorID { //Sjekker at calgt heis ikke er unconnected
			return true
		}
	}
	return false
}

func UpdateDatabase(aliveMsg elevator.IAmAliveMessageStruct, database ElevatorDatabase) {

	if aliveMsg.Elevator.Operating != elevator.WS_NoMotor {
		aliveMsg.Elevator.Operating = elevator.WS_Running //OBS! Nå håndterer vi running-state som connected
	}

	for i := 0; i < database.NumElevators; i++ {
		if database.ElevatorsInNetwork[i].ElevatorID == aliveMsg.ElevatorID {
			database.ElevatorsInNetwork[i] = aliveMsg.Elevator
		}
	}
}

func WhatFloorIsElevatorFromStringID(database ElevatorDatabase, elevatorID string) int {

	for i := 0; i < database.NumElevators; i++ {
		if database.ElevatorsInNetwork[i].ElevatorID == elevatorID {
			return database.ElevatorsInNetwork[i].Floor
		}
	}
	return -1
}

func WhatStateIsElevatorFromStringID(database ElevatorDatabase, elevatorID string) elevator.ElevatorBehaviour {
	for i := 0; i < database.NumElevators; i++ {
		if database.ElevatorsInNetwork[i].ElevatorID == elevatorID {
			return database.ElevatorsInNetwork[i].Behaviour
		}
	}
	return elevator.EB_Undefined
}

func UpdateElevatorNetworkStateInDatabase(peerUpdate peers.PeerUpdate, database ElevatorDatabase) {
	for i := 0; i < len(database.ElevatorsInNetwork); i++ {
		if !peers.IsPeerOnNetwork(database.ElevatorsInNetwork[i], peerUpdate) {
			database.ElevatorsInNetwork[i].Operating = elevator.WS_Unconnected
		}

	}
}

func GetElevatorFromID(database ElevatorDatabase, elevatorID string) elevator.Elevator {
	var e elevator.Elevator
	for i := 0; i < database.NumElevators; i++ {
		if database.ElevatorsInNetwork[i].ElevatorID == elevatorID {
			return database.ElevatorsInNetwork[i]
		}
	}
	return e
}

func SendOrderMessage(orderTx chan elevator.OrderMessageStruct, button elevio.ButtonEvent, database ElevatorDatabase) {
	chosenElevator := AssignOrderToElevator(database, button)

	orderMsg := elevator.OrderMessageStruct{SystemID: "Gruppe10",
		MessageID:      "Order",
		ElevatorID:     elevator.MyID,
		OrderedButton:  button,
		ChosenElevator: chosenElevator}

	orderTx <- orderMsg
}
