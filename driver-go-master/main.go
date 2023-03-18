package main

import (
	"Driver-go/elevator"
	"Driver-go/elevio"
	"Driver-go/manager"
	"Driver-go/network/bcast"
	"Driver-go/network/localip"
	"Driver-go/network/peers"
	"flag"
	"fmt"
	"os"
	"time"
)

const nFloors = 4

//const nButtons = 3

func main() {
	// Our id can be anything. Here we pass it on the command line, using
	//  `go run main.go -id=our_id`

	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	// ... or alternatively, we can use the local IP address.
	// (But since we can run multiple programs on the same PC, we also append the
	//  process ID)
	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}

	// We make a channel for receiving updates on the id's of the peers that are
	//  alive on the network
	peerUpdateCh := make(chan peers.PeerUpdate)
	// We can disable/enable the transmitter after it has been started.
	// This could be used to signal that we are somehow "unavailable".
	peerTxEnable := make(chan bool)
	go peers.Transmitter(15600, id, peerTxEnable) //15647
	go peers.Receiver(15600, peerUpdateCh)

	/*
		orderTx := make(chan manager.MessageStruct)
		orderRx := make(chan manager.MessageStruct)

		floorArrivalTx := make(chan manager.MessageStruct)
		floorArrivalRx := make(chan manager.MessageStruct)
	*/
	msgTx := make(chan manager.MessageStruct)
	msgRx := make(chan manager.MessageStruct)

	ackTx := make(chan manager.AckMessageStruct) //Disse kanalene for å sende acks
	ackRx := make(chan manager.AckMessageStruct)

	go bcast.Transmitter(11771, msgTx, ackTx)
	go bcast.Receiver(11771, msgRx, ackRx)
	//port: 16569

	fmt.Println("Started!")

	mainDatabase := manager.ElevatorDatabase{
		//hardkodede verdier vi alltid bruker når vi flagger
		NumElevators: 0,
	}

	mainTimer := time.NewTimer(3 * time.Second)
	mainTimer.Stop()

	inputPollRateMs := 25

	elevio.Init("localhost:"+id, nFloors) //endre denne for å bruke flere sockets for elevcd //15657

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	updateDatabaseChan := make(chan manager.ElevatorDatabase)
	receivedAckChan := make(chan manager.AckMessageStruct) //ack communication between receiver and transmitter
	initiateSendAckChan := make(chan manager.AckMessageStruct)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	go manager.ReceiveMessages(msgRx, ackRx, mainDatabase, updateDatabaseChan, *mainTimer, receivedAckChan, initiateSendAckChan)
	go manager.TransmitMessages(drv_buttons, drv_floors, msgTx, ackTx, mainDatabase, receivedAckChan, initiateSendAckChan)

	//input := elevioGetInputDevice()

	if elevio.GetFloor() == -1 {
		elevator.Fsm_onInitBetweenFloors()
	}

	//prev := [nFloors][nButtons]int{}

	//lage en for select hvor drv_button sender istedet for request buttons
	for {

		select {
		case obstruction := <-drv_obstr:
			if elevator.IsDoorOpen() && obstruction {
				mainTimer.Stop()
			} else if !obstruction && elevator.IsDoorOpen() {
				mainTimer.Reset(3 * time.Second)
			}

		case timedOut := <-mainTimer.C:

			fmt.Print(timedOut)
			fmt.Println("INNNNNNNNNNE I TIMEDOUT")

			elevator.Fsm_onDoorTimeout(mainTimer)
			//SendFloorArrival(floorMsg,floorArrivalTx)

		case updatedDatabase := <-updateDatabaseChan:
			mainDatabase = updatedDatabase
			fmt.Println("The database has been updated!")

		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

			manager.UpdateElevatorNetworkStateInDatabase(p, mainDatabase)

			//legg dette inn i updatenetwork state
			if len(p.Lost) != 0 {
				for i := 0; i < len(p.Lost); i++ {
					manager.ReassignDeadOrders(msgTx, mainDatabase, p.Lost[i])
				}
			}

			if p.New != "" {
				msgTx <- manager.MakeNewElevator()

			}

			//for i := 0; i < len(p.New); i++ {
			// 	reload orders

		}

		time.Sleep(time.Duration(inputPollRateMs) * time.Millisecond)
	}

}
