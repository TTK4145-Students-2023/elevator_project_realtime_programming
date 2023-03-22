package manager

import (
	"Driver-go/elevator"
	"Driver-go/elevio"

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
	fmt.Println("The connected elevators are: ", len(connectedElevators))
	fmt.Println("And the number of connected elevators is: ", database.NumElevators)

	if order.Button == elevio.BT_Cab || elevator.GetIAmAlone() {
		elevatorID = elevator.MyID
	} else {
		for i := 0; i < database.NumElevators; i++ {
			c := calculateCost(&connectedElevators[i], order)                            //OBS! Blanding av pekere og ikke pekere
			if c < lowCost && connectedElevators[i].Operating == elevator.WS_Connected { //Sjekker at calgt heis ikke er unconnected
				lowCost = c
				elevatorID = connectedElevators[i].ElevatorID
			}
		}
	}

	fmt.Println("Assigned order to: ", elevatorID)
	return elevatorID
}

func ReassignDeadOrders(orderTx chan elevator.OrderMessageStruct, database ElevatorDatabase, deadElevatorID string) {
	deadElev := GetElevatorFromID(database, deadElevatorID)
	fmt.Println(" -----dead elevator id -----")
	fmt.Println(deadElev.ElevatorID)
	fmt.Println(("here are the orders"))
	for floor := 0; floor < elevator.NumFloors; floor++ {
		for button := elevio.BT_HallUp; button < elevio.BT_Cab; button++ {
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

func UpdateDatabase(elevatorToBeUpdated elevator.Elevator, database ElevatorDatabase) {

	if elevatorToBeUpdated.Operating != elevator.WS_Immobile {
		elevatorToBeUpdated.Operating = elevator.WS_Connected //OBS! Nå håndterer vi running-state som connected
	}

	for i := 0; i < database.NumElevators; i++ {
		if database.ElevatorsInNetwork[i].ElevatorID == elevatorToBeUpdated.ElevatorID {
			database.ElevatorsInNetwork[i] = elevatorToBeUpdated
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

func UpdateElevatorNetworkStateInDatabase(elevatorID string, database ElevatorDatabase, newState elevator.WorkingState) {
	for i := 0; i < len(database.ElevatorsInNetwork); i++ {
		if elevatorID == database.ElevatorsInNetwork[i].ElevatorID {
			database.ElevatorsInNetwork[i].Operating = newState
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

func SendCabCallsForElevator(orderTx chan elevator.OrderMessageStruct, database ElevatorDatabase, newPeer string) {
	for i := 0; i < len(database.ElevatorsInNetwork); i++ {
		if database.ElevatorsInNetwork[i].ElevatorID == newPeer {
			fmt.Println("her har den matchende id som den som kom tilbake")
			for floor := 0; floor < elevator.NumFloors; floor++ {
				if database.ElevatorsInNetwork[i].Requests[floor][elevio.BT_Cab].ElevatorID == newPeer {
					var button elevio.ButtonEvent
					button.Floor = floor
					button.Button = elevio.BT_Cab
					fmt.Println("Cab call to be sent: ", button)
					orderTx <- elevator.MakeOrderMessage(newPeer, button)
				}
			}
		}
	}
}

//ny meldinger oppdtaeres i databasen, og heisen henter inn fra databasen hvor den skal kjøre

func SearchMessageOrderUpdate(aliveMessage elevator.IAmAliveMessageStruct, database ElevatorDatabase) elevator.Elevator {
	localElevator := GetElevatorFromID(database, elevator.MyID)

	for floor := 0; floor < elevator.NumFloors; floor++ {
		for button := elevio.BT_HallUp; button < elevio.BT_Cab; button++ {

			if aliveMessage.Elevator.Requests[floor][button].OrderState != localElevator.Requests[floor][button].OrderState ||
				aliveMessage.Elevator.Requests[floor][button].ElevatorID != localElevator.Requests[floor][button].ElevatorID {
				//Kan oppstå forskjellige IDer ved reassignment
				switch aliveMessage.Elevator.Requests[floor][button].OrderState {
				case elevator.SO_NoOrder:
					if localElevator.Requests[floor][button].ElevatorID == aliveMessage.ElevatorID {
						localElevator.Requests[floor][button].OrderState = elevator.SO_NoOrder
						localElevator.Requests[floor][button].ElevatorID = ""
					}
				case elevator.SO_NewOrder:
					if aliveMessage.Elevator.Requests[floor][button].ElevatorID == localElevator.ElevatorID {
						localElevator.Requests[floor][button].OrderState = elevator.SO_Confirmed
						localElevator.Requests[floor][button].ElevatorID = localElevator.ElevatorID
					} else if aliveMessage.Elevator.Requests[floor][button].ElevatorID != localElevator.ElevatorID {
						localElevator.Requests[floor][button] = aliveMessage.Elevator.Requests[floor][button]
					}
				case elevator.SO_Confirmed:
					if localElevator.Requests[floor][button].ElevatorID == aliveMessage.ElevatorID {
						localElevator.Requests[floor][button].OrderState = elevator.SO_Confirmed
					}
				}
			}
		}
	}

	return localElevator
}
