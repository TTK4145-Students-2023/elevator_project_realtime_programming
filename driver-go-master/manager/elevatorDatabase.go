package manager

import (
	"Driver-go/elevator"
	"Driver-go/elevio"
	"time"

	"fmt"
)

type ElevatorDatabase struct {
	ConnectedElevators int
	ElevatorList       []elevator.Elevator
}

func UpdateDatabase(elevatorToBeUpdated elevator.Elevator, database ElevatorDatabase) ElevatorDatabase {
	/*if elevatorToBeUpdated.Operating != elevator.WS_Immobile {
		elevatorToBeUpdated.Operating = elevator.WS_Connected
		fmt.Println("Her setter jeg operating staten til ", elevatorToBeUpdated.ElevatorID, " til ", elevatorToBeUpdated.Operating) //OBS! Nå håndterer vi running-state som connected
	}*/

	for i := 0; i < len(database.ElevatorList); i++ {
		if database.ElevatorList[i].ElevatorID == elevatorToBeUpdated.ElevatorID {
			database.ElevatorList[i] = elevatorToBeUpdated
		}
	}
	return database
}

func UpdateElevatorNetworkStateInDatabase(elevatorID string, database ElevatorDatabase, newState elevator.WorkingState) ElevatorDatabase {
	for i := 0; i < len(database.ElevatorList); i++ {
		if elevatorID == database.ElevatorList[i].ElevatorID {
			database.ElevatorList[i].Operating = newState
			fmt.Println("Her setter jeg operating staten til ", database.ElevatorList[i].ElevatorID, " til ", database.ElevatorList[i].Operating)
			if newState == elevator.WS_Unconnected {
				database.ConnectedElevators--
			} else if newState == elevator.WS_Connected {
				database.ConnectedElevators++
			}
		}

	}

	return database
}

func SearchMessageForOrderUpdate(aliveMessage elevator.StateUpdateStruct, database ElevatorDatabase) []elevator.OrderStruct {

	var newChangedOrders []elevator.OrderStruct

	localElevator := GetElevatorFromID(database, elevator.MyID)

	for floor := 0; floor < elevator.NumFloors; floor++ {
		for button := elevio.BT_HallUp; button < elevio.BT_Cab; button++ {

			currentButtonEvent := elevio.ButtonEvent{Floor: floor, Button: button}

			receivedOrderState := aliveMessage.Elevator.Requests[floor][button].OrderState
			localOrderState := localElevator.Requests[floor][button].OrderState

			receivedRequestID := aliveMessage.Elevator.Requests[floor][button].ElevatorID
			localRequestID := localElevator.Requests[floor][button].ElevatorID

			if receivedOrderState != localOrderState { // || receivedRequestID != localRequestID {

				/*changedOwner := receivedRequestID != localRequestID
				fmt.Println("Endring i eier av ordre:", (receivedRequestID != localRequestID))
				if changedOwner {
					fmt.Println("Old owner:", localRequestID)
					fmt.Println("New owner:", receivedRequestID)
				}*/
				if receivedOrderState == elevator.SO_NoOrder {

					if localRequestID == aliveMessage.ElevatorID &&
						localOrderState == elevator.SO_Confirmed {

						panelPair := elevator.OrderpanelPair{ElevatorID: aliveMessage.ElevatorID, OrderState: elevator.SO_NoOrder}
						newChangedOrders = append(newChangedOrders, elevator.MakeOrder(panelPair, currentButtonEvent))
						fmt.Println("I found a no order so I will erase it")

					} else if localRequestID == aliveMessage.ElevatorID &&
						localOrderState == elevator.SO_NewOrder {

						panelPair := elevator.OrderpanelPair{ElevatorID: aliveMessage.ElevatorID, OrderState: elevator.SO_NoOrder}
						newChangedOrders = append(newChangedOrders, elevator.MakeOrder(panelPair, currentButtonEvent))

					}
				} else if receivedOrderState == elevator.SO_NewOrder {

					if receivedRequestID == localElevator.ElevatorID {
						fmt.Println("Her fikk jeg en new order og den var til meg så jeg setter den confirmed")
						panelPair := elevator.OrderpanelPair{ElevatorID: localElevator.ElevatorID, OrderState: elevator.SO_Confirmed}
						newChangedOrders = append(newChangedOrders, elevator.MakeOrder(panelPair, currentButtonEvent))

					} else {
						panelPair := elevator.OrderpanelPair{ElevatorID: aliveMessage.ElevatorID, OrderState: elevator.SO_NewOrder}
						newChangedOrders = append(newChangedOrders, elevator.MakeOrder(panelPair, currentButtonEvent))
					}
				} else if receivedOrderState == elevator.SO_Confirmed {

					if receivedRequestID == aliveMessage.ElevatorID {
						fmt.Println("Her fikk jeg en confirmed order og den var fra eieren så jeg setter den confirmed")
						panelPair := elevator.OrderpanelPair{ElevatorID: aliveMessage.ElevatorID, OrderState: elevator.SO_Confirmed}
						newChangedOrders = append(newChangedOrders, elevator.MakeOrder(panelPair, currentButtonEvent))

					}
				}
				time.Sleep(25 * time.Millisecond)
			}
		}
	}
	return newChangedOrders
}
