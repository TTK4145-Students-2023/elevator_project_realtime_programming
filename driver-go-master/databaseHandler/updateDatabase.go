package databaseHandler

import (
	"Driver-go/elevatorHardware"
	"Driver-go/singleElevator"
	"time"
)

//Wrapper functions for updating database. Designed for increasing readability in main. 

func HandleNewFloorAndUpdateDatabase(floor int, database ElevatorDatabase, doorTimer *time.Timer, immobilityTimer *time.Timer) ElevatorDatabase {
	newDatabaseEntry := singleElevator.FloorArrival(floor, doorTimer, immobilityTimer)
	database = UpdateDatabase(newDatabaseEntry, database)
	return database
}

func HandleNewButtonAndUpdateDatabase(button elevatorHardware.ButtonEvent, database ElevatorDatabase, doorTimer *time.Timer, immobilityTimer *time.Timer) ElevatorDatabase {
	chosenElevator := AssignOrderToElevator(database, button)

	var newDatabaseEntry singleElevator.Elevator
		if chosenElevator == singleElevator.MyID {
			newDatabaseEntry = singleElevator.ExecuteAssignedOrder(button, chosenElevator, doorTimer, immobilityTimer)
		} else {
			newDatabaseEntry = singleElevator.SaveLocalNewOrder(button, chosenElevator)
		}

	database = UpdateDatabase(newDatabaseEntry, database)
	return database
}


//Functions to update the database based on incoming messages and network connection.

func UpdateDatabase(newDatabaseEntry singleElevator.Elevator, database ElevatorDatabase) ElevatorDatabase {
	for i := 0; i < len(database.ElevatorList); i++ {
		if database.ElevatorList[i].ElevatorID == newDatabaseEntry.ElevatorID {
			database.ElevatorList[i] = newDatabaseEntry
		}
	}
	return database
}

func UpdateElevatorNetworkStateInDatabase(newState singleElevator.WorkingState, elevatorID string, database ElevatorDatabase) ElevatorDatabase {
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

		var newDatabaseEntry singleElevator.Elevator
		if chosenElevator == singleElevator.MyID {
			newDatabaseEntry = singleElevator.ExecuteAssignedOrder(deadOrders[j], chosenElevator, doorTimer, immobilityTimer)
		} else {
			newDatabaseEntry = singleElevator.SaveLocalNewOrder(deadOrders[j], chosenElevator)
		}
		database = UpdateDatabase(newDatabaseEntry, database)



	}
	return database
}

func UpdateDatabaseFromIncomingMessages(stateUpdateMessage singleElevator.ElevatorStateUpdate, database ElevatorDatabase, immobilityTimer *time.Timer, doorTimer *time.Timer) ElevatorDatabase {
	database = UpdateDatabase(stateUpdateMessage.Elevator, database)

	OrderDifferencesFound := FindChangesBetweenIncomingMessageAndLocalDatabase(stateUpdateMessage, database)

	for i := 0; i < len(OrderDifferencesFound); i++ {
		var newDatabaseEntry singleElevator.Elevator

		orderDifference := OrderDifferencesFound[i]
		orderState := GetStateOfOrder(orderDifference)
		
		if orderState == singleElevator.ConfirmedOrder { 
			chosenElevator := orderDifference.PanelPair.AssingedElevatorID
			orderedButton := orderDifference.OrderedButton

			if chosenElevator == singleElevator.MyID {
				newDatabaseEntry = singleElevator.ExecuteAssignedOrder(orderedButton, chosenElevator, doorTimer, immobilityTimer)
			} else {
				newDatabaseEntry = singleElevator.SaveLocalConfirmedOrder(orderedButton, chosenElevator) //endre navn mer deskrriptivt
			}
		
		} else if orderState == singleElevator.NoOrder {
			newDatabaseEntry = singleElevator.ClearOrderAtThisFloor(orderDifference.PanelPair.AssingedElevatorID, orderDifference.OrderedButton.Floor)
		}

		database = UpdateDatabase(newDatabaseEntry, database)
	}
	return database

}


