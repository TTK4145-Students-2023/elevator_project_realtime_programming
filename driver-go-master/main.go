package main

import (
	"Driver-go/databaseHandler"
	"Driver-go/elevatorHardware"
	"Driver-go/elevatorInterface"
	"Driver-go/network/bcast"
	"Driver-go/network/localip"
	"Driver-go/network/peers"
	"Driver-go/orderDelegation"
	"Driver-go/peerStatus"
	"Driver-go/singleElevator"
	"flag"
	"fmt"
	"os"
	"time"
)

const nFloors = 4

const inputPollRateMs = 25

func main() {

	//hardware channel init
	

	buttonChannel := make(chan elevatorHardware.ButtonEvent)
	floorSensorChannel := make(chan int)
	obstructionChannel := make(chan bool)

	go elevatorHardware.PollButtons(buttonChannel)
	go elevatorHardware.PollFloorSensor(floorSensorChannel)
	go elevatorHardware.PollObstructionSwitch(obstructionChannel)

	
	//network init
	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}

	elevatorHardware.Init("localhost:"+id, nFloors) 
	
	if elevatorHardware.GetFloor() == -1 {
		singleElevator.Fsm_onInitBetweenFloors()
	}


	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)

	go peers.Transmitter(15600, id, peerTxEnable) //15647
	go peers.Receiver(15600, peerUpdateCh)

	restoredCabsChannelRx := make(chan databaseHandler.OrderStruct)
	restoredCabsChannelTx := make(chan databaseHandler.OrderStruct)
	stateUpdateTx := make(chan singleElevator.ElevatorStateUpdate)
	stateUpdateRx := make(chan singleElevator.ElevatorStateUpdate)

	go bcast.Transmitter(16569, restoredCabsChannelTx, stateUpdateTx)
	go bcast.Receiver(16569, restoredCabsChannelRx, stateUpdateRx) //port: 16569

	//endre denne for å bruke flere sockets for elevcd //15657

	go singleElevator.TransmitStateUpdate(stateUpdateTx)

	//database init
	database := databaseHandler.ElevatorDatabase{
		ConnectedElevators: 0,
	}

	//timer init
	doorTimer := time.NewTimer(3 * time.Second)
	doorTimer.Stop()

	immobilityTimer := time.NewTimer(3 * time.Second)
	immobilityTimer.Stop()

	for {

		select {
		case floor := <-floorSensorChannel:
			database = elevatorInterface.HandleNewFloorAndUpdateDatabase(floor, database, doorTimer, immobilityTimer)

		case button := <-buttonChannel:
			database = elevatorInterface.HandleNewButtonAndUpdateDatabase(button, database, doorTimer, immobilityTimer)

		case obstruction := <-obstructionChannel:
			elevatorInterface.HandleObstruction(obstruction, doorTimer, immobilityTimer)

		case stateUpdateMessage := <-stateUpdateRx:
			if !databaseHandler.MessageIDEqualsMyID(stateUpdateMessage.ElevatorID) {
				database = databaseHandler.UpdateDatabaseFromIncomingMessages(stateUpdateMessage, database, immobilityTimer, doorTimer)
			}

		case restoredCabs := <-restoredCabsChannelRx:
			newDatabaseEntry := orderDelegation.HandleRestoredCabs(restoredCabs, doorTimer, immobilityTimer)
			database = databaseHandler.UpdateDatabase(newDatabaseEntry, database)

		case peerUpdateInfo := <-peerUpdateCh:
			lostPeers := peerUpdateInfo.Lost
			newPeer := peerUpdateInfo.New

			if len(lostPeers) != 0 {
				database = peerStatus.HandlePeerLoss(lostPeers, database, immobilityTimer, doorTimer)
			}

			if newPeer != "" {
				database = peerStatus.HandleNewPeer(newPeer, database, restoredCabsChannelTx)
			}

		case <-doorTimer.C:
			singleElevator.DoorTimeout(doorTimer)

		case <-immobilityTimer.C:
			database = databaseHandler.UpdateElevatorNetworkStateInDatabase(singleElevator.MyID, database, singleElevator.Immobile)

			database = databaseHandler.UpdateDatabaseWithDeadOrders(singleElevator.MyID, immobilityTimer, doorTimer, database)

		}

		time.Sleep(time.Duration(inputPollRateMs) * time.Millisecond)
	}

}
