package databaseHandler

import (
	"Driver-go/elevatorHardware"
	"Driver-go/singleElevator"
	"time"
)

func UpdateDatabase(elevatorToBeUpdated singleElevator.Elevator, database ElevatorDatabase) ElevatorDatabase {

	for i := 0; i < len(database.ElevatorList); i++ {
		if database.ElevatorList[i].ElevatorID == elevatorToBeUpdated.ElevatorID {
			database.ElevatorList[i] = elevatorToBeUpdated
		}
	}
	return database
}

func UpdateElevatorNetworkStateInDatabase(elevatorID string, database ElevatorDatabase, newState singleElevator.WorkingState) ElevatorDatabase {
	for i := 0; i < len(database.ElevatorList); i++ {
		if elevatorID == database.ElevatorList[i].ElevatorID {
			database.ElevatorList[i].Operating = newState
			if newState == singleElevator.Unconnected {
				database.ConnectedElevators--
			} else if newState == singleElevator.Connected {
				database.ConnectedElevators++
			}
		}

	}

	return database
}

func UpdateDatabaseWithDeadOrders(deadElevatorID string, immobilityTimer *time.Timer, doorTimer *time.Timer, database ElevatorDatabase) ElevatorDatabase {
	deadOrders := FindDeadOrders(database, deadElevatorID)
	for j := 0; j < len(deadOrders); j++ {
		chosenElevator := AssignOrderToElevator(database, deadOrders[j])
		newElevatorUpdate := singleElevator.HandleNewOrder(chosenElevator, deadOrders[j], doorTimer, immobilityTimer)
		database = UpdateDatabase(newElevatorUpdate, database)
	}
	return database
}

func UpdateDatabaseFromIncomingMessages(stateUpdateMessage singleElevator.ElevatorStateUpdate, database ElevatorDatabase, immobilityTimer *time.Timer, doorTimer *time.Timer) ElevatorDatabase {
	database = UpdateDatabase(stateUpdateMessage.Elevator, database)

	OrderDifferencesFound := FindChangesBetweenIncomingmessageAndLocalDatabase(stateUpdateMessage, database)

	for i := 0; i < len(OrderDifferencesFound); i++ {
		orderDifference := OrderDifferencesFound[i]
		var newDatabaseEntry singleElevator.Elevator
		stateOfOrderDifference := orderDifference.PanelPair.OrderState

		if stateOfOrderDifference == singleElevator.ConfirmedOrder { //Dårlig navn
			chosenElevator := orderDifference.PanelPair.ElevatorID
			newButton := orderDifference.OrderedButton

			newDatabaseEntry = singleElevator.HandleConfirmedOrder(chosenElevator, newButton, doorTimer, immobilityTimer)

		} else if orderDifference.PanelPair.OrderState == singleElevator.NoOrder {
			newDatabaseEntry = singleElevator.ClearOrderAtThisFloor(orderDifference.PanelPair.ElevatorID, orderDifference.OrderedButton.Floor)
		}

		database = UpdateDatabase(newDatabaseEntry, database)
	}
	return database

}

func HandleNewFloorAndUpdateDatabase(floor int, database ElevatorDatabase, doorTimer *time.Timer, immobilityTimer *time.Timer) ElevatorDatabase {
	newDatabaseEntry := singleElevator.FloorArrival(floor, doorTimer, immobilityTimer)
	database = UpdateDatabase(newDatabaseEntry, database)
	return database
}

func HandleNewButtonAndUpdateDatabase(button elevatorHardware.ButtonEvent, database ElevatorDatabase, doorTimer *time.Timer, immobilityTimer *time.Timer) ElevatorDatabase {
	chosenElevator := AssignOrderToElevator(database, button)
	newDatabaseEntry := singleElevator.HandleNewOrder(chosenElevator, button, doorTimer, immobilityTimer)
	database = UpdateDatabase(newDatabaseEntry, database)
	return database
}
