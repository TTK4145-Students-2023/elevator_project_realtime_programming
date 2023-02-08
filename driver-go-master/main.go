package main

import (
	"Driver-go/elevio"
	"fmt"
)

func main() {

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

	for {
		select {

		//trykker på knapp, setter på lyset
		case btnLight := <-drv_buttons:
			fmt.Printf("%+v\n", btnLight)
			elevio.SetButtonLamp(btnLight.Button, btnLight.Floor, true)

		//går inn i case hver gang den kommer til en floor, skal printe, sjekke calls hvis det svarer til riktig state stopper den
		case newFloor := <-drv_floors:
			fmt.Printf("%+v\n", elevio.GetFloor())

			//Fsm_onFloorArrival(elevio.GetFloor())
			//fmt.Printf("%+v\n", a)
			if newFloor == numFloors-1 {
				d = elevio.MD_Down
			} else if newFloor == 0 {
				d = elevio.MD_Up
			}
			elevio.SetMotorDirection(d)
			elevio.SetFloorIndicator(newFloor)

			//if( state order på newfloor  == state orders[newFloor])

		//sjekker om døra har obstr, hvis ja holder åpen, hvis ikke tell ned fra sekunder
		case a := <-drv_obstr:
			fmt.Printf("%+v\n", a)
			if a {
				elevio.SetMotorDirection(elevio.MD_Stop)
			} else {
				elevio.SetMotorDirection(d)
			}

		//sjekker stop
		case stop := <-drv_stop:
			if stop {
				elevio.SetMotorDirection(elevio.MD_Stop)
			} else {
				elevio.SetMotorDirection(d)
			}

			//print
			fmt.Printf("%+v\n", stop)
			for f := 0; f < numFloors; f++ {
				for b := elevio.ButtonType(0); b < 3; b++ {
					elevio.SetButtonLamp(b, f, false)
				}
			}

		}
	}
}
