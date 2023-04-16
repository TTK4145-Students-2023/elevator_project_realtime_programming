package databaseHandler

import (
	"Driver-go/elevatorHardware"
	"Driver-go/singleElevator"
)

//Functions to analyze incoming stateUpdateMessages from other elevators. 
//Used to find differences in order states between local database and incoming message.

func FindChangesBetweenIncomingMessageAndLocalDatabase(stateUpdateMessage singleElevator.ElevatorStateUpdate, database ElevatorDatabase) []OrderStruct { // comparemessagewithlocaldatabase

	var DifferencesFound []OrderStruct

	localElevator := GetElevatorFromID(database, singleElevator.MyID)

	receivedElevatorID := stateUpdateMessage.ElevatorID

	for floor := 0; floor < singleElevator.NumFloors; floor++ {
		for button := elevatorHardware.BT_HallUp; button < elevatorHardware.BT_Cab; button++ {

			currentButtonEvent := elevatorHardware.ButtonEvent{Floor: floor, Button: button}

			receivedOrderState := stateUpdateMessage.Elevator.Requests[floor][button].OrderState
			localOrderState := localElevator.Requests[floor][button].OrderState

			receivedRequestID := stateUpdateMessage.Elevator.Requests[floor][button].AssingedElevatorID
			localRequestID := localElevator.Requests[floor][button].AssingedElevatorID

			if receivedOrderState != localOrderState {
				DifferencesFound = CompareIncomningOrderStateAndLocalOrderState(receivedElevatorID, currentButtonEvent, receivedOrderState, localOrderState, receivedRequestID, localRequestID)
			}
		}

	}
	return DifferencesFound
}

func CompareIncomningOrderStateAndLocalOrderState(receivedElevatorID string, currentButtonEvent elevatorHardware.ButtonEvent, receivedOrderState singleElevator.StateOfOrder, localOrderState singleElevator.StateOfOrder, receivedRequestID string, localRequestID string) []OrderStruct {

	var DifferencesFound []OrderStruct
	localElevatorID := singleElevator.MyID

	if receivedOrderState == singleElevator.NoOrder {

		if localRequestID == receivedElevatorID &&
			localOrderState == singleElevator.ConfirmedOrder {

			panelPair := singleElevator.StateAndChosenElevator{AssingedElevatorID: receivedElevatorID, OrderState: singleElevator.NoOrder}
			DifferencesFound = append(DifferencesFound, MakeOrder(panelPair, currentButtonEvent))

		} else if localRequestID == receivedElevatorID &&
			localOrderState == singleElevator.NewOrder {

			panelPair := singleElevator.StateAndChosenElevator{AssingedElevatorID: receivedElevatorID, OrderState: singleElevator.NoOrder}
			DifferencesFound = append(DifferencesFound, MakeOrder(panelPair, currentButtonEvent))

		}
	} else if receivedOrderState == singleElevator.NewOrder {

		if receivedRequestID == localElevatorID {
			panelPair := singleElevator.StateAndChosenElevator{AssingedElevatorID: localElevatorID, OrderState: singleElevator.ConfirmedOrder}
			DifferencesFound = append(DifferencesFound, MakeOrder(panelPair, currentButtonEvent))

		} else {
			panelPair := singleElevator.StateAndChosenElevator{AssingedElevatorID: receivedElevatorID, OrderState: singleElevator.NewOrder}
			DifferencesFound = append(DifferencesFound, MakeOrder(panelPair, currentButtonEvent))
		}
	} else if receivedOrderState == singleElevator.ConfirmedOrder {

		if receivedRequestID == receivedElevatorID {
			panelPair := singleElevator.StateAndChosenElevator{AssingedElevatorID: receivedElevatorID, OrderState: singleElevator.ConfirmedOrder}
			DifferencesFound = append(DifferencesFound, MakeOrder(panelPair, currentButtonEvent))

		}
	}
	return DifferencesFound
}
