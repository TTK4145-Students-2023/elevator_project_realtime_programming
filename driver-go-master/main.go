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

	var elev = elevator.Elevator_uninitialized()

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)
	drv_timeOut := make(chan bool)

	go elevator.Timer_runTimer(drv_timeOut, elev)
	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	if elevio.GetFloor() == -1 {
		elev = elevator.Fsm_onInitBetweenFloors(elev)
		fmt.Printf("INNE I INIT BETWEEN")
	} else {
		elev = elevator.Fsm_init(elev)
	}

	//prev := [nFloors][nButtons]int{}

	//elev.Timer := time.NewTimer(time.Second*1)

	//lage en for select hvor drv_button sender istedet for request buttons
	for {

		select {
		case floor := <-drv_floors:
			fmt.Printf("%+v\n", floor)
			elev = elevator.Fsm_onFloorArrival(floor, elev)

			fmt.Printf(" ---------case floor----------")

		case button := <-drv_buttons:
			elev = elevator.Fsm_onRequestButtonPress(button.Floor, button.Button, elev)
			fmt.Printf("inne i poll buttons")

		case timer := <-drv_timeOut:
			elev = elevator.Fsm_onDoorTimeout(elev)
			fmt.Println("-------------------I TIMER CASE:", timer)

		case obstruction := <-drv_obstr:
			if obstruction {
				elevio.SetMotorDirection(elevio.MD_Stop)
			}
			fmt.Printf(" ---------case obstruction----------")

		}
		time.Sleep(time.Duration(inputPollRateMs) * time.Millisecond)
	}

}
