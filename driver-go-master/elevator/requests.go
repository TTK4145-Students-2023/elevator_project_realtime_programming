package elevator

import (
	"Driver-go/elevio"
)

type DirnBehaviourPair struct {
	dirn      elevio.MotorDirection //Kanskje lage en egen type som er Direction
	behaviour ElevatorBehaviour
}

func Requests_above(e *Elevator) bool {
	for Floor := e.Floor + 1; Floor < numFloors; Floor++ {
		for btn := 0; btn < numButtons; btn++ {
			if e.requests[Floor][btn] { //Antar at requests har verdi 1 om bestilling og null ellers
				return true
			}
		}
	}

	return false
}

func Requests_below(e *Elevator) bool {
	for Floor := 0; Floor < e.Floor; Floor++ {
		for btn := 0; btn < numButtons; btn++ {
			if e.requests[Floor][btn] {
				return true
			}
		}
	}

	return false
}

func Requests_here(e *Elevator) bool {

	for btn := 0; btn < numButtons; btn++ {
		if e.requests[e.Floor][btn] {
			return true
		}
	}

	return false
}

func Requests_chooseDirection(e *Elevator) DirnBehaviourPair {

	switch e.dirn {
	case elevio.MD_Up:
		if Requests_above(e) {
			return DirnBehaviourPair{elevio.MD_Up, EB_Moving}
		} else if Requests_here(e) {
			return DirnBehaviourPair{elevio.MD_Stop, EB_DoorOpen} // stop?
		} else if Requests_below(e) {
			return DirnBehaviourPair{elevio.MD_Down, EB_Moving}
		} else {
			return DirnBehaviourPair{elevio.MD_Stop, EB_Idle}
		}
	case elevio.MD_Down:
		if Requests_below(e) {
			return DirnBehaviourPair{elevio.MD_Down, EB_Moving}
		} else if Requests_here(e) {
			return DirnBehaviourPair{elevio.MD_Stop, EB_DoorOpen} ///stop?
		} else if Requests_above(e) {
			return DirnBehaviourPair{elevio.MD_Up, EB_Moving}
		} else {
			return DirnBehaviourPair{elevio.MD_Stop, EB_Idle}
		}
	case elevio.MD_Stop:
		if Requests_here(e) {
			return DirnBehaviourPair{elevio.MD_Stop, EB_DoorOpen}
		} else if Requests_above(e) {
			return DirnBehaviourPair{elevio.MD_Up, EB_Moving}
		} else if Requests_below(e) {
			return DirnBehaviourPair{elevio.MD_Down, EB_Moving}
		} else {
			return DirnBehaviourPair{elevio.MD_Stop, EB_Idle}
		}
	default:
		return DirnBehaviourPair{elevio.MD_Stop, EB_Idle}
	}
}

func Requests_shouldStop(e *Elevator) bool {
	switch e.dirn {
	case elevio.MD_Down:
		return e.requests[e.Floor][elevio.BT_HallDown] ||
			e.requests[e.Floor][elevio.BT_Cab] ||
			!Requests_below(e)

	case elevio.MD_Up:
		return e.requests[e.Floor][elevio.BT_HallUp] ||
			e.requests[e.Floor][elevio.BT_Cab] ||
			!Requests_above(e)

	default:
		return true
	}
}

func Requests_shouldClearImmediately(e Elevator, btn_Floor int, btn_type elevio.ButtonType) bool {
	return e.Floor == btn_Floor && ((e.dirn == elevio.MD_Up && btn_type == elevio.BT_HallUp) || (e.dirn == elevio.MD_Down && btn_type == elevio.BT_HallDown) || (e.dirn == elevio.MD_Stop || btn_type == elevio.BT_Cab))

}

func Requests_clearAtCurrentFloor(e *Elevator) Elevator {

	e.requests[e.Floor][elevio.BT_Cab] = false
	switch e.dirn {
	case elevio.MD_Up:
		if !Requests_above(e) && !e.requests[e.Floor][elevio.BT_HallUp] {
			e.requests[e.Floor][elevio.BT_HallDown] = false
		}
		e.requests[e.Floor][elevio.BT_HallUp] = false
	case elevio.MD_Down:
		if !Requests_below(e) && !e.requests[e.Floor][elevio.BT_HallDown] {
			e.requests[e.Floor][elevio.BT_HallUp] = false
		}
		e.requests[e.Floor][elevio.BT_HallDown] = false
	case elevio.MD_Stop:
		fallthrough
	default:
		e.requests[e.Floor][elevio.BT_HallUp] = false
		e.requests[e.Floor][elevio.BT_HallDown] = false
	}
	return *e
}

// ////////////////////////
/*func Requests_clearFloorOrders(floor int) {
	elevio.SetMotorDirection(elevio.MD_Stop)
	elevio.SetButtonLamp(elevio.BT_Cab, floor, false)

	//Hvilken hall-dir som skal cleares kommer an på heis-tilstand
	elevio.SetButtonLamp(elevio.BT_HallUp, floor, false)
	elevio.SetButtonLamp(elevio.BT_HallDown, floor, false)

	elevio.SetDoorOpenLamp(true)
}
*/