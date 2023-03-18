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

	orderTx := make(chan manager.OrderMessageStruct)
	orderRx := make(chan manager.OrderMessageStruct)

	floorArrivalTx := make(chan manager.FloorArrivalMessageStruct)
	floorArrivalRx := make(chan manager.FloorArrivalMessageStruct)

	newElevatorTx := make(chan manager.SingleElevatorMessageStruct)
	newElevatorRx := make(chan manager.SingleElevatorMessageStruct)

	//msgTx := make(chan manager.MessageStruct)
	//msgRx := make(chan manager.MessageStruct)

	ackTx := make(chan manager.AckMessageStruct) //Disse kanalene for å sende acks
	ackRx := make(chan manager.AckMessageStruct)

	go bcast.Transmitter(11771, orderTx, floorArrivalTx, newElevatorTx, ackTx)
	go bcast.Receiver(11771, orderRx, floorArrivalRx, newElevatorRx, ackRx)
	//port: 16569

	fmt.Println("Started!")

	mainTimer := time.NewTimer(3 * time.Second)
	mainTimer.Stop()

	mainDatabase := manager.ElevatorDatabase{
		//hardkodede verdier vi alltid bruker når vi flagger
		NumElevators: 0,
	}

	inputPollRateMs := 25

	elevio.Init("localhost:"+id, nFloors) //endre denne for å bruke flere sockets for elevcd //15657

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	databaseChannel := make(chan manager.ElevatorDatabase)
	initiateResend := make(chan manager.OrderMessageStruct)
	shouldAcknowledge := make(chan manager.AckMessageStruct)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	go manager.ReceiveMessages(*mainTimer, mainDatabase, databaseChannel, ackTx, orderRx, floorArrivalRx, ackRx, newElevatorRx, shouldAcknowledge)

	if elevio.GetFloor() == -1 {
		elevator.Fsm_onInitBetweenFloors()
	}

	for {

		select {
		case floor := <-drv_floors:
			floorMessage := manager.MakeFloorMessage(floor)
			/*if !manager.IsElevatorInDatabase(floorMessage.SenderID, mainDatabase) {
				floorMessage.MyElevator.Floor = floor
				mainDatabase.ElevatorsInNetwork = append(mainDatabase.ElevatorsInNetwork, floorMessage.MyElevator)
				mainDatabase.NumElevators++
				fmt.Println("The elevator has been put into the database by a floor sensor.\nDatabase: ", mainDatabase)
			}*/
			floorArrivalTx <- floorMessage //ackRx) //Ta inn ack kanalen og håndtere den inn i funksjonen?
			fmt.Println("sent floor message")
			manager.FloorArrivalBroadcast_timeAcknowledgementAndResend(floorArrivalTx, floorMessage, ackRx, mainDatabase.NumElevators)

		case button := <-drv_buttons:
			chosenElevator := manager.AssignOrderToElevator(mainDatabase, button)
			orderMessage := manager.MakeOrderMessage(chosenElevator, button)

			orderTx <- orderMessage
			fmt.Println("sent button message!")
			fmt.Println("The order was from floor: ", orderMessage.OrderedButton.Floor, "\nAnd the order was assigned to: ", orderMessage.ChosenElevator)

			manager.OrderBroadcast_timeAcknowledgementAndResend(initiateResend, orderMessage, ackRx, mainDatabase.NumElevators)

			//elevator.Fsm_onRequestButtonPress(button.Floor, button.Button) //droppe denne

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

			//fmt.Printf("Received database: %#v\n", database

		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

			manager.UpdateElevatorNetworkStateInDatabase(p, mainDatabase)

			//legg dette inn i updatenetwork state
			if len(p.Lost) != 0 {
				for i := 0; i < len(p.Lost); i++ {
					manager.ReassignDeadOrders(orderTx, mainDatabase, p.Lost[i])
				}
			}

			if p.New != "" {
				newElevatorTx <- manager.MakeNewElevator()

			}

			//for i := 0; i < len(p.New); i++ {
			// 	reload orders

		case updatedDatabase := <-databaseChannel:
			mainDatabase = updatedDatabase

		case shouldResend := <-initiateResend:
			orderTx <- shouldResend
		case ack := <-shouldAcknowledge:
			ackTx <- ack
		}

		time.Sleep(time.Duration(inputPollRateMs) * time.Millisecond)
	}

}

