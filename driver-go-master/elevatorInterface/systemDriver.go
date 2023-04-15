package elevatorInterface

import (
	"Driver-go/elevio"
	"Driver-go/manager"
	"Driver-go/singleElevator"
	"time"
)

/*

var database = manager.ElevatorDatabase{
	ConnectedElevators: 0,
}

func RunElevatorSystem(database manager.ElevatorDatabase) {
	doorTimer := time.NewTimer(3 * time.Second)
	doorTimer.Stop()

	immobilityTimer := time.NewTimer(3 * time.Second)
	immobilityTimer.Stop()

	for {
		select {
		//case: elevinput <-inputChannel
		//	SingleElevatorInput(elevinput)

		}

	}

}
*/

func HandleNewFloorAndUpdateDatabase(floor int, database manager.ElevatorDatabase, doorTimer *time.Timer, immobilityTimer *time.Timer) manager.ElevatorDatabase {
	var newElevatorUpdate singleElevator.Elevator
	newElevatorUpdate = singleElevator.Fsm_onFloorArrival(floor, doorTimer, immobilityTimer)
	database = manager.UpdateDatabase(newElevatorUpdate, database)
	return database
}

func HandleNewButtonAndUpdateDatabase(button elevio.ButtonEvent, database manager.ElevatorDatabase, doorTimer *time.Timer, immobilityTimer *time.Timer) manager.ElevatorDatabase {
	chosenElevator := manager.AssignOrderToElevator(database, button)
	newElevatorUpdate := singleElevator.HandleNewOrder(chosenElevator, button, doorTimer, immobilityTimer)
	database = manager.UpdateDatabase(newElevatorUpdate, database)
	return database
}

func HandleObstruction(obstruction bool, doorTimer *time.Timer, immobilityTimer *time.Timer) {
	if singleElevator.IsDoorOpen() && obstruction {
		doorTimer.Stop()
		immobilityTimer.Reset(3 * time.Second)
	} else if !obstruction && singleElevator.IsDoorOpen() {
		immobilityTimer.Stop()
		singleElevator.SetWorkingState(singleElevator.WS_Connected)
		doorTimer.Reset(3 * time.Second)
	}
}

func HandleStopButton(database manager.ElevatorDatabase) {
	singleElevator.ElevatorPrint(singleElevator.GetSingleEleavtorObject())
	for i := 0; i < len(database.ElevatorList); i++ {
		singleElevator.ElevatorPrint(database.ElevatorList[i])
	}
}

/*

func TimerInput(doorTimer *time.Timer, immobilityTimer *time.Timer) manager.ElevatorDatabase {
	for {
		select {
		case <-doorTimer.C:
			singleElevator.Fsm_onDoorTimeout(doorTimer)

		case <-immobilityTimer.C:
			fmt.Println("Iam immobile", singleElevator.MyID)
			database = manager.UpdateElevatorNetworkStateInDatabase(singleElevator.MyID, database, singleElevator.WS_Immobile)

			var deadOrders []elevio.ButtonEvent
			deadOrders = manager.FindDeadOrders(database, singleElevator.MyID)
			for j := 0; j < len(deadOrders); j++ {
				chosenElevator := manager.AssignOrderToElevator(database, deadOrders[j])
				newElevatorUpdate := singleElevator.HandleNewOrder(chosenElevator, deadOrders[j], doorTimer, immobilityTimer)
				database = manager.UpdateDatabase(newElevatorUpdate, database)
			}
			return database
		}
	}
}

func StateUpdateHandler(stateUpdateMessage singleElevator.ElevatorUpdateToDatabase, doorTimer *time.Timer, immobilityTimer *time.Timer) manager.ElevatorDatabase {
	if stateUpdateMessage.ElevatorID != singleElevator.MyID {
		database = manager.UpdateDatabase(stateUpdateMessage.Elevator, database)

		newChangedOrders := manager.SearchMessageForOrderUpdate(stateUpdateMessage, database)

		for i := 0; i < len(newChangedOrders); i++ {
			newOrder := newChangedOrders[i]
			var newElevatorUpdate singleElevator.Elevator

			if newOrder.PanelPair.OrderState == singleElevator.SO_Confirmed {
				chosenElevator := newOrder.PanelPair.ElevatorID
				newButton := newOrder.OrderedButton

				newElevatorUpdate = singleElevator.HandleConfirmedOrder(chosenElevator, newButton, doorTimer, immobilityTimer)

			} else if newOrder.PanelPair.OrderState == singleElevator.SO_NoOrder {
				fmt.Println("Inne no order ifen")
				newElevatorUpdate = singleElevator.Requests_clearOnFloor(newOrder.PanelPair.ElevatorID, newOrder.OrderedButton.Floor)
			}

			database = manager.UpdateDatabase(newElevatorUpdate, database)
		}

	}
	return database
}

func CabCallAssigner(newCabs manager.OrderStruct, database manager.ElevatorDatabase, doorTimer *time.Timer, immobilityTimer *time.Timer) manager.ElevatorDatabase {
	var newElevatorUpdate singleElevator.Elevator
	if newCabs.PanelPair.ElevatorID == singleElevator.MyID {
		newElevatorUpdate = singleElevator.Fsm_onRequestButtonPress(newCabs.OrderedButton.Floor, newCabs.OrderedButton.Button, singleElevator.MyID, doorTimer, immobilityTimer)
		database = manager.UpdateDatabase(newElevatorUpdate, database)
	}

	return database
}

func peerNetworkUpdateHandler(peerUpdate peers.PeerUpdate, database manager.ElevatorDatabase, cabsChannelTx chan manager.OrderStruct) manager.ElevatorDatabase {

	if len(peerUpdate.Lost) != 0 {
		var deadOrders []elevio.ButtonEvent
		for i := 0; i < len(peerUpdate.Lost); i++ {
			database = manager.UpdateElevatorNetworkStateInDatabase(p.Lost[i], database, singleElevator.WS_Unconnected)
			if database.ConnectedElevators <= 1 {
				singleElevator.SetIsAlone(true)
			}
			singleElevator.ElevatorPrint(manager.GetElevatorFromID(database, p.Lost[i]))
			deadOrders = manager.FindDeadOrders(database, p.Lost[i])
			singleElevator.ElevatorPrint(manager.GetElevatorFromID(database, p.Lost[i]))

		}

		for j := 0; j < len(deadOrders); j++ {
			chosenElevator := manager.AssignOrderToElevator(database, deadOrders[j])
			newElevatorUpdate := singleElevator.HandleNewOrder(chosenElevator, deadOrders[j], doorTimer, immobilityTimer)
			database = manager.UpdateDatabase(newElevatorUpdate, database)
		}

	}

	if p.New != "" {
		if !singleElevator.GetIsAlone() {
			cabsToBeSent := manager.FindCabCallsForElevator(database, peerUpdate.New)
			fmt.Println("Ready to send the following CABs:", cabsToBeSent)
			for k := 0; k < len(cabsToBeSent); k++ {
				cabsChannelTx <- cabsToBeSent[k]
				time.Sleep(time.Duration(25) * time.Millisecond)
			}
		}

		if !manager.IsElevatorInDatabase(peerUpdate.New, database) {
			database.ElevatorList = append(database.ElevatorList, singleElevator.Elevator{ElevatorID: p.New, Operating: singleElevator.WS_Connected})
		}

		database = manager.UpdateElevatorNetworkStateInDatabase(peerUpdate.New, database, singleElevator.WS_Connected)
		if database.ConnectedElevators > 1 {
			singleElevator.SetIsAlone(false)
		}

	}

}
*/
