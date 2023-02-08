package elev

import "Driver-go/elevio"

var elevator Elevator

func Fsm_init() {
	elevator = Elevator_uninitialized()
}

func SetAllLights(es Elevator) {
}

// /////////////////////////////////////////////////////////////////////////////
func Fsm_onInitBetweenFloors() {
	elevio.SetMotorDirection(elevio.MD_Down)
}

func Fsm_onButtonPress(floor int, buttonType elevio.ButtonType) {
	//Skru på lys for å signalisere bestilling
	elevio.SetButtonLamp(buttonType, floor, true)
}

func Fsm_onFloorArrival(newFloor int) {
	//hvis knapp lyser i etg vi ankommer,
	//OG button-type er samme som heisens tilstand - stopp

}
