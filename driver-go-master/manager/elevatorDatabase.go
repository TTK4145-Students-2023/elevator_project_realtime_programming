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

			if receivedOrderState != localOrderState { // || receivedRequestID != localRequestID {

				/*changedOwner := receivedRequestID != localRequestID
				fmt.Println("Endring i eier av ordre:", (receivedRequestID != localRequestID))
				if changedOwner {
					fmt.Println("Old owner:", localRequestID)
					fmt.Println("New owner:", receivedRequestID)
				}*/
				if receivedOrderState == singleElevator.SO_NoOrder {

					if localRequestID == stateUpdateMessage.ElevatorID &&
						localOrderState == singleElevator.SO_Confirmed {

						panelPair := singleElevator.OrderpanelPair{ElevatorID: stateUpdateMessage.ElevatorID, OrderState: singleElevator.SO_NoOrder}
						newChangedOrders = append(newChangedOrders, MakeOrder(panelPair, currentButtonEvent))
						fmt.Println("I found a no order so I will erase it")

					} else if localRequestID == stateUpdateMessage.ElevatorID &&
						localOrderState == singleElevator.SO_NewOrder {

						panelPair := singleElevator.OrderpanelPair{ElevatorID: stateUpdateMessage.ElevatorID, OrderState: singleElevator.SO_NoOrder}
						newChangedOrders = append(newChangedOrders, MakeOrder(panelPair, currentButtonEvent))

					}
				} else if receivedOrderState == singleElevator.SO_NewOrder {

					if receivedRequestID == localElevator.ElevatorID {
						fmt.Println("Her fikk jeg en new order og den var til meg så jeg setter den confirmed")
						panelPair := singleElevator.OrderpanelPair{ElevatorID: localElevator.ElevatorID, OrderState: singleElevator.SO_Confirmed}
						newChangedOrders = append(newChangedOrders, MakeOrder(panelPair, currentButtonEvent))

					} else {
						panelPair := singleElevator.OrderpanelPair{ElevatorID: stateUpdateMessage.ElevatorID, OrderState: singleElevator.SO_NewOrder}
						newChangedOrders = append(newChangedOrders, MakeOrder(panelPair, currentButtonEvent))
					}
				} else if receivedOrderState == singleElevator.SO_Confirmed {

					if receivedRequestID == stateUpdateMessage.ElevatorID {
						fmt.Println("Her fikk jeg en confirmed order og den var fra eieren så jeg setter den confirmed")
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
