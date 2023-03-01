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
			fmt.Printf("%+v\n", floor)
			elevator.Fsm_onFloorArrival(floor)

		case button := <-drv_buttons:
			elevator.Fsm_onRequestButtonPress(button.Floor, button.Button)

			
		
		case timer := <-drv_timer:
			fmt.Print(timer)
	
			elevator.Fsm_onDoorTimeout()
		
		case obstruction := <-drv_obstr:
			if obstruction {
				elevio.SetMotorDirection(elevio.MD_Stop)
			}

		}
		time.Sleep(time.Duration(inputPollRateMs) * time.Millisecond)
	}
}
