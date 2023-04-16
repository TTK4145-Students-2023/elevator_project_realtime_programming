package connectionHandler

import (
	"Driver-go/databaseHandler"
	"Driver-go/singleElevator"
	"time"
)
//Functions to handle elevator disconnects or reconnects in the network.
//Every elevator is a peer.


func HandleDisconnectedPeer(lostPeers []string, database databaseHandler.ElevatorDatabase, immobilityTimer *time.Timer, doorTimer *time.Timer) databaseHandler.ElevatorDatabase {

	for i := 0; i < len(lostPeers); i++ {
		database = databaseHandler.UpdateElevatorNetworkStateInDatabase(singleElevator.Unconnected, lostPeers[i], database)
		if database.ConnectedElevators <= 1 {
			singleElevator.SetIsAlone(true)
		}
		database = databaseHandler.UpdateDatabaseWithDeadOrders(lostPeers[i], immobilityTimer, doorTimer, database)

	}

	return database

}

func HandleReconnectedPeer(newPeer string, database databaseHandler.ElevatorDatabase, cabsChannelTx chan databaseHandler.OrderStruct) databaseHandler.ElevatorDatabase {
	if !singleElevator.GetIsAlone() {
		cabsToBeSent := FindCabCallsForElevator(database, newPeer)
		SendCabCalls(cabsToBeSent, cabsChannelTx)
	}

	if !databaseHandler.IsElevatorInDatabase(newPeer, database) {
		database.ElevatorList = append(database.ElevatorList, singleElevator.Elevator{ElevatorID: newPeer, Operating: singleElevator.Connected})
	}

	database = databaseHandler.UpdateElevatorNetworkStateInDatabase(singleElevator.Connected, newPeer, database)
	if database.ConnectedElevators > 1 {
		singleElevator.SetIsAlone(false)
	}
	return database
}
