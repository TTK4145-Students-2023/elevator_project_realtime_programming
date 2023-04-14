package manager

import (
	"Driver-go/elevator"
	"Driver-go/elevio"
)

func IsElevatorInDatabase(elevatorID string, database ElevatorDatabase) bool {
	for i := 0; i < len(database.ElevatorList); i++ {
		if database.ElevatorList[i].ElevatorID == elevatorID { //Sjekker at valgt heis ikke er unconnected
			return true
		}
	}
	return false
}

func shouldITakeTheOrder(order elevio.ButtonEvent) bool {
	if order.Button == elevio.BT_Cab || elevator.GetIAmAlone() || elevator.AvailableAtCurrFloor(order.Floor) {
		return true
	} else {
		return false
	}
}

func GetElevatorFromID(database ElevatorDatabase, elevatorID string) elevator.Elevator {
	var e elevator.Elevator
	for i := 0; i < len(database.ElevatorList); i++ {
		if database.ElevatorList[i].ElevatorID == elevatorID {
			return database.ElevatorList[i]
		}
	}
	return e
}
