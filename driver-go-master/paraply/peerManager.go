package paraply

import (
	"Driver-go/elevator"
	"Driver-go/elevio"
	"Driver-go/manager"
	"Driver-go/network/peers"
	"time"
)

func ManagePeers(p peers.PeerUpdate, database manager.ElevatorDatabase, doorTimer *time.Timer, immobilityTimer *time.Timer, cabsChannelTx chan elevator.OrderStruct) manager.ElevatorDatabase {
	if len(p.Lost) != 0 {
		handleLostPeer(p.Lost, database, doorTimer, immobilityTimer)
	}

	if p.New != "" {
		handleNewPeer(p.New, database)
		sendLostCabs(cabsChannelTx, p.New, database)
	}
	return database
}

func sendLostCabs(cabsChannelTx chan elevator.OrderStruct, newPeer string, database manager.ElevatorDatabase) {
	if !elevator.GetIAmAlone() {
		cabsToBeSent := manager.FindCabCallsForElevator(database, newPeer)
		for k := 0; k < len(cabsToBeSent); k++ {
			cabsChannelTx <- cabsToBeSent[k]
			time.Sleep(time.Duration(25) * time.Millisecond)
		}
	}

}

func handleNewPeer(newPeer string, database manager.ElevatorDatabase) manager.ElevatorDatabase {
	newPeerUpdate := manager.GetElevatorFromID(database, elevator.MyID)

	if !manager.IsElevatorInDatabase(newPeer, database) {
		database.ElevatorList = append(database.ElevatorList, elevator.Elevator{ElevatorID: newPeer, Operating: elevator.WS_Connected})
	}

	database = manager.UpdateElevatorNetworkStateInDatabase(newPeer, database, elevator.WS_Connected)
	if database.ConnectedElevators > 1 {
		newPeerUpdate = elevator.SetIAmAlone(false, newPeerUpdate)
		database = manager.UpdateDatabase(newPeerUpdate, database)
	}

	return database

}

func handleLostPeer(lostPeers []string, database manager.ElevatorDatabase, doorTimer *time.Timer, immobilityTimer *time.Timer) manager.ElevatorDatabase {
	newPeerUpdate := manager.GetElevatorFromID(database, elevator.MyID)

	var deadOrders []elevio.ButtonEvent

	for i := 0; i < len(lostPeers); i++ {
		database = manager.UpdateElevatorNetworkStateInDatabase(lostPeers[i], database, elevator.WS_Unconnected)

		if database.ConnectedElevators <= 1 {
			newPeerUpdate = elevator.SetIAmAlone(true, newPeerUpdate)
			database = manager.UpdateDatabase(newPeerUpdate, database)
		}
		deadOrders = manager.FindDeadOrders(database, lostPeers[i])
	}

	for j := 0; j < len(deadOrders); j++ {
		chosenElevator := manager.AssignOrderToElevator(database, deadOrders[j])
		newElevatorUpdate := elevator.HandleNewOrder(chosenElevator, deadOrders[j], doorTimer, immobilityTimer)
		database = manager.UpdateDatabase(newElevatorUpdate, database)
	}
	return database
}
