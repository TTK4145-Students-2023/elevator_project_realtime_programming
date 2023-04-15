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

	cabsChannelTx := make(chan manager.OrderStruct)
	cabsChannelRx := make(chan manager.OrderStruct)

	stateUpdateTx := make(chan singleElevator.ElevatorUpdateToDatabase)
	stateUpdateRx := make(chan singleElevator.ElevatorUpdateToDatabase)

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
		singleElevator.Fsm_onInitBetweenFloors()
	}

	go singleElevator.SendElevatorToDatabase(stateUpdateTx)

	for {

		select {
		case floor := <-drv_floors:
			database = elevatorInterface.HandleNewFloorAndUpdateDatabase(floor, database, doorTimer, immobilityTimer)

		case button := <-drv_buttons:
			chosenElevator := manager.AssignOrderToElevator(database, button)
			newElevatorUpdate := singleElevator.HandleNewOrder(chosenElevator, button, doorTimer, immobilityTimer)
			database = manager.UpdateDatabase(newElevatorUpdate, database)

		case obstruction := <-drv_obstr:
			if singleElevator.IsDoorOpen() && obstruction {
				doorTimer.Stop()
				immobilityTimer.Reset(3 * time.Second)
				fmt.Println("Nå har jeg resetet immobilityTimer i obstruction")
			} else if !obstruction && singleElevator.IsDoorOpen() {
				immobilityTimer.Stop()
				fmt.Println("Stoppet immobilityTimer i obstruction")
				singleElevator.SetWorkingState(singleElevator.WS_Connected)
				doorTimer.Reset(3 * time.Second)
			}
		case <-drv_stop:
			fmt.Println("Jeg er heis, ", singleElevator.MyID, "her er min heis: ")
			singleElevator.ElevatorPrint(singleElevator.GetSingleEleavtorObject())
			fmt.Println("..og her er databasen min: ")
			for i := 0; i < len(database.ElevatorList); i++ {
				singleElevator.ElevatorPrint(database.ElevatorList[i])
			}

		case <-doorTimer.C:
			singleElevator.Fsm_onDoorTimeout(doorTimer)

		case <-immobilityTimer.C:
			fmt.Println("Iam immobile", singleElevator.MyID)
			database = manager.UpdateElevatorNetworkStateInDatabase(singleElevator.MyID, database, singleElevator.WS_Immobile)

			var deadOrders []elevio.ButtonEvent
			deadOrders = manager.FindDeadOrders(database, singleElevator.MyID)
			for j := 0; j < len(deadOrders); j++ {
				chosenElevator := manager.AssignOrderToElevator(database, deadOrders[j])
				newElevatorUpdate := singleElevator.HandleNewOrder(chosenElevator, deadOrders[j], doorTimer, immobilityTimer)
				database = manager.UpdateDatabase(newElevatorUpdate, database)
			}

		case stateUpdateMessage := <-stateUpdateRx:
			if stateUpdateMessage.ElevatorID != singleElevator.MyID {
				database = manager.UpdateDatabase(stateUpdateMessage.Elevator, database)

				newChangedOrders := manager.SearchMessageForOrderUpdate(stateUpdateMessage, database)

				for i := 0; i < len(newChangedOrders); i++ {
					newOrder := newChangedOrders[i]
					var newElevatorUpdate singleElevator.Elevator

					if newOrder.PanelPair.OrderState == singleElevator.SO_Confirmed {
						chosenElevator := newOrder.PanelPair.ElevatorID
						newButton := newOrder.OrderedButton

						newElevatorUpdate = singleElevator.HandleConfirmedOrder(chosenElevator, newButton, doorTimer, immobilityTimer)

					} else if newOrder.PanelPair.OrderState == singleElevator.SO_NoOrder {
						fmt.Println("Inne no order ifen")
						newElevatorUpdate = singleElevator.Requests_clearOnFloor(newOrder.PanelPair.ElevatorID, newOrder.OrderedButton.Floor)
					}

					database = manager.UpdateDatabase(newElevatorUpdate, database)
				}

			}

		case newCabs := <-cabsChannelRx:
			var newElevatorUpdate singleElevator.Elevator
			fmt.Println("I got a message update from cabs")
			if newCabs.PanelPair.ElevatorID == singleElevator.MyID {
				fmt.Println("I got a cab update and it is for me so I will go through the orders")
				newElevatorUpdate = singleElevator.Fsm_onRequestButtonPress(newCabs.OrderedButton.Floor, newCabs.OrderedButton.Button, singleElevator.MyID, doorTimer, immobilityTimer)
			} else {
				fmt.Println("I got a cab updtae but it is not for me so i will just break")
				break
			}
			database = manager.UpdateDatabase(newElevatorUpdate, database)

		case p := <-peerUpdateCh:

			if len(p.Lost) != 0 {
				var deadOrders []elevio.ButtonEvent
				for i := 0; i < len(p.Lost); i++ {
					database = manager.UpdateElevatorNetworkStateInDatabase(p.Lost[i], database, singleElevator.WS_Unconnected)
					if database.ConnectedElevators <= 1 {
						singleElevator.SetIsAlone(true)
					}
					singleElevator.ElevatorPrint(manager.GetElevatorFromID(database, p.Lost[i]))
					deadOrders = manager.FindDeadOrders(database, p.Lost[i])
					singleElevator.ElevatorPrint(manager.GetElevatorFromID(database, p.Lost[i]))

				}

				for j := 0; j < len(deadOrders); j++ {
					chosenElevator := manager.AssignOrderToElevator(database, deadOrders[j])
					newElevatorUpdate := singleElevator.HandleNewOrder(chosenElevator, deadOrders[j], doorTimer, immobilityTimer)
					database = manager.UpdateDatabase(newElevatorUpdate, database)
				}

			}

			if p.New != "" {
				if !singleElevator.GetIsAlone() {
					cabsToBeSent := manager.FindCabCallsForElevator(database, p.New)
					fmt.Println("Ready to send the following CABs:", cabsToBeSent)
					for k := 0; k < len(cabsToBeSent); k++ {
						cabsChannelTx <- cabsToBeSent[k]
						time.Sleep(time.Duration(inputPollRateMs) * time.Millisecond)
					}
				}

				if !manager.IsElevatorInDatabase(p.New, database) {
					database.ElevatorList = append(database.ElevatorList, singleElevator.Elevator{ElevatorID: p.New, Operating: singleElevator.WS_Connected})
				}

				database = manager.UpdateElevatorNetworkStateInDatabase(p.New, database, singleElevator.WS_Connected)
				if database.ConnectedElevators > 1 {
					singleElevator.SetIsAlone(false)
				}

			}

		}

		time.Sleep(time.Duration(inputPollRateMs) * time.Millisecond)
	}

}
