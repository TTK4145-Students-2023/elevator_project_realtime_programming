package manager

import (
	"Driver-go/elevator"
	"Driver-go/elevio"
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

func WhatFloorIsElevator(database ElevatorDatabase, elevatorID string) int{

	for i := 0; i < database.NumElevators; i++{
		if database.ElevatorsInNetwork[i].ElevatorID == elevatorID{
			return database.ElevatorsInNetwork[i].Floor
		}
	}
	return -1
}