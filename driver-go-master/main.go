package main

import (
	"Driver-go/elevator"
	"Driver-go/elevio"
	"fmt"
	"time"
)

const nFloors = 4
const nButtons = 3

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
	//drv_timedOut := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)
	//go elevator.Timer_timedOut(drv_timedOut, &elevator.TimerActive)

	//input := elevioGetInputDevice()
	var elev = elevator.Elevator_uninitialized()

	if elevio.GetFloor() == -1 {
		elevator.Fsm_onInitBetweenFloors()
	} else {
		elevator.Fsm_init()
	}

	//prev := [nFloors][nButtons]int{}

	//timer := time.NewTimer(time.Second*3)

	//lage en for select hvor drv_button sender istedet for request buttons
	for {

		select {
		case floor := <-drv_floors:
			fmt.Printf("%+v\n", floor)
			elevator.Fsm_onFloorArrival(floor)

		case button := <-drv_buttons:
			elevator.Fsm_onRequestButtonPress(button.Floor, button.Button)

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

		case <-elev.Timer.C:
			fmt.Printf("-----før timerout if---------")

			elev.Timer.Stop()
			elevator.Fsm_onDoorTimeout()

			fmt.Printf("----etter if timedout")

		}
		time.Sleep(time.Duration(inputPollRateMs) * time.Millisecond)
	}
}
