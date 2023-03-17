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

	timer := time.NewTimer(3 * time.Second)
	timer.Stop()

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

	//input := elevioGetInputDevice()

	if elevio.GetFloor() == -1 {
		elevator.Fsm_onInitBetweenFloors()
	}

	//prev := [nFloors][nButtons]int{}

	//lage en for select hvor drv_button sender istedet for request buttons
	for {

		select {
		case floor := <-drv_floors:
			floorArrivalTx <- elevator.MakeFloorMessage(floor) //ackRx) //Ta inn ack kanalen og håndtere den inn i funksjonen?

			//SendFloorArrival(floorMsg,floorArrivalTx)

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

		case floorArrivalBroadcast := <-floorArrivalRx:
			if floorArrivalBroadcast.ElevatorID == elevator.MyID {
				elevator.Fsm_onFloorArrival(floorArrivalBroadcast.ArrivedFloor, timer)
			} else {
				elevator.Requests_clearOnFloor(floorArrivalBroadcast.ElevatorID, floorArrivalBroadcast.ArrivedFloor)
			}
		//case: mottatt melding om at kommet til floor
		//if msg.ID == MYID: fsm_onFloorArrival
		//else Requests_clearOnFloor

		case button := <-drv_buttons:
			//Heis tilhørende panelet regner ut cost for alle tre heiser
			//Broadcaster fordelt ordre (med elevatorID)
			//Hvis CAB-order: håndter internt (ikke broadcast)
			//CAB-order deles ikke som en ordre, men som del av heis-tilstand/info
			chosenElevator := manager.AssignOrderToElevator(database, button)

			orderTx <- elevator.MakeOrderMessage(chosenElevator, button)

			//manager.SendOrderMessage(orderTx, button, database)

			/*chosenElevator := manager.AssignOrderToElevator(database, button)

			//Husk at vi skal fikse CAB som en egen greie
			//pakk inn i melding og send
			orderMsg := elevator.OrderMessageStruct{SystemID: "Gruppe10",
				MessageID:      "Order",
				ElevatorID:     elevator.MyID,
				OrderedButton:  button,
				ChosenElevator: chosenElevator}

			orderTx <- orderMsg
			*/

			//elevator.Fsm_onRequestButtonPress(button.Floor, button.Button) //droppe denne

		case obstruction := <-drv_obstr:
			if elevator.IsDoorOpen() && obstruction {
				timer.Stop()
			} else if !obstruction && elevator.IsDoorOpen() {
				timer.Reset(3 * time.Second)
			}

		case timedOut := <-timer.C:
			fmt.Println("fått lest fra timer.C")

			fmt.Print(timedOut)

			elevator.Fsm_onDoorTimeout(timer)

			//if obstruction {
			//	elevio.SetMotorDirection(elevio.MD_Stop)
			//}

		case orderBroadcast := <-orderRx:

			fmt.Printf("Received: %#v\n", orderBroadcast)
			if orderBroadcast.OrderedButton.Button != elevio.BT_Cab {
				elevator.Elevator_increaseOrderNumber()
			}

			//if chosenElev already on floor -> Request_clearOnFloor
			if (orderBroadcast.OrderedButton.Button == elevio.BT_Cab && orderBroadcast.ChosenElevator == elevator.MyID) ||
				orderBroadcast.OrderedButton.Button != elevio.BT_Cab {
				elevator.Fsm_onRequestButtonPress(orderBroadcast.OrderedButton.Floor, orderBroadcast.OrderedButton.Button, orderBroadcast.ChosenElevator, timer)
			}

			//HER LA VI TIL EN SJEKK OM CHOSEN ELEVTAOR ER I ETASJEN TIL BESTILLINGEN ALLEREDE, hvis den er det skal bestillingen cleares med en gang.
			//burde sikkert være innbakt et annet sted.
			if manager.WhatFloorIsElevatorFromStringID(database, orderBroadcast.ChosenElevator) == orderBroadcast.OrderedButton.Floor &&
				manager.WhatStateIsElevatorFromStringID(database, orderBroadcast.ChosenElevator) != elevator.EB_Moving {
				elevator.Requests_clearOnFloor(orderBroadcast.ChosenElevator, orderBroadcast.OrderedButton.Floor)
			}

			//fmt.Printf("Received database: %#v\n", database)

		case aliveMsg := <-aliveRx:
			//oppdater tilhørende heis i databasestruct (dette er for å regne cost)

			if !manager.IsElevatorInDatabase(aliveMsg.ElevatorID, database) {
				database.ElevatorsInNetwork = append(database.ElevatorsInNetwork, aliveMsg.Elevator)
				database.NumElevators++
			}

			manager.UpdateDatabase(aliveMsg, database)

		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

			manager.UpdateElevatorNetworkStateInDatabase(p, database)

			//legg dette inn i updatenetwork state
			if len(p.Lost) != 0 {
				for i := 0; i < len(p.Lost); i++ {
					manager.ReassignDeadOrders(orderTx, database, p.Lost[i])
				}
			}

			//if p.New != ""

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
