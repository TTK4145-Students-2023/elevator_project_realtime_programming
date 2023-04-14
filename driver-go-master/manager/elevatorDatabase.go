package manager

import (
	"Driver-go/elevator"
	"Driver-go/elevio"
	"time"

	"fmt"
)

type ElevatorDatabase struct {
	ConnectedElevators int
	ElevatorsInNetwork []elevator.Elevator
}

func AssignOrderToElevator(database ElevatorDatabase, order elevio.ButtonEvent) string {

	lowCost := 100000.0
	lowestCostElevator := ""

	connectedElevators := database.ElevatorsInNetwork
	fmt.Println("The number of elevators that we have data on in the database is: ", len(connectedElevators))
	fmt.Println("And the number of connected elevators is: ", database.ConnectedElevators)

	if order.Button == elevio.BT_Cab || elevator.GetIAmAlone() {
		lowestCostElevator = elevator.MyID
	} else if elevator.AvailableAtCurrFloor(order.Floor) {
		lowestCostElevator = elevator.MyID
	} else {
		for i := 0; i < database.ConnectedElevators; i++ {
			c := calculateCost(connectedElevators[i], order)

			if c < lowCost && connectedElevators[i].Operating == elevator.WS_Connected {
				lowCost = c
				lowestCostElevator = connectedElevators[i].ElevatorID
			} else if c == lowCost && connectedElevators[i].Operating == elevator.WS_Connected {
				if lowestCostElevator > connectedElevators[i].ElevatorID {
					lowCost = c
					lowestCostElevator = connectedElevators[i].ElevatorID
				}
			}

		}

	}

	fmt.Println("Assigned order to:", lowestCostElevator)
	return lowestCostElevator
}

func FindDeadOrders(database ElevatorDatabase, deadElevatorID string) []elevio.ButtonEvent {
	deadElev := GetElevatorFromID(database, deadElevatorID)
	var deadOrders []elevio.ButtonEvent

	for floor := 0; floor < elevator.NumFloors; floor++ {
		for button := elevio.BT_HallUp; button < elevio.BT_Cab; button++ {
			var order elevio.ButtonEvent
			order.Button = elevio.ButtonType(button)
			order.Floor = floor

			if deadElev.Requests[floor][button].ElevatorID == deadElevatorID {
				deadOrders = append(deadOrders, order)
			}
		}

	}
	return deadOrders
}

func IsElevatorInDatabase(elevatorID string, database ElevatorDatabase) bool {
	for i := 0; i < database.ConnectedElevators; i++ {
		if database.ElevatorsInNetwork[i].ElevatorID == elevatorID { //Sjekker at calgt heis ikke er unconnected
			return true
		}
	}
	return false
}

func UpdateDatabase(elevatorToBeUpdated elevator.Elevator, database ElevatorDatabase) ElevatorDatabase {
	/*if elevatorToBeUpdated.Operating != elevator.WS_Immobile {
		elevatorToBeUpdated.Operating = elevator.WS_Connected
		fmt.Println("Her setter jeg operating staten til ", elevatorToBeUpdated.ElevatorID, " til ", elevatorToBeUpdated.Operating) //OBS! Nå håndterer vi running-state som connected
	}*/

	for i := 0; i < database.ConnectedElevators; i++ {
		if database.ElevatorsInNetwork[i].ElevatorID == elevatorToBeUpdated.ElevatorID {
			database.ElevatorsInNetwork[i] = elevatorToBeUpdated
		}
	}
	return database
}

func WhatFloorIsElevatorFromStringID(database ElevatorDatabase, elevatorID string) int {
	for i := 0; i < database.ConnectedElevators; i++ {
		if database.ElevatorsInNetwork[i].ElevatorID == elevatorID {
			return database.ElevatorsInNetwork[i].Floor
		}
	}
	return -1
}

func WhatStateIsElevatorFromStringID(database ElevatorDatabase, elevatorID string) elevator.ElevatorBehaviour {
	for i := 0; i < database.ConnectedElevators; i++ {
		if database.ElevatorsInNetwork[i].ElevatorID == elevatorID {
			return database.ElevatorsInNetwork[i].Behaviour
		}
	}
	return elevator.EB_Undefined
}

func UpdateElevatorNetworkStateInDatabase(elevatorID string, database ElevatorDatabase, newState elevator.WorkingState) ElevatorDatabase {
	for i := 0; i < len(database.ElevatorsInNetwork); i++ {
		if elevatorID == database.ElevatorsInNetwork[i].ElevatorID {
			database.ElevatorsInNetwork[i].Operating = newState
			fmt.Println("Her setter jeg operating staten til ", database.ElevatorsInNetwork[i].ElevatorID, " til ", database.ElevatorsInNetwork[i].Operating)
			if newState == elevator.WS_Unconnected {
				database.ConnectedElevators--
			} else if newState == elevator.WS_Connected {
				database.ConnectedElevators++
			}
		}

	}

	return database
}

