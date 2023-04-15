package manager

import (
	"Driver-go/elevio"
	"Driver-go/singleElevator"
	"time"

	"fmt"
)

type ElevatorDatabase struct {
	ConnectedElevators int
	ElevatorList       []singleElevator.Elevator
}

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
			fmt.Println("Her setter jeg operating staten til ", database.ElevatorList[i].ElevatorID, " til ", database.ElevatorList[i].Operating)
			if newState == singleElevator.WS_Unconnected {
				database.ConnectedElevators--
			} else if newState == singleElevator.WS_Connected {
				database.ConnectedElevators++
			}
		}

	}

	return database
}

func UpdateDatabaseWithDeadOrders(deadElevatorID string, immobilityTimer *time.Timer, doorTimer *time.Timer, database ElevatorDatabase) ElevatorDatabase {
	var deadOrders []elevio.ButtonEvent
	deadOrders = FindDeadOrders(database, deadElevatorID)
	for j := 0; j < len(deadOrders); j++ {
		chosenElevator := AssignOrderToElevator(database, deadOrders[j])
		newElevatorUpdate := singleElevator.HandleNewOrder(chosenElevator, deadOrders[j], doorTimer, immobilityTimer)
		database = UpdateDatabase(newElevatorUpdate, database)
	}
	return database
}

func UpdateDatabaseFromIncomingMessages(stateUpdateMessage singleElevator.ElevatorUpdateToDatabase, database ElevatorDatabase, immobilityTimer *time.Timer, doorTimer *time.Timer) ElevatorDatabase {
	database = UpdateDatabase(stateUpdateMessage.Elevator, database)

	newChangedOrders := SearchMessageForOrderUpdate(stateUpdateMessage, database)
	for i := 0; i < len(newChangedOrders); i++ {
		newOrder := newChangedOrders[i]
		var newElevatorUpdate singleElevator.Elevator

		if newOrder.PanelPair.OrderState == singleElevator.SO_Confirmed {
			chosenElevator := newOrder.PanelPair.ElevatorID
			newButton := newOrder.OrderedButton

			newElevatorUpdate = singleElevator.HandleConfirmedOrder(chosenElevator, newButton, doorTimer, immobilityTimer)

		} else if newOrder.PanelPair.OrderState == singleElevator.SO_NoOrder {
			newElevatorUpdate = singleElevator.Requests_clearOnFloor(newOrder.PanelPair.ElevatorID, newOrder.OrderedButton.Floor)
		}

		database = UpdateDatabase(newElevatorUpdate, database)
	}
	return database

}

func HandleRestoredCabs(newCabs OrderStruct, doorTimer *time.Timer, immobilityTimer *time.Timer) singleElevator.Elevator {
	var newElevatorUpdate singleElevator.Elevator
	if MessageIDEqualsMyID(newCabs.PanelPair.ElevatorID) {
		newElevatorUpdate = singleElevator.Fsm_onRequestButtonPress(newCabs.OrderedButton.Floor, newCabs.OrderedButton.Button, singleElevator.MyID, doorTimer, immobilityTimer)
	}
	return newElevatorUpdate
}

func SearchMessageForOrderUpdate(stateUpdateMessage singleElevator.ElevatorUpdateToDatabase, database ElevatorDatabase) []OrderStruct {

	var newChangedOrders []OrderStruct

	localElevator := GetElevatorFromID(database, singleElevator.MyID)

	for floor := 0; floor < singleElevator.NumFloors; floor++ {
		for button := elevio.BT_HallUp; button < elevio.BT_Cab; button++ {

			currentButtonEvent := elevio.ButtonEvent{Floor: floor, Button: button}

			receivedOrderState := stateUpdateMessage.Elevator.Requests[floor][button].OrderState
			localOrderState := localElevator.Requests[floor][button].OrderState

			receivedRequestID := stateUpdateMessage.Elevator.Requests[floor][button].ElevatorID
			localRequestID := localElevator.Requests[floor][button].ElevatorID

			if receivedOrderState != localOrderState {
				if receivedOrderState == singleElevator.SO_NoOrder {

					if localRequestID == stateUpdateMessage.ElevatorID &&
						localOrderState == singleElevator.SO_Confirmed {

						panelPair := singleElevator.OrderpanelPair{ElevatorID: stateUpdateMessage.ElevatorID, OrderState: singleElevator.SO_NoOrder}
						newChangedOrders = append(newChangedOrders, MakeOrder(panelPair, currentButtonEvent))

					} else if localRequestID == stateUpdateMessage.ElevatorID &&
						localOrderState == singleElevator.SO_NewOrder {

						panelPair := singleElevator.OrderpanelPair{ElevatorID: stateUpdateMessage.ElevatorID, OrderState: singleElevator.SO_NoOrder}
						newChangedOrders = append(newChangedOrders, MakeOrder(panelPair, currentButtonEvent))

					}
				} else if receivedOrderState == singleElevator.SO_NewOrder {

					if receivedRequestID == localElevator.ElevatorID {
						panelPair := singleElevator.OrderpanelPair{ElevatorID: localElevator.ElevatorID, OrderState: singleElevator.SO_Confirmed}
						newChangedOrders = append(newChangedOrders, MakeOrder(panelPair, currentButtonEvent))

					} else {
						panelPair := singleElevator.OrderpanelPair{ElevatorID: stateUpdateMessage.ElevatorID, OrderState: singleElevator.SO_NewOrder}
						newChangedOrders = append(newChangedOrders, MakeOrder(panelPair, currentButtonEvent))
					}
				} else if receivedOrderState == singleElevator.SO_Confirmed {

					if receivedRequestID == stateUpdateMessage.ElevatorID {
						panelPair := singleElevator.OrderpanelPair{ElevatorID: stateUpdateMessage.ElevatorID, OrderState: singleElevator.SO_Confirmed}
						newChangedOrders = append(newChangedOrders, MakeOrder(panelPair, currentButtonEvent))

					}
				}
				time.Sleep(25 * time.Millisecond)
			}
		}
	}
	return newChangedOrders
}
