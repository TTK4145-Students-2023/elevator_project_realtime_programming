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

	orderTx := make(chan elevator.OrderMessageStruct)
	orderRx := make(chan elevator.OrderMessageStruct)

	floorArrivalTx := make(chan elevator.FloorArrivalMessageStruct)
	floorArrivalRx := make(chan elevator.FloorArrivalMessageStruct)

	aliveTx := make(chan elevator.IAmAliveMessageStruct)
	aliveRx := make(chan elevator.IAmAliveMessageStruct)

	newOrder := make(chan elevator.Elevator)
	confirmedOrderChan := make(chan elevator.OrderMessageStruct)

	//ackTx := make(chan manager.AckMessage) Disse kanalene for å sende acks
	//ackRx := make(chan manager.AckMessage)

	go bcast.Transmitter(16569, orderTx, aliveTx, floorArrivalTx)
	go bcast.Receiver(16569, orderRx, aliveRx, floorArrivalRx)

	go elevator.SendIAmAlive(aliveTx)
	//port: 16569

	fmt.Println("Started!")

	database := manager.ElevatorDatabase{
		//hardkodede verdier vi alltid bruker når vi flagger
		NumElevators: 0,
	}

	doorTimer := time.NewTimer(3 * time.Second)
	doorTimer.Stop()

	immobilityTimer := time.NewTimer(3 * time.Second)
	immobilityTimer.Stop()

	inputPollRateMs := 25

	elevio.Init("localhost:"+id, nFloors) //endre denne for å bruke flere sockets for elevcd //15657

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	go manager.AliveMessageReceiver(aliveRx, database, newOrder, confirmedOrderChan)

	//input := elevioGetInputDevice()

	if elevio.GetFloor() == -1 {
		elevator.Fsm_onInitBetweenFloors()
	}

	//prev := [nFloors][nButtons]int{}

	//lage en for select hvor drv_button sender istedet for request buttons
	for {

		select {
		case floor := <-drv_floors:
			//immobilityTimer.Stop()
			//fmt.Println("Stoppet immobility timer on floorArrival")
			//elevator.SetWorkingState(elevator.WS_Connected)
			elevator.Fsm_onFloorArrival(floor, doorTimer, immobilityTimer)

		case button := <-drv_buttons:
			chosenElevator := manager.AssignOrderToElevator(database, button)

			if chosenElevator == elevator.MyID { // || orderBroadcast.OrderedButton.Button != elevio.BT_Cab {
				elevator.Fsm_onRequestButtonPress(button.Floor, button.Button, chosenElevator, doorTimer, immobilityTimer) //En toer vil bli satt
			} else {
				elevator.Fsm_localNewOrder(button, chosenElevator)
			}

		case obstruction := <-drv_obstr:
			if elevator.IsDoorOpen() && obstruction {
				doorTimer.Stop()
				immobilityTimer.Reset(3 * time.Second)
				fmt.Println("Nå har jeg resetet immobilityTimer i obstruction")
			} else if !obstruction && elevator.IsDoorOpen() {
				immobilityTimer.Stop()
				fmt.Println("Stoppet immobilityTimer i obstruction")
				elevator.SetWorkingState(elevator.WS_Connected)
				doorTimer.Reset(3 * time.Second)
			}

		case <-doorTimer.C:

			elevator.Fsm_onDoorTimeout(doorTimer)

		case <-immobilityTimer.C:
			//skal redestribute ordre
			elevator.SetWorkingState(elevator.WS_Immobile)
			manager.UpdateElevatorNetworkStateInDatabase(elevator.MyID, database, elevator.WS_Immobile)
			manager.ReassignDeadOrders(orderTx, database, elevator.MyID)
			fmt.Println("Iam immobile", elevator.MyID)

		case order := <-confirmedOrderChan:
			elevator.Fsm_onRequestButtonPress(order.OrderedButton.Floor, order.OrderedButton.Button, order.ChosenElevator, doorTimer, immobilityTimer)

		case update := <-newOrder:
			elevator.Fsm_updateLocalRequests(update)
			database = manager.UpdateDatabase(update, database)

		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

			//legg dette inn i updatenetwork state
			if len(p.Lost) != 0 {
				for i := 0; i < len(p.Lost); i++ {
					manager.UpdateElevatorNetworkStateInDatabase(p.Lost[i], database, elevator.WS_Unconnected)
					fmt.Println("Lost elevator requests: ", manager.GetElevatorFromID(database, p.Lost[i]).Requests)
					fmt.Println("This is the database: ")
					for i := 0; i < len(database.ElevatorsInNetwork); i++ {
						elevator.ElevatorPrint(database.ElevatorsInNetwork[i])
					}
					manager.ReassignDeadOrders(orderTx, database, p.Lost[i])
					database.NumElevators--
				}
				if database.NumElevators <= 1 {
					elevator.SetIAmAlone(true)
					fmt.Println("I am alone", elevator.GetIAmAlone())
				}
			}

			if p.New != "" {
				if !manager.IsElevatorInDatabase(p.New, database) {
					database.ElevatorsInNetwork = append(database.ElevatorsInNetwork, elevator.Elevator{ElevatorID: p.New, Operating: elevator.WS_Connected})
				}

				for i := 0; i < len(database.ElevatorsInNetwork); i++ {
					if database.ElevatorsInNetwork[i].ElevatorID == p.New {
						database.ElevatorsInNetwork[i].Operating = elevator.WS_Connected
					}
				}

				if !elevator.GetIAmAlone() {
					fmt.Println("Ready to send CABs")
					manager.SendCabCallsForElevator(orderTx, database, p.New)
				}

				database.NumElevators++
				if database.NumElevators > 1 {
					elevator.SetIAmAlone(false)
				}
			}

			//for i := 0; i < len(p.New); i++ {
			// 	reload orders

		}

		//case: mottatt broadcast-ordre
		//putt i array (for å stoppe ved onFloorArrival)
		//Hvis mottatt ordre har min elevatorID:
		//Fsm_onReq

		time.Sleep(time.Duration(inputPollRateMs) * time.Millisecond)
	}

}
