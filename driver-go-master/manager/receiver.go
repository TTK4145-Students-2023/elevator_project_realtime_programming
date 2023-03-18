package manager

import (
	"Driver-go/elevator"
	"Driver-go/elevio"
	"fmt"
	"time"
)

func ReceiveMessages(mainTimer time.Timer, database ElevatorDatabase, databaseChan chan ElevatorDatabase,
	ackTx chan AckMessageStruct, orderRx chan OrderMessageStruct, floorArrivalRx chan FloorArrivalMessageStruct,
	ackRx chan AckMessageStruct, newElevatorRx chan SingleElevatorMessageStruct, shouldAck chan AckMessageStruct) {

	for {
		select {
		case floorArrivalBroadcast := <-floorArrivalRx:
			
			shouldAck <- SendAcknowledge(floorArrivalBroadcast.SenderID)
			databaseChan <- UpdateDatabase(floorArrivalBroadcast.MyElevator, database, floorArrivalBroadcast.SenderID)

			if floorArrivalBroadcast.SenderID == elevator.MyID {
				elevator.Fsm_onFloorArrival(floorArrivalBroadcast.ArrivedFloor, &mainTimer)
			} else {
				elevator.Requests_clearOnFloor(floorArrivalBroadcast.SenderID, floorArrivalBroadcast.ArrivedFloor)
			}

		case orderBroadcast := <-orderRx:
			fmt.Println("I have received an order!")
			shouldAck <- SendAcknowledge(orderBroadcast.SenderID)
			databaseChan <- UpdateDatabase(orderBroadcast.MyElevator, database, orderBroadcast.SenderID)

			if orderBroadcast.OrderedButton.Button != elevio.BT_Cab {
				elevator.Elevator_increaseOrderNumber()
			}

			//if chosenElev already on floor -> Request_clearOnFloor
			if (orderBroadcast.OrderedButton.Button == elevio.BT_Cab && orderBroadcast.ChosenElevator == elevator.MyID) ||
				orderBroadcast.OrderedButton.Button != elevio.BT_Cab {
				elevator.Fsm_onRequestButtonPress(orderBroadcast.OrderedButton.Floor, orderBroadcast.OrderedButton.Button, orderBroadcast.ChosenElevator, &mainTimer)
			}

			//HER LA VI TIL EN SJEKK OM CHOSEN ELEVTAOR ER I ETASJEN TIL BESTILLINGEN ALLEREDE, hvis den er det skal bestillingen cleares med en gang.
			//burde sikkert være innbakt et annet sted.
			if WhatFloorIsElevatorFromStringID(database, orderBroadcast.ChosenElevator) == orderBroadcast.OrderedButton.Floor &&
				WhatStateIsElevatorFromStringID(database, orderBroadcast.ChosenElevator) != elevator.EB_Moving {
				elevator.Requests_clearOnFloor(orderBroadcast.ChosenElevator, orderBroadcast.OrderedButton.Floor)
			}

		case newElevatorBroadcast := <-newElevatorRx:
			if !IsElevatorInDatabase(newElevatorBroadcast.SenderID, database) {
				database.ElevatorsInNetwork = append(database.ElevatorsInNetwork, newElevatorBroadcast.MyElevator)
				database.NumElevators++
				fmt.Println(" number of elevators ---  ")
				fmt.Println(database.NumElevators)
				databaseChan <- database
			} else {
				databaseChan <- UpdateDatabase(newElevatorBroadcast.MyElevator, database, newElevatorBroadcast.SenderID)
			}
		}
		time.Sleep(time.Duration(25) * time.Millisecond)
	}
}
