package main

import (
	"Driver-go/elevatorInterface"
	"Driver-go/elevio"
	"Driver-go/manager"
	"Driver-go/network/bcast"
	"Driver-go/network/localip"
	"Driver-go/network/peers"
	"Driver-go/singleElevator"
	"flag"
	"fmt"
	"os"
	"time"
)

const nFloors = 4

//const nButtons = 3

func main() {
	

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


	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)
	
	go peers.Transmitter(15600, id, peerTxEnable) //15647
	go peers.Receiver(15600, peerUpdateCh)

	cabsChannelTx := make(chan manager.OrderStruct)
	cabsChannelRx := make(chan manager.OrderStruct)
	stateUpdateTx := make(chan singleElevator.ElevatorStateUpdate)
	stateUpdateRx := make(chan singleElevator.ElevatorStateUpdate)
	
	go bcast.Transmitter(16569, cabsChannelTx, stateUpdateTx)
	go bcast.Receiver(16569, cabsChannelRx, stateUpdateRx)  //port: 16569
	
	elevio.Init("localhost:"+id, nFloors) //endre denne for å bruke flere sockets for elevcd //15657

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	if elevio.GetFloor() == -1 {
		singleElevator.Fsm_onInitBetweenFloors()
	}


	database := manager.ElevatorDatabase{
		ConnectedElevators: 0,
	}
	
	doorTimer := time.NewTimer(3 * time.Second)
	doorTimer.Stop()
	
	immobilityTimer := time.NewTimer(3 * time.Second)
	immobilityTimer.Stop()
	
	var inputPollRateMs = 25
	

	go singleElevator.TransmittStateUpdate(stateUpdateTx) //TransmitStateUpdate



	for {

		select {
		case floor := <-drv_floors:
			database = elevatorInterface.HandleNewFloorAndUpdateDatabase(floor, database, doorTimer, immobilityTimer)

		case button := <-drv_buttons:
			database = elevatorInterface.HandleNewButtonAndUpdateDatabase(button, database, doorTimer, immobilityTimer)

		case obstruction := <-drv_obstr:
			elevatorInterface.HandleObstruction(obstruction, doorTimer, immobilityTimer)

		case <-drv_stop:
			elevatorInterface.HandleStopButton(database) //Extra print function

		case <-doorTimer.C:
			singleElevator.Fsm_onDoorTimeout(doorTimer)

		case <-immobilityTimer.C:
			database = manager.UpdateElevatorNetworkStateInDatabase(singleElevator.MyID, database, singleElevator.WS_Immobile)

			database = manager.UpdateDatabaseWithDeadOrders(singleElevator.MyID, immobilityTimer, doorTimer, database)

		case stateUpdateMessage := <-stateUpdateRx:
			if !manager.MessageIDEqualsMyID(stateUpdateMessage.ElevatorID) {
				database = manager.UpdateDatabaseFromIncomingMessages(stateUpdateMessage, database, immobilityTimer, doorTimer)
			}

		case newCabs := <-cabsChannelRx:
			newElevatorUpdate := manager.HandleRestoredCabs(newCabs, doorTimer, immobilityTimer)
			database = manager.UpdateDatabase(newElevatorUpdate, database)

		case p := <-peerUpdateCh:

			if len(p.Lost) != 0 {
				database = manager.HandlePeerLoss(p.Lost, database, immobilityTimer, doorTimer)
			}

			if p.New != "" {
				database = manager.HandleNewPeer(p.New, database, cabsChannelTx)
			}

		}

		time.Sleep(time.Duration(inputPollRateMs) * time.Millisecond)
	}

}