func GetElevatorFromID(database ElevatorDatabase, elevatorID string) elevator.Elevator {
	var e elevator.Elevator
	for i := 0; i < database.ConnectedElevators; i++ {
		if database.ElevatorsInNetwork[i].ElevatorID == elevatorID {
			return database.ElevatorsInNetwork[i]
		}
	}
	return e
}

func SendCabCallsForElevator(database ElevatorDatabase, newPeer string) []elevator.OrderMessageStruct {
	var cabsToBeSent []elevator.OrderMessageStruct
	for i := 0; i < len(database.ElevatorsInNetwork); i++ {
		if database.ElevatorsInNetwork[i].ElevatorID == newPeer && newPeer != elevator.MyID {
			for floor := 0; floor < elevator.NumFloors; floor++ {
				if database.ElevatorsInNetwork[i].Requests[floor][elevio.BT_Cab].ElevatorID == newPeer {
					var button elevio.ButtonEvent
					button.Floor = floor
					button.Button = elevio.BT_Cab
					panelPair := elevator.OrderpanelPair{ElevatorID: newPeer, OrderState: elevator.SO_Confirmed}
					cabsToBeSent = append(cabsToBeSent, elevator.MakeOrderMessage(panelPair, button))
				}
			}
		}
	}
	return cabsToBeSent
}

func SearchMessageOrderUpdate(aliveMessage elevator.IAmAliveMessageStruct, database ElevatorDatabase) []elevator.OrderMessageStruct {

	var newChangedOrders []elevator.OrderMessageStruct

	localElevator := GetElevatorFromID(database, elevator.MyID)

	for floor := 0; floor < elevator.NumFloors; floor++ {
		for button := elevio.BT_HallUp; button < elevio.BT_Cab; button++ {

			currentButtonEvent := elevio.ButtonEvent{Floor: floor, Button: button}

			receivedOrderState := aliveMessage.Elevator.Requests[floor][button].OrderState
			localOrderState := localElevator.Requests[floor][button].OrderState

			receivedRequestID := aliveMessage.Elevator.Requests[floor][button].ElevatorID
			localRequestID := localElevator.Requests[floor][button].ElevatorID

			if receivedOrderState != localOrderState ||
				receivedRequestID != localRequestID {

				if receivedOrderState == elevator.SO_NoOrder {

					if localRequestID == aliveMessage.ElevatorID &&
						localOrderState == elevator.SO_Confirmed {

						panelPair := elevator.OrderpanelPair{ElevatorID: aliveMessage.ElevatorID, OrderState: elevator.SO_NoOrder}
						newChangedOrders = append(newChangedOrders, elevator.MakeOrderMessage(panelPair, currentButtonEvent))
						fmt.Println("I found a no order so I will erase it")

					} else if localRequestID == aliveMessage.ElevatorID &&
						localOrderState == elevator.SO_NewOrder {

						panelPair := elevator.OrderpanelPair{ElevatorID: aliveMessage.ElevatorID, OrderState: elevator.SO_NoOrder}
						newChangedOrders = append(newChangedOrders, elevator.MakeOrderMessage(panelPair, currentButtonEvent))

					}
				} else if receivedOrderState == elevator.SO_NewOrder {

					if receivedRequestID == localElevator.ElevatorID {
						fmt.Println("Her fikk jeg en new order og den var til meg så jeg setter den confirmed")
						panelPair := elevator.OrderpanelPair{ElevatorID: localElevator.ElevatorID, OrderState: elevator.SO_Confirmed}
						newChangedOrders = append(newChangedOrders, elevator.MakeOrderMessage(panelPair, currentButtonEvent))

					} else {
						panelPair := elevator.OrderpanelPair{ElevatorID: aliveMessage.ElevatorID, OrderState: elevator.SO_NewOrder}
						newChangedOrders = append(newChangedOrders, elevator.MakeOrderMessage(panelPair, currentButtonEvent))
					}
				} else if receivedOrderState == elevator.SO_Confirmed {

					if receivedRequestID == aliveMessage.ElevatorID {
						fmt.Println("Her fikk jeg en confirmed order og den var fra eieren så jeg setter den confirmed")
						panelPair := elevator.OrderpanelPair{ElevatorID: aliveMessage.ElevatorID, OrderState: elevator.SO_Confirmed}
						newChangedOrders = append(newChangedOrders, elevator.MakeOrderMessage(panelPair, currentButtonEvent))

					}
				}
				time.Sleep(25 * time.Millisecond)
			}
		}
	}
	return newChangedOrders
}
