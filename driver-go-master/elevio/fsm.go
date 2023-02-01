package main

import "Driver-go/elevio"

func fsm_onInitBetweenFloors() {
	elevio.SetMotorDirection(elevio.MD_Down)
}

func fsm_onButtonPress(floor int, buttonType elevio.ButtonType) {
	//Skru på lys for å signalisere bestilling
	elevio.SetButtonLamp(buttonType, floor, true)
}

func fsm_onFloorArrival(newFloor int) {
	//hvis knapp lyser i etg vi ankommer,
	//OG button-type er samme som heisens tilstand - stopp

	ch_stop := make(chan bool)

	ch_stop <- elevio.GetButton(elevio.BT_Cab, newFloor)
	ch_stop <- elevio.GetButton(elevio.BT_HallUp, newFloor)

	select {
	case <-ch_stop:
		requests_clearFloorOrders(newFloor)
	default:
		elevio.SetMotorDirection(elevio.MD_Up)
	}
}
