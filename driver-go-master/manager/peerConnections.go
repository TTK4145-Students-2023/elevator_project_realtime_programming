package manager

import (
	"Driver-go/singleElevator"
	"time"
)

func HandlePeerLoss(lostPeers []string, database ElevatorDatabase, immobilityTimer *time.Timer, doorTimer *time.Timer) ElevatorDatabase {

	for i := 0; i < len(lostPeers); i++ {
		database = UpdateElevatorNetworkStateInDatabase(lostPeers[i], database, singleElevator.WS_Unconnected)
		if database.ConnectedElevators <= 1 {
			singleElevator.SetIsAlone(true)
		}
		database = UpdateDatabaseWithDeadOrders(lostPeers[i], immobilityTimer, doorTimer, database)

	}

	return database

}

func HandleNewPeer(newPeer string, database ElevatorDatabase, cabsChannelTx chan OrderStruct) ElevatorDatabase{
	if !singleElevator.GetIsAlone() {
		cabsToBeSent := FindCabCallsForElevator(database, newPeer)
		SendCabCalls(cabsToBeSent, cabsChannelTx)
	}

	if !IsElevatorInDatabase(newPeer, database) {
		database.ElevatorList = append(database.ElevatorList, singleElevator.Elevator{ElevatorID: newPeer, Operating: singleElevator.WS_Connected})
	}

	database = UpdateElevatorNetworkStateInDatabase(newPeer, database, singleElevator.WS_Connected)
	if database.ConnectedElevators > 1 {
		singleElevator.SetIsAlone(false)
	}
	return database
}


