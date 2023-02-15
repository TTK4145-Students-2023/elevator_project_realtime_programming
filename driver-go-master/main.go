package main

import (
	"Driver-go/elevio"
	"fmt"
	"time"
)

const nFloors = 4
const nButtons = 3

func main() {
	fmt.Println("Started!")

	inputPollRateMs := 25
	var orders [4]int

	elevio.Init("localhost:15657", numFloors)

	var d elevio.MotorDirection = elevio.MD_Up

	elevio.SetMotorDirection(d)

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	input := elevioGetInputDevice()

	if input.FloorSensor() == -1 {
		Fsm_onInitBetweenFloors()
	}

	prev := [nFloors][nButtons]int{}

	//lage en for select hvor drv_button sender istedet for request buttons
	for {
		// Request button
		for f := 0; f < nFloors; f++ {
			for b := 0; b < nButtons; b++ {
				v := input.RequestButton(f, b)
				if v != 0 && v != prev[f][b] {
					Fsm_onRequestButtonPress(f, b)
				}
				prev[f][b] = v
			}
		}

		// Floor sensor
		var prevFloor int = -1
		f := input.FloorSensor()
		if f != -1 && f != prevFloor {
			Fsm_onFloorArrival(f)
		}
		prevFloor = f

		// Timer
		if Timer_timedOut() {
			Timer_stop()
			Fsm_onDoorTimeout()
		}

		time.Sleep(time.Duration(inputPollRateMs) * time.Millisecond)
	}
}
