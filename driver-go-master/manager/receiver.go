package manager

import (
	"Driver-go/elevator"
	"Driver-go/elevio"
	"fmt"
	"time"
)

func ReceiveMessages(mainTimer time.Timer, database ElevatorDatabase, databaseChan chan ElevatorDatabase,
	ackTx chan AckMessageStruct, orderRx chan OrderMessageStruct, floorArrivalRx chan FloorArrivalMessageStruct,
	ackRx chan AckMessageStruct, newElevatorRx chan SingleElevatorMessageStruct) {

	acknowledgementCount := 0

	for {
		select {
		case floorArrivalBroadcast := <-floorArrivalRx:
			fmt.Println("I have received a floor arrival!")
			ackTx <- SendAcknowledge(floorArrivalBroadcast.SenderID)
			databaseChan <- UpdateDatabase(floorArrivalBroadcast.MyElevator, database, floorArrivalBroadcast.SenderID)

			if floorArrivalBroadcast.SenderID != elevator.MyID {
				elevator.Requests_clearOnFloor(floorArrivalBroadcast.SenderID, floorArrivalBroadcast.ArrivedFloor)
			}

		case orderBroadcast := <-orderRx:
			fmt.Println("I have received an order!")
			ackTx <- SendAcknowledge(orderBroadcast.SenderID)
			updatedDatabase := UpdateDatabase(orderBroadcast.MyElevator, database, orderBroadcast.SenderID)
			databaseChan <- updatedDatabase

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
			if WhatFloorIsElevatorFromStringID(updatedDatabase, orderBroadcast.ChosenElevator) == orderBroadcast.OrderedButton.Floor &&
				WhatStateIsElevatorFromStringID(updatedDatabase, orderBroadcast.ChosenElevator) != elevator.EB_Moving {
				elevator.Requests_clearOnFloor(orderBroadcast.ChosenElevator, orderBroadcast.OrderedButton.Floor)
			}

		case newElevatorBroadcast := <-newElevatorRx:
			ackTx <- SendAcknowledge(newElevatorBroadcast.SenderID)
			if !IsElevatorInDatabase(newElevatorBroadcast.SenderID, database) {
				database.ElevatorsInNetwork = append(database.ElevatorsInNetwork, newElevatorBroadcast.MyElevator)
				database.NumElevators++
				fmt.Println(" number of elevators ---  ")
				fmt.Println(database.NumElevators)
				databaseChan <- database
			} else {
				databaseChan <- UpdateDatabase(newElevatorBroadcast.MyElevator, database, newElevatorBroadcast.SenderID)
			}
		case ack := <-ackRx:
			//Hvis det er en ack som er til meg
			if ack.IDOfAckReciever == elevator.MyID {
				acknowledgementCount++ //må man vite hvem som har sendt acks?
				fmt.Println("Order_Number of acknowledgements: ", acknowledgementCount)
				if acknowledgementCount == 2 { //numOfExpectedAcks
					fmt.Println("I have received all acks")
					acknowledgementCount = 0
					//timer.Stop()
					//break
				}
			}
		}
		time.Sleep(time.Duration(20) * time.Millisecond)
	}
}