//---------------15.03-----------------
/*
	DETTE SKAL SKJE VED HVER HARDWARE EVENT: FLOORSENSOR, BUTTONS, OBSTRUCTION
	1. sende standardisert melding med "fysisk" heis (struct), avsenderID, meldingsID og info for gjeldende event
	2. starte venting på acks som har din egen ID (sette en readDeadline ellerno)

	Om å løse samme problem to ganger: Om vi resender ti ganger og fortsatt ikke har mottatt acks fra begge heiser,
	så BØR det ha skjedd noe i peers. Sannsyligheten for samme packet loss ti ganger, er veldig liten.

	JEG HAR MOTTATT MELDING
	1. Send acknowledge med avsenderIDen som er med i den mottatte meldingen (Inkludere meldingscounter)
	2. Switch case på meldingsID:
		2a. meldingsID: FLOORSENSOR -> Fsm_onFloorArrival() (OBS! Logikk vi har)
		2b. osv.....

	ACKNOWLEDGEMENTKANAL
	1. Mottatt ack. -> ack-count++
	2. Hvis to tre acks, så er alt good
	3. Hvis timer har gått ut og mottatt under tre acks, resend

	MESSAGECOUNTERS
	- Hver heis inkrementerer srivalTin messagesSent når de sender en melding. Denne counten er inkludert i meldingen som sendes.
	- Hver individuelle count lagres til tilhørende heis i databasen og sjekkes opp imot den mottatte meldingscounten.
	- Ikke ack om du mottar en meldingscount som er >1 over din lagrede info:
		- Vent på resending av alle manglende meldinger


		------- Litt søppel -------ø
		2a. ikke mottatt ack: resend (x5 feks vet ikke helt) (OBS! Duplicate messages!! Counters?)
		2b. mottatt ack -> bra
	3. har ikke mottatt noen acks etter x antall resends -> jeg er alene og død
		(OBS! Løser dette samme problem to ganger ettersom at vi har peers?)
	4. Kun mottatt ack fra 1 heis (selv etter flere resends) -> Den andre heisen er død
		4a. Initiere en sjekk som gjør at begge heiser erklærer den død?? (PEERS???)

*/
/*case MessageBroadcast := <-msgRx:
fmt.Println("I have received a message! It should go to case: ", MessageBroadcast.MessageType)
switch MessageBroadcast.MessageType {
case 17:
	fmt.Println("received floor")
	manager.SendAcknowledge(ackTx, MessageBroadcast)
	//if !manager.IsElevatorInDatabase(MessageBroadcast.SenderID, database) {
		database.ElevatorsInNetwork = append(database.ElevatorsInNetwork, MessageBroadcast.MyElevator)
		database.NumElevators++
	}
	manager.UpdateDatabase(MessageBroadcast, database)
	//------------------------ SLUTT: DETTE ER FRA IAMALIVE-CASE ------------------------

	if MessageBroadcast.SenderID == elevator.MyID {
		elevator.Fsm_onFloorArrival(MessageBroadcast.ArrivedFloor, mainTimer)
	} else {
		elevator.Requests_clearOnFloor(MessageBroadcast.SenderID, MessageBroadcast.ArrivedFloor)
	}

case 69:
	fmt.Println("received order")
	manager.SendAcknowledge(ackTx, MessageBroadcast)
	//if !manager.IsElevatorInDatabase(MessageBroadcast.SenderID, database) {
		database.ElevatorsInNetwork = append(database.ElevatorsInNetwork, MessageBroadcast.MyElevator)
		database.NumElevators++
		fmt.Println(" number of elevatrs ---  ")
		fmt.Println(database.NumElevators)
	}
	manager.UpdateDatabase(MessageBroadcast, database)
	//------------------------ SLUTT: DETTE ER FRA IAMALIVE-CASE ------------------------

	//fmt.Printf("Received: %#v\n", MessageBroadcast)
	if MessageBroadcast.OrderedButton.Button != elevio.BT_Cab {
		elevator.Elevator_increaseOrderNumber()
	}

	//if chosenElev already on floor -> Request_clearOnFloor
	if (MessageBroadcast.OrderedButton.Button == elevio.BT_Cab && MessageBroadcast.ChosenElevator == elevator.MyID) ||
		MessageBroadcast.OrderedButton.Button != elevio.BT_Cab {
		elevator.Fsm_onRequestButtonPress(MessageBroadcast.OrderedButton.Floor, MessageBroadcast.OrderedButton.Button, MessageBroadcast.ChosenElevator, mainTimer)
	}

	//HER LA VI TIL EN SJEKK OM CHOSEN ELEVTAOR ER I ETASJEN TIL BESTILLINGEN ALLEREDE, hvis den er det skal bestillingen cleares med en gang.
	//burde sikkert være innbakt et annet sted.
	if manager.WhatFloorIsElevatorFromStringID(database, MessageBroadcast.ChosenElevator) == MessageBroadcast.OrderedButton.Floor &&
		manager.WhatStateIsElevatorFromStringID(database, MessageBroadcast.ChosenElevator) != elevator.EB_Moving {
		elevator.Requests_clearOnFloor(MessageBroadcast.ChosenElevator, MessageBroadcast.OrderedButton.Floor)
	}
case 666:
	if !manager.IsElevatorInDatabase(MessageBroadcast.SenderID, database) {
		database.ElevatorsInNetwork = append(database.ElevatorsInNetwork, MessageBroadcast.MyElevator)
		database.NumElevators++
		fmt.Println(" number of elevatrs ---  ")
		fmt.Println(database.NumElevators)
	}
default:
	fmt.Println("...but I don't know what to do with it.")
}*/
