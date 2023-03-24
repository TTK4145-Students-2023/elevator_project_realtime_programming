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

	cabsChannelTx := make(chan elevator.OrderMessageStruct)
	cabsChannelRx := make(chan elevator.OrderMessageStruct)

	floorArrivalTx := make(chan elevator.FloorArrivalMessageStruct)
	floorArrivalRx := make(chan elevator.FloorArrivalMessageStruct)

	aliveTx := make(chan elevator.IAmAliveMessageStruct)
	aliveRx := make(chan elevator.IAmAliveMessageStruct)

	//ackTx := make(chan manager.AckMessage) Disse kanalene for å sende acks
	//ackRx := make(chan manager.AckMessage)

	go bcast.Transmitter(16569, cabsChannelTx, aliveTx, floorArrivalTx)
	go bcast.Receiver(16569, cabsChannelRx, aliveRx, floorArrivalRx)

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

	//go manager.AliveMessageReceiver(aliveRx, database, newOrder, confirmedOrderChan)

	//input := elevioGetInputDevice()

	if elevio.GetFloor() == -1 {
		elevator.Fsm_onInitBetweenFloors()
	}

	go elevator.SendIAmAlive(aliveTx)

	//prev := [nFloors][nButtons]int{}

	//lage en for select hvor drv_button sender istedet for request buttons
	for {

		select {
		case floor := <-drv_floors:
			//immobilityTimer.Stop()
			//fmt.Println("Stoppet immobility timer on floorArrival")
			//elevator.SetWorkingState(elevator.WS_Connected)
			var newUpdate elevator.Elevator
			newUpdate = elevator.Fsm_onFloorArrival(floor, doorTimer, immobilityTimer)
			database = manager.UpdateDatabase(newUpdate, database)

		case button := <-drv_buttons:
			var newUpdate elevator.Elevator
			chosenElevator := manager.AssignOrderToElevator(database, button)

			if chosenElevator == elevator.MyID { // || orderBroadcast.OrderedButton.Button != elevio.BT_Cab {
				newUpdate = elevator.Fsm_onRequestButtonPress(button.Floor, button.Button, chosenElevator, doorTimer, immobilityTimer) //En toer vil bli satt
			} else {
				newUpdate = elevator.Fsm_setLocalNewOrder(button, chosenElevator)
			}

			database = manager.UpdateDatabase(newUpdate, database)

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
			database = manager.UpdateElevatorNetworkStateInDatabase(elevator.MyID, database, elevator.WS_Immobile)
			manager.ReassignDeadOrders(database, elevator.MyID)
			fmt.Println("Iam immobile", elevator.MyID)

		case aliveMessage := <-aliveRx:
			if aliveMessage.ElevatorID != elevator.MyID {
				database = manager.UpdateDatabase(aliveMessage.Elevator, database)

				newChangedOrders := manager.SearchMessageOrderUpdate(aliveMessage, database)

				for i := 0; i < len(newChangedOrders); i++ {
					newOrder := newChangedOrders[i]
					var newUpdate elevator.Elevator

					if newOrder.PanelPair.OrderState == elevator.SO_Confirmed {
						fmt.Println("Jeg har nå fått bekreftet denne orderen: Floor ", newOrder.OrderedButton.Floor, " - Button ",
							newOrder.OrderedButton.Button, " - chosenElevator ", newOrder.PanelPair.ElevatorID)
						if newOrder.PanelPair.ElevatorID == elevator.MyID {
							fmt.Println("...og den skulle jeg ta selv, så nå kjører jeg Fsm_onReqButPress")
							newUpdate = elevator.Fsm_onRequestButtonPress(newOrder.OrderedButton.Floor, newOrder.OrderedButton.Button, newOrder.PanelPair.ElevatorID, doorTimer, immobilityTimer)
						} else {
							fmt.Println("...men den skulle tas av noen andre, så jeg setter bare en localConfirmedOrder.")
							newUpdate = elevator.Fsm_setLocalConfirmedOrder(newOrder.OrderedButton, newOrder.PanelPair.ElevatorID)
						}
					} else if newOrder.PanelPair.OrderState == elevator.SO_NoOrder {
						fmt.Println("Jeg har nå fått beskjed om å cleare denne orderen: Floor ", newOrder.OrderedButton.Floor, " - Button ",
							newOrder.OrderedButton.Button, " - chosenElevator ", newOrder.PanelPair.ElevatorID)
						newUpdate = elevator.Requests_clearOnFloor(newOrder.PanelPair.ElevatorID, newOrder.OrderedButton.Floor)
					}

					database = manager.UpdateDatabase(newUpdate, database)
				}
				//database = manager.UpdateDatabase(elevatorFromSearch, database)
				//elevator.Fsm_updateLocalRequests(elevatorFromSearch)
			}

		case newCabs := <-cabsChannelRx:
			var newUpdate elevator.Elevator
			fmt.Println("I got a message update from cabs")
			if newCabs.PanelPair.ElevatorID == elevator.MyID {
				fmt.Println("I got a cab update and it is for me so I will go through the orders")
				newUpdate = elevator.Fsm_onRequestButtonPress(newCabs.OrderedButton.Floor, newCabs.OrderedButton.Button, elevator.MyID, doorTimer, immobilityTimer)
			} else {
				fmt.Println("I got a cab updtae but it is not for me so i will just break")
				break
			}
			database = manager.UpdateDatabase(newUpdate, database)

		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

			//legg dette inn i updatenetwork state
			if len(p.Lost) != 0 {
				var ordersToBeReassigned []elevio.ButtonEvent
				for i := 0; i < len(p.Lost); i++ {
					fmt.Println("Lost elevator requests: ", manager.GetElevatorFromID(database, p.Lost[i]).Requests)
					manager.UpdateElevatorNetworkStateInDatabase(p.Lost[i], database, elevator.WS_Unconnected)
					fmt.Println("This is the database: ")
					for i := 0; i < len(database.ElevatorsInNetwork); i++ {
						elevator.ElevatorPrint(database.ElevatorsInNetwork[i])
					}
					ordersToBeReassigned = manager.ReassignDeadOrders(database, p.Lost[i])
					database.NumElevators--

					if database.NumElevators <= 1 {
						elevator.SetIAmAlone(true)
						fmt.Println("I am alone", elevator.GetIAmAlone())
					}
				}
				for j := 0; j < len(ordersToBeReassigned); j++ {
					//newUpdate := elevator.Requests_clearOnFloor(p.Lost[i], ordersToBeReassigned[j].Floor) //OBS! Denne fjerner vel lys i etasjen, men kanskje det settes umerkbart fort igjen?
					//database = manager.UpdateDatabase(newUpdate, database)

					button := ordersToBeReassigned[j] //prøvde her å skrive direkte til drv_buttons, men da oppstod det en lock
					var newUpdate elevator.Elevator
					chosenElevator := manager.AssignOrderToElevator(database, button)
					fmt.Println("The chosen elevator for the reassignment is: ")
					if chosenElevator == elevator.MyID { // || orderBroadcast.OrderedButton.Button != elevio.BT_Cab {
						newUpdate = elevator.Fsm_onRequestButtonPress(button.Floor, button.Button, chosenElevator, doorTimer, immobilityTimer) //En toer vil bli satt
					} else {
						newUpdate = elevator.Fsm_setLocalNewOrder(button, chosenElevator)
					}

					database = manager.UpdateDatabase(newUpdate, database)
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
					cabsToBeSent := manager.SendCabCallsForElevator(database, p.New)
					fmt.Println("Ready to send the following CABs:", cabsToBeSent)
					for k := 0; k < len(cabsToBeSent); k++ {
						cabsChannelTx <- cabsToBeSent[k]
					}
					//OBS! Kanskje vi må lage en egen kanal for dette?
				}

				database.NumElevators++
				if database.NumElevators > 1 {
					elevator.SetIAmAlone(false)
				}
			}

			//for i := 0; i < len(p.New); i++ {
			// 	reload orders

		}

		time.Sleep(time.Duration(inputPollRateMs) * time.Millisecond)
	}

}
