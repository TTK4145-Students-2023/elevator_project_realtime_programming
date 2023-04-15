package peerStatus

import (
	"Driver-go/databaseHandler"
	"Driver-go/orderDelegation"
	"Driver-go/singleElevator"
	"time"
)

func HandlePeerLoss(lostPeers []string, database databaseHandler.ElevatorDatabase, immobilityTimer *time.Timer, doorTimer *time.Timer) databaseHandler.ElevatorDatabase {

	for i := 0; i < len(lostPeers); i++ {
		database = databaseHandler.UpdateElevatorNetworkStateInDatabase(lostPeers[i], database, singleElevator.Unconnected)
		if database.ConnectedElevators <= 1 {
			singleElevator.SetIsAlone(true)
		}
		database = databaseHandler.UpdateDatabaseWithDeadOrders(lostPeers[i], immobilityTimer, doorTimer, database)

	}

	return database

}

func HandleNewPeer(newPeer string, database databaseHandler.ElevatorDatabase, cabsChannelTx chan databaseHandler.OrderStruct) databaseHandler.ElevatorDatabase {
	if !singleElevator.GetIsAlone() {
		cabsToBeSent := orderDelegation.FindCabCallsForElevator(database, newPeer)
		orderDelegation.SendCabCalls(cabsToBeSent, cabsChannelTx)
	}

	if !databaseHandler.IsElevatorInDatabase(newPeer, database) {
		database.ElevatorList = append(database.ElevatorList, singleElevator.Elevator{ElevatorID: newPeer, Operating: singleElevator.Connected})
	}

	database = databaseHandler.UpdateElevatorNetworkStateInDatabase(newPeer, database, singleElevator.Connected)
	if database.ConnectedElevators > 1 {
		singleElevator.SetIsAlone(false)
	}
	return database
}
