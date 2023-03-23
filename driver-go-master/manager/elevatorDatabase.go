package manager

import (
	"Driver-go/elevator"
	"Driver-go/elevio"
	"time"

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
	} else if elevator.AvailableAtCurrFloor(order.Floor) {
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
				//SendOrderMessage(orderTx, order, database)
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

func UpdateDatabase(elevatorToBeUpdated elevator.Elevator, database ElevatorDatabase) ElevatorDatabase {
	if elevatorToBeUpdated.Operating != elevator.WS_Immobile {
		elevatorToBeUpdated.Operating = elevator.WS_Connected //OBS! Nå håndterer vi running-state som connected
	}

	for i := 0; i < database.NumElevators; i++ {
		if database.ElevatorsInNetwork[i].ElevatorID == elevatorToBeUpdated.ElevatorID {
			database.ElevatorsInNetwork[i] = elevatorToBeUpdated
		}
	}
	return database
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

func UpdateElevatorNetworkStateInDatabase(elevatorID string, database ElevatorDatabase, newState elevator.WorkingState) ElevatorDatabase {
	for i := 0; i < len(database.ElevatorsInNetwork); i++ {
		if elevatorID == database.ElevatorsInNetwork[i].ElevatorID {
			database.ElevatorsInNetwork[i].Operating = newState
		}

	}
	return database
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

/*func SendOrderMessage(orderTx chan elevator.OrderMessageStruct, button elevio.ButtonEvent, database ElevatorDatabase) {
	chosenElevator := AssignOrderToElevator(database, button)

	orderMsg := elevator.OrderMessageStruct{SystemID: "Gruppe10",
		MessageID:      "Order",
		ElevatorID:     elevator.MyID,
		OrderedButton:  button,
		ChosenElevator: chosenElevator}

	orderTx <- orderMsg
}*/

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
					//orderTx <- elevator.MakeOrderMessage(newPeer, button)
				}
			}
		}
	}
}

//ny meldinger oppdtaeres i databasen, og heisen henter inn fra databasen hvor den skal kjøre

func SearchMessageOrderUpdate(aliveMessage elevator.IAmAliveMessageStruct, database ElevatorDatabase) []elevator.OrderMessageStruct {

	var newChangedOrders []elevator.OrderMessageStruct

	localElevator := GetElevatorFromID(database, elevator.MyID)
	//fmt.Println("This is the staus of the elevator in the database:")
	//elevator.ElevatorPrint(localElevator)

	for floor := 0; floor < elevator.NumFloors; floor++ {
		for button := elevio.BT_HallUp; button < elevio.BT_Cab; button++ {

			if aliveMessage.Elevator.Requests[floor][button].OrderState != localElevator.Requests[floor][button].OrderState ||
				aliveMessage.Elevator.Requests[floor][button].ElevatorID != localElevator.Requests[floor][button].ElevatorID {

				fmt.Println("Jeg har funnet en forskjell!\nMottatt array: ", aliveMessage.Elevator.Requests)
				fmt.Println("Mitt lokale array: ", localElevator.Requests)
				//OBS! Kan oppstå forskjellige IDer ved reassignment
				//switch aliveMessage.Elevator.Requests[floor][button].OrderState {
				//case elevator.SO_NoOrder:
				if aliveMessage.Elevator.Requests[floor][button].OrderState == elevator.SO_NoOrder {
					fmt.Println("Jeg har funnet en endring til SO_NoOrder")
					if localElevator.Requests[floor][button].ElevatorID == aliveMessage.ElevatorID &&
						localElevator.Requests[floor][button].OrderState == elevator.SO_Confirmed {
						fmt.Println("... og det kom fra den som eide orderen, så nå sletter jeg den")
						//localElevator = elevator.Requests_clearOnFloor(aliveMessage.ElevatorID, floor)
						panelPair := elevator.OrderpanelPair{ElevatorID: aliveMessage.ElevatorID, OrderState: elevator.SO_NoOrder}
						newChangedOrders = append(newChangedOrders, elevator.MakeOrderMessage(panelPair, elevio.ButtonEvent{Floor: floor, Button: button}))

						localElevator.Requests[floor][button].OrderState = elevator.SO_NoOrder
						localElevator.Requests[floor][button].ElevatorID = ""
						//Vi må legge noe append her sånn at det oppdateres utenfor scopet også
					} else if localElevator.Requests[floor][button].ElevatorID == aliveMessage.ElevatorID &&
					localElevator.Requests[floor][button].OrderState == elevator.SO_NewOrder{
						panelPair := elevator.OrderpanelPair{ElevatorID: aliveMessage.ElevatorID, OrderState: elevator.SO_NoOrder}
						newChangedOrders = append(newChangedOrders, elevator.MakeOrderMessage(panelPair, elevio.ButtonEvent{Floor: floor, Button: button}))

						localElevator.Requests[floor][button].OrderState = elevator.SO_NoOrder
						localElevator.Requests[floor][button].ElevatorID = ""
					}
					//Kanal til Requests_clearOnFloor()?
				} else if aliveMessage.Elevator.Requests[floor][button].OrderState == elevator.SO_NewOrder {
					fmt.Println("Jeg har funnet en endring til SO_NewOrder")
					if aliveMessage.Elevator.Requests[floor][button].ElevatorID == localElevator.ElevatorID {
						fmt.Println("...og orderen var fordelt til meg, så nå bekrefter jeg den")
						localElevator.Requests[floor][button].OrderState = elevator.SO_Confirmed
						localElevator.Requests[floor][button].ElevatorID = localElevator.ElevatorID

						panelPair := elevator.OrderpanelPair{ElevatorID: localElevator.ElevatorID, OrderState: elevator.SO_Confirmed}
						newChangedOrders = append(newChangedOrders, elevator.MakeOrderMessage(panelPair, elevio.ButtonEvent{Floor: floor, Button: button}))

						//confirmedOrderChan <- elevator.MakeOrderMessage(localElevator.ElevatorID, elevio.ButtonEvent{Floor: floor, Button: button})
					} else if aliveMessage.Elevator.Requests[floor][button].ElevatorID != localElevator.ElevatorID {
						fmt.Println("...men orderen var ikke til meg, så jeg bekrefter den ikke")
						localElevator.Requests[floor][button] = aliveMessage.Elevator.Requests[floor][button]
					}
				} else if aliveMessage.Elevator.Requests[floor][button].OrderState == elevator.SO_Confirmed {
					fmt.Println("Jeg har funnet en endring til SO_Confirmed")
					if aliveMessage.Elevator.Requests[floor][button].ElevatorID == aliveMessage.ElevatorID {
						fmt.Println("...og det er en bekreftelse fra den som eide orderen, så nå bekrefter jeg den også.")
						localElevator.Requests[floor][button].OrderState = elevator.SO_Confirmed

						panelPair := elevator.OrderpanelPair{ElevatorID: aliveMessage.ElevatorID, OrderState: elevator.SO_Confirmed}
						newChangedOrders = append(newChangedOrders, elevator.MakeOrderMessage(panelPair, elevio.ButtonEvent{Floor: floor, Button: button}))

						//localElevator.Requests[floor][button].ElevatorID = aliveMessage.ElevatorID
						//Kanal til Fsm_onRequestButtonPress()?
					} else {
						fmt.Println("...men det var jeg som eide denne orderen, så jeg bare chiller til den andre heisen har skjønt greia.")
					}
				}
				time.Sleep(25 * time.Millisecond)
			}
		}
	}

	return newChangedOrders
}
