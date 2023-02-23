package main

import (
	"Driver-go/elevator"
	"Driver-go/elevio"
	"fmt"
	"time"
)

const nFloors = 4
const nButtons = 3

//elev.Timer = *time.NewTimer(time.Second * 1) hvordam lage timer

func main() {
	fmt.Println("Started!")

	inputPollRateMs := 25
	//var orders [4]int

	elevio.Init("localhost:15657", nFloors)

	//var d elevio.MotorDirection = elevio.MD_Up

	//elevio.SetMotorDirection(d)

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	var elevPtr = elevator.Elevator_uninitialized()

	if elevio.GetFloor() == -1 {
		elevator.Fsm_onInitBetweenFloors(elevPtr)
		fmt.Printf("INNE I INIT BETWEEN")
	} else {
		elevPtr = elevator.Fsm_init(elevPtr)
	}

	//prev := [nFloors][nButtons]int{}

	//elev.Timer := time.NewTimer(time.Second*1)

	//lage en for select hvor drv_button sender istedet for request buttons
	for {

		select {
		case floor := <-drv_floors:
			fmt.Printf("%+v\n", floor)
			elevator.Fsm_onFloorArrival(floor, elevPtr)

			fmt.Printf(" ---------case floor----------")

		case button := <-drv_buttons:
			elevator.Fsm_onRequestButtonPress(button.Floor, button.Button, elevPtr)
			fmt.Printf("inne i poll buttons")
		/*
			case timer := <-drv_timer:
				fmt.Print(timer)
				elevator.Timer_stop()
				elevator.Fsm_onDoorTimeout()
		*/
		case obstruction := <-drv_obstr:
			if obstruction {
				elevio.SetMotorDirection(elevio.MD_Stop)
			}
			fmt.Printf(" ---------case obstruction----------")

		case <-elevPtr.Timer.C:
			fmt.Printf("-----før timerout if---------")

			elevPtr.Timer.Stop()
			elevator.Fsm_onDoorTimeout(elevPtr)

			fmt.Printf("----etter if timedout")

		}
		time.Sleep(time.Duration(inputPollRateMs) * time.Millisecond)
	}

}
