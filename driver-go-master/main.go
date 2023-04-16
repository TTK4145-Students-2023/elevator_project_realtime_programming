package main

import (
	"Driver-go/databaseHandler"
	"Driver-go/elevatorHardware"
	"Driver-go/network/bcast"
	"Driver-go/network/localip"
	"Driver-go/network/peers"
	"Driver-go/peerUpdateHandler"
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

	//tcp connection init
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
		singleElevator.InitializeElevatorBetweenFloors()
	}

	//network init
	peerUpdateChannel := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)

	go peers.Transmitter(15600, id, peerTxEnable) //15647
	go peers.Receiver(15600, peerUpdateChannel)

	restoredCabsChannelRx := make(chan databaseHandler.OrderStruct)
	restoredCabsChannelTx := make(chan databaseHandler.OrderStruct)
	stateUpdateChannelTx := make(chan singleElevator.ElevatorStateUpdate)
	stateUpdateChannelRx := make(chan singleElevator.ElevatorStateUpdate)

	go bcast.Transmitter(16569, restoredCabsChannelTx, stateUpdateChannelTx)
	go bcast.Receiver(16569, restoredCabsChannelRx, stateUpdateChannelRx) //port: 16569

	//endre denne for å bruke flere sockets for elevcd //15657

	go singleElevator.TransmitStateUpdate(stateUpdateChannelTx)

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
			database = databaseHandler.HandleNewFloorAndUpdateDatabase(floor, database, doorTimer, immobilityTimer)

		case button := <-buttonChannel:
			database = databaseHandler.HandleNewButtonAndUpdateDatabase(button, database, doorTimer, immobilityTimer)

		case obstruction := <-obstructionChannel:
			singleElevator.HandleObstruction(obstruction, doorTimer, immobilityTimer)

		case stateUpdateMessage := <-stateUpdateChannelRx:
			if !databaseHandler.MessageIDEqualsMyID(stateUpdateMessage.ElevatorID) {
				database = databaseHandler.UpdateDatabaseFromIncomingMessages(stateUpdateMessage, database, immobilityTimer, doorTimer)
			}

		case restoredCabs := <-restoredCabsChannelRx:
			newDatabaseEntry := peerUpdateHandler.HandleRestoredCabs(restoredCabs, doorTimer, immobilityTimer)
			database = databaseHandler.UpdateDatabase(newDatabaseEntry, database)

		case peerUpdateInfo := <-peerUpdateChannel:
			lostPeers := peerUpdateInfo.Lost
			newPeer := peerUpdateInfo.New

			if len(lostPeers) != 0 {
				database = peerUpdateHandler.HandlePeerLoss(lostPeers, database, immobilityTimer, doorTimer)
			}

			if newPeer != "" {
				database = peerUpdateHandler.HandleNewPeer(newPeer, database, restoredCabsChannelTx)
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
