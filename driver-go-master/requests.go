package elev

import "Driver-go/elevio"

type DirnBehaviourPair struct {
	dirn      elevio.MotorDirection //Kanskje lage en egen type som er Direction
	behaviour ElevatorBehaviour
}

func Requests_above(e Elevator) int {
	for floor := e.floor + 1; floor < numFloors; floor++ {
		for btn := 0; btn < numButtons; btn++ {
			if e.requests[floor][btn] >= 1 { //Antar at requests har verdi 1 om bestilling og null ellers
				return 1
			}
		}
	}

	return 0
}

func Requests_below(e Elevator) int {
	for floor := 0; floor < e.floor; floor++ {
		for btn := 0; btn < numButtons; btn++ {
			if e.requests[floor][btn] >= 1 {
				return 1
			}
		}
	}

	return 0
}

func Requests_here(e Elevator) int {

	for btn := 0; btn < numButtons; btn++ {
		if e.requests[e.floor][btn] >= 1 {
			return 1
		}
	}

	return 0
}

// ////////////////////////
func Requests_clearFloorOrders(floor int) {
	elevio.SetMotorDirection(elevio.MD_Stop)
	elevio.SetButtonLamp(elevio.BT_Cab, floor, false)

	//Hvilken hall-dir som skal cleares kommer an på heis-tilstand
	elevio.SetButtonLamp(elevio.BT_HallUp, floor, false)
	elevio.SetButtonLamp(elevio.BT_HallDown, floor, false)

	elevio.SetDoorOpenLamp(true)
}
