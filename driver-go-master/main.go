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
const nButtons = 3

type HelloMsg struct {
	Message string
	Iter    int
}

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
	go peers.Transmitter(15647, id, peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)

	// We make channels for sending and receiving our custom data types
	//helloTx := make(chan HelloMsg)
	//helloRx := make(chan HelloMsg)

	orderTx := make(chan elevator.OrderMessageStruct)
	orderRx := make(chan elevator.OrderMessageStruct)

	floorArrivalTx := make(chan elevator.FloorArrivalMessageStruct)
	floorArrivalRx := make(chan elevator.FloorArrivalMessageStruct)

	aliveTx := make(chan elevator.IAmAliveMessageStruct)
	aliveRx := make(chan elevator.IAmAliveMessageStruct)
	// ... and start the transmitter/receiver pair on some port
	// These functions can take any number of channels! It is also possible to
	//  start multiple transmitters/receivers on the same port.
	go bcast.Transmitter(16569, orderTx, aliveTx, floorArrivalTx)
	go bcast.Receiver(16569, orderRx, aliveRx, floorArrivalRx)

	go elevator.SendIAmAlive(aliveTx)
	//port: 16569

	// The example message. We just send one of these every second.
	/*o func() {
		helloMsg := HelloMsg{"Dette er den andre pc'en " + id, 0}
		for {
			helloMsg.Iter++
			//helloTx <- helloMsg
			time.Sleep(1 * time.Second)
		}
	}()*/

	/*fmt.Println("Started")
	for {
		select {
		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

		case a := <-helloRx:
			fmt.Printf("Received: %#v\n", a)
		}
	}*/

	fmt.Println("Started!")

	database := manager.ElevatorDatabase{
		//hardkodede verdier vi alltid bruker når vi flagger
		NumElevators: 3,

		Elevator13520: elevator.Elevator_uninitialized("13520"),
		Elevator70310: elevator.Elevator_uninitialized("70310"),
		Elevator54321: elevator.Elevator_uninitialized("54321"),
	}
	

	inputPollRateMs := 25
	//var orders [4]int

	elevio.Init("localhost: 15657", nFloors)

	//var d elevio.MotorDirection = elevio.MD_Up

	//elevio.SetMotorDirection(d)

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)
	drv_timer := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)
	go elevator.Timer_runTimer(drv_timer)

	//input := elevioGetInputDevice()

	if elevio.GetFloor() == -1 {
		elevator.Fsm_onInitBetweenFloors()
	}

	//prev := [nFloors][nButtons]int{}

	//lage en for select hvor drv_button sender istedet for request buttons
	for {

		select {
		case floor := <-drv_floors:
			floorMsg := elevator.FloorArrivalMessageStruct{SystemID: "Gruppe10",
				MessageID:    "Order",
				ElevatorID:   elevator.MyID,
				ArrivedFloor: floor}

			floorArrivalTx <- floorMsg
			fmt.Printf("%+v\n", floor)
			//elevator.Fsm_onFloorArrival(floor)
			//send melding om ankomst til floor

		case floorArrivalBroadcast := <-floorArrivalRx:
			if floorArrivalBroadcast.ElevatorID == elevator.MyID {
				elevator.Fsm_onFloorArrival(floorArrivalBroadcast.ArrivedFloor)
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
			
			//Husk at vi skal fikse CAB som en egen greie
			//pakk inn i melding og send
			orderMsg := elevator.OrderMessageStruct{SystemID: "Gruppe10",
				MessageID:      "Order",
				ElevatorID:     elevator.MyID,
				OrderedButton:  button,
				ChosenElevator: chosenElevator}

			orderTx <- orderMsg

			//elevator.Fsm_onRequestButtonPress(button.Floor, button.Button) //droppe denne

		case timer := <-drv_timer:
			fmt.Print(timer)

			elevator.Fsm_onDoorTimeout()

		case obstruction := <-drv_obstr:
			if obstruction {
				elevio.SetMotorDirection(elevio.MD_Stop)
			}

		case orderBroadcast := <-orderRx:
			fmt.Printf("Received: %#v\n", orderBroadcast)
			
			if (orderBroadcast.OrderedButton.Button == elevio.BT_Cab && orderBroadcast.ChosenElevator == elevator.MyID) ||
				orderBroadcast.OrderedButton.Button != elevio.BT_Cab {
				elevator.Fsm_onRequestButtonPress(orderBroadcast.OrderedButton.Floor, orderBroadcast.OrderedButton.Button, orderBroadcast.ChosenElevator)
			}

			fmt.Printf("Received database: %#v\n", database)
			

		case aliveMsg := <-aliveRx:
			//oppdater tilhørende heis i databasestruct (dette er for å regne cost)
			switch aliveMsg.ElevatorID {
			case "13520":
				if aliveMsg.Elevator.Operating != elevator.WS_NoMotor {
					aliveMsg.Elevator.Operating = elevator.WS_Running
				}
				database.Elevator13520 = aliveMsg.Elevator
			case "70310":
				if aliveMsg.Elevator.Operating != elevator.WS_NoMotor {
					aliveMsg.Elevator.Operating = elevator.WS_Running
				}
				database.Elevator70310 = aliveMsg.Elevator
			case "54321":
				if aliveMsg.Elevator.Operating != elevator.WS_NoMotor {
					aliveMsg.Elevator.Operating = elevator.WS_Running
				}
				database.Elevator54321 = aliveMsg.Elevator
			}

			//fmt.Printf("Recivied IAMALIVE:  %#v\n", aliveMsg)
		}

		//case: mottatt broadcast-ordre
		//putt i array (for å stoppe ved onFloorArrival)
		//Hvis mottatt ordre har min elevatorID:
		//Fsm_onReq

		time.Sleep(time.Duration(inputPollRateMs) * time.Millisecond)
	}

}
