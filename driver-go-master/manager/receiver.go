package manager

import (
	"Driver-go/elevator"
	"Driver-go/elevio"
	"fmt"
	"time"
)

func ReceiveMessages(msgRx chan MessageStruct, ackRx chan AckMessageStruct,
	database ElevatorDatabase, updateDatabaseChan chan ElevatorDatabase, mainTimer time.Timer,
	receivedAckChan chan AckMessageStruct, initiateSendAckChan chan AckMessageStruct) {
	for {
		select {
		case messageBroadcast := <-msgRx:
			switch messageBroadcast.MessageType {
			case "FloorArrival":
				fmt.Println("received floor")
				initiateSendAckChan <- SendAcknowledge(messageBroadcast)
				updateDatabaseChan <- UpdateDatabase(messageBroadcast, database)

				if messageBroadcast.SenderID == elevator.MyID {
					elevator.Fsm_onFloorArrival(messageBroadcast.ArrivedFloor, &mainTimer)
				} else {
					elevator.Requests_clearOnFloor(messageBroadcast.SenderID, messageBroadcast.ArrivedFloor)
				}

			case "Order":
				fmt.Println("The receiver has received an order at floor: ", messageBroadcast.OrderedButton.Floor)
				initiateSendAckChan <- SendAcknowledge(messageBroadcast)
				fmt.Println("The need for sending an acknowledgment has been sent.")
				updateDatabaseChan <- UpdateDatabase(messageBroadcast, database)

				if messageBroadcast.OrderedButton.Button != elevio.BT_Cab {
					elevator.Elevator_increaseOrderNumber()
				}

				//if chosenElev already on floor -> Request_clearOnFloor
				if (messageBroadcast.OrderedButton.Button == elevio.BT_Cab && messageBroadcast.ChosenElevator == elevator.MyID) ||
					messageBroadcast.OrderedButton.Button != elevio.BT_Cab {
					elevator.Fsm_onRequestButtonPress(messageBroadcast.OrderedButton.Floor, messageBroadcast.OrderedButton.Button,
						messageBroadcast.ChosenElevator, &mainTimer)
				}

				//HER LA VI TIL EN SJEKK OM CHOSEN ELEVTAOR ER I ETASJEN TIL BESTILLINGEN ALLEREDE, hvis den er det skal bestillingen cleares med en gang.
				//burde sikkert være innbakt et annet sted.
				if WhatFloorIsElevatorFromStringID(database, messageBroadcast.ChosenElevator) == messageBroadcast.OrderedButton.Floor &&
					WhatStateIsElevatorFromStringID(database, messageBroadcast.ChosenElevator) != elevator.EB_Moving {
					elevator.Requests_clearOnFloor(messageBroadcast.ChosenElevator, messageBroadcast.OrderedButton.Floor)
				}
			case "NewElevatorOnNetwork":
				if !IsElevatorInDatabase(messageBroadcast.SenderID, database) {
					database.ElevatorsInNetwork = append(database.ElevatorsInNetwork, messageBroadcast.MyElevator)
					database.NumElevators++
					fmt.Println(" number of elevators ---  ")
					fmt.Println(database.NumElevators)
				}
			default:
				fmt.Println("I have received a message that I don't know hoe to handle.")
			}
		case ack := <-ackRx:
			receivedAckChan <- ack
		}
	}
}
