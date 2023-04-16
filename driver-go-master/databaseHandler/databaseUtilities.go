package databaseHandler

import (
	"Driver-go/elevatorHardware"
	"Driver-go/singleElevator"
)

//Contains the database that holds information about all the elevators states and orders.
//Also contains helperfunctions.

type ElevatorDatabase struct {
	ConnectedElevators int
	ElevatorList       []singleElevator.Elevator
}



func IsElevatorInDatabase(elevatorID string, database ElevatorDatabase) bool {
	for i := 0; i < len(database.ElevatorList); i++ {
		if database.ElevatorList[i].ElevatorID == elevatorID { //Sjekker at valgt heis ikke er unconnected
			return true
		}
	}
	return false
}

func MessageIDEqualsMyID(messageUpdateID string) bool {
	if messageUpdateID == singleElevator.MyID {
		return true
	} else {
		return false
	}
}

func GetElevatorFromID(database ElevatorDatabase, elevatorID string) singleElevator.Elevator {
	var e singleElevator.Elevator
	for i := 0; i < len(database.ElevatorList); i++ {
		if database.ElevatorList[i].ElevatorID == elevatorID {
			return database.ElevatorList[i]
		}
	}
	return e
}

func GetStateOfOrder(order OrderStruct) singleElevator.StateOfOrder{
	return order.PanelPair.OrderState
}

func FindDeadOrders(database ElevatorDatabase, deadElevatorID string) []elevatorHardware.ButtonEvent {
	deadElevator := GetElevatorFromID(database, deadElevatorID)
	var deadOrders []elevatorHardware.ButtonEvent
	var order elevatorHardware.ButtonEvent

	for floor := 0; floor < singleElevator.NumFloors; floor++ {
		for button := elevatorHardware.BT_HallUp; button < elevatorHardware.BT_Cab; button++ {
			ownerOfOrder := deadElevator.Requests[floor][button].AssingedElevatorID
			order.Button = elevatorHardware.ButtonType(button)
			order.Floor = floor

			if ownerOfOrder == deadElevatorID {
				deadOrders = append(deadOrders, order)
			}
		}

	}
	return deadOrders
}
