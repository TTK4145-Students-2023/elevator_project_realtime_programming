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

	cabsChannelTx := make(chan elevator.OrderStruct)
	cabsChannelRx := make(chan elevator.OrderStruct)

	stateUpdateTx := make(chan elevator.StateUpdateStruct)
	stateUpdateRx := make(chan elevator.StateUpdateStruct)


	go bcast.Transmitter(16569, cabsChannelTx, stateUpdateTx)
	go bcast.Receiver(16569, cabsChannelRx, stateUpdateRx)
	//port: 16569

	fmt.Println("Started!")

	database := manager.ElevatorDatabase{
		ConnectedElevators: 0,
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


	if elevio.GetFloor() == -1 {
		elevator.Fsm_onInitBetweenFloors()
	}

	go elevator.SendStateUpdate(stateUpdateTx)

	
	for {

		select {
		case floor := <-drv_floors:
			var newElevatorUpdate elevator.Elevator
			newElevatorUpdate = elevator.Fsm_onFloorArrival(floor, doorTimer, immobilityTimer)
			fmt.Println("Her oppdaterer jeg databasen med en slettet ordre")
			database = manager.UpdateDatabase(newElevatorUpdate, database)

		case button := <-drv_buttons:
			chosenElevator := manager.AssignOrderToElevator(database, button)
			newElevatorUpdate := elevator.HandleNewOrder(chosenElevator, button, doorTimer, immobilityTimer)
			database = manager.UpdateDatabase(newElevatorUpdate, database)

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
		case <-drv_stop:
			fmt.Println("Jeg er heis, ", elevator.MyID, "her er min heis: ")
			elevator.ElevatorPrint(elevator.GetSingleEleavtorStruct())
			fmt.Println("..og her er databasen min: ")
			for i := 0; i < len(database.ElevatorList); i++ {
				elevator.ElevatorPrint(database.ElevatorList[i])
			}

		case <-doorTimer.C:
			elevator.Fsm_onDoorTimeout(doorTimer)

		case <-immobilityTimer.C:
			fmt.Println("Iam immobile", elevator.MyID)
			database = manager.UpdateElevatorNetworkStateInDatabase(elevator.MyID, database, elevator.WS_Immobile)

			var deadOrders []elevio.ButtonEvent
			deadOrders = manager.FindDeadOrders(database, elevator.MyID)
			for j := 0; j < len(deadOrders); j++ {
				chosenElevator := manager.AssignOrderToElevator(database, deadOrders[j])
				newElevatorUpdate := elevator.HandleNewOrder(chosenElevator, deadOrders[j], doorTimer, immobilityTimer)
				database = manager.UpdateDatabase(newElevatorUpdate, database)
			}

		case aliveMessage := <-stateUpdateRx:
			if aliveMessage.ElevatorID != elevator.MyID {
				database = manager.UpdateDatabase(aliveMessage.Elevator, database)

				newChangedOrders := manager.SearchMessageForOrderUpdate(aliveMessage, database)

				for i := 0; i < len(newChangedOrders); i++ {
					newOrder := newChangedOrders[i]
					var newElevatorUpdate elevator.Elevator

					if newOrder.PanelPair.OrderState == elevator.SO_Confirmed {
						chosenElevator := newOrder.PanelPair.ElevatorID
						newButton := newOrder.OrderedButton

						newElevatorUpdate = elevator.HandleConfirmedOrder(chosenElevator, newButton, doorTimer, immobilityTimer)

					} else if newOrder.PanelPair.OrderState == elevator.SO_NoOrder {
						fmt.Println("Inne no order ifen")
						newElevatorUpdate = elevator.Requests_clearOnFloor(newOrder.PanelPair.ElevatorID, newOrder.OrderedButton.Floor)
					}

					database = manager.UpdateDatabase(newElevatorUpdate, database)
				}

			}

		case newCabs := <-cabsChannelRx:
			var newElevatorUpdate elevator.Elevator
			fmt.Println("I got a message update from cabs")
			if newCabs.PanelPair.ElevatorID == elevator.MyID {
				fmt.Println("I got a cab update and it is for me so I will go through the orders")
				newElevatorUpdate = elevator.Fsm_onRequestButtonPress(newCabs.OrderedButton.Floor, newCabs.OrderedButton.Button, elevator.MyID, doorTimer, immobilityTimer)
			} else {
				fmt.Println("I got a cab updtae but it is not for me so i will just break")
				break
			}
			database = manager.UpdateDatabase(newElevatorUpdate, database)

		case p := <-peerUpdateCh:
			

			if len(p.Lost) != 0 {
				var deadOrders []elevio.ButtonEvent
				for i := 0; i < len(p.Lost); i++ {
					database = manager.UpdateElevatorNetworkStateInDatabase(p.Lost[i], database, elevator.WS_Unconnected)
					if database.ConnectedElevators <= 1 {
						elevator.SetIAmAlone(true)
					}
					elevator.ElevatorPrint(manager.GetElevatorFromID(database, p.Lost[i]))
					deadOrders = manager.FindDeadOrders(database, p.Lost[i])
					elevator.ElevatorPrint(manager.GetElevatorFromID(database, p.Lost[i]))

				}

				for j := 0; j < len(deadOrders); j++ {
					chosenElevator := manager.AssignOrderToElevator(database, deadOrders[j])
					newElevatorUpdate := elevator.HandleNewOrder(chosenElevator, deadOrders[j], doorTimer, immobilityTimer)
					database = manager.UpdateDatabase(newElevatorUpdate, database)
				}

			}

			if p.New != "" {
				if !elevator.GetIAmAlone() {
					cabsToBeSent := manager.FindCabCallsForElevator(database, p.New)
					fmt.Println("Ready to send the following CABs:", cabsToBeSent)
					for k := 0; k < len(cabsToBeSent); k++ {
						cabsChannelTx <- cabsToBeSent[k]
						time.Sleep(time.Duration(inputPollRateMs) * time.Millisecond)
					}
				}

				if !manager.IsElevatorInDatabase(p.New, database) {
					database.ElevatorList = append(database.ElevatorList, elevator.Elevator{ElevatorID: p.New, Operating: elevator.WS_Connected})
				}

				database = manager.UpdateElevatorNetworkStateInDatabase(p.New, database, elevator.WS_Connected)
				if database.ConnectedElevators > 1 {
					elevator.SetIAmAlone(false)
				}

			}

		}

		time.Sleep(time.Duration(inputPollRateMs) * time.Millisecond)
	}

}
