package elevator

import (
	"Driver-go/elevio"
)

type DirnBehaviourPair struct {
	dirn      elevio.MotorDirection //Kanskje lage en egen type som er Direction
	behaviour ElevatorBehaviour
}

func Requests_above(e Elevator) bool {
	for floor := e.Floor + 1; floor < NumFloors; floor++ {
		for btn := 0; btn < NumButtons; btn++ {
			if e.Requests[floor][btn].OrderState == SO_Confirmed &&
				e.Requests[floor][btn].ElevatorID == e.ElevatorID { //Antar at requests har verdi 1 om bestilling og null ellers
				return true
			}
		}
	}

	return false
}

func Requests_below(e Elevator) bool {
	for floor := 0; floor < e.Floor; floor++ {
		for btn := 0; btn < NumButtons; btn++ {
			if e.Requests[floor][btn].OrderState == SO_Confirmed &&
				e.Requests[floor][btn].ElevatorID == e.ElevatorID {
				return true
			}
		}
	}

	return false
}

func Requests_here(e Elevator) bool {

	for btn := 0; btn < NumButtons; btn++ {
		if e.Requests[e.Floor][btn].OrderState == SO_Confirmed &&
			e.Requests[e.Floor][btn].ElevatorID == e.ElevatorID {
			return true
		}
	}

	return false
}

func Requests_chooseDirection(e Elevator) DirnBehaviourPair {
	switch e.Dirn {
	case elevio.MD_Up:
		if Requests_above(e) {
			return DirnBehaviourPair{elevio.MD_Up, EB_Moving}
		} else if Requests_here(e) {
			return DirnBehaviourPair{elevio.MD_Down, EB_DoorOpen}
		} else if Requests_below(e) {
			return DirnBehaviourPair{elevio.MD_Down, EB_Moving}
		} else {
			return DirnBehaviourPair{elevio.MD_Stop, EB_Idle}
		}
	case elevio.MD_Down:
		if Requests_below(e) {
			return DirnBehaviourPair{elevio.MD_Down, EB_Moving}
		} else if Requests_here(e) {
			return DirnBehaviourPair{elevio.MD_Up, EB_DoorOpen}
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

func Requests_shouldStop(e Elevator) bool {
	switch e.Dirn {
	case elevio.MD_Down:
		return (e.Requests[e.Floor][elevio.BT_HallDown].OrderState == SO_Confirmed && e.Requests[e.Floor][elevio.BT_HallDown].ElevatorID == e.ElevatorID) ||
			e.Requests[e.Floor][elevio.BT_Cab].OrderState == SO_Confirmed ||
			!Requests_below(e) //mulig vi må legge til ID-sjekk

	case elevio.MD_Up:
		return (e.Requests[e.Floor][elevio.BT_HallUp].OrderState == SO_Confirmed && e.Requests[e.Floor][elevio.BT_HallUp].ElevatorID == e.ElevatorID)||
			e.Requests[e.Floor][elevio.BT_Cab].OrderState == SO_Confirmed ||
			!Requests_above(e)

	default:
		return true
	}
}

func Requests_shouldClearImmediately(e Elevator, btn_floor int, btn_type elevio.ButtonType) bool {
	return e.Floor == btn_floor && ((e.Dirn == elevio.MD_Up && btn_type == elevio.BT_HallUp) ||
		(e.Dirn == elevio.MD_Down && btn_type == elevio.BT_HallDown) || e.Dirn == elevio.MD_Stop || btn_type == elevio.BT_Cab)

}

func Requests_clearAtCurrentFloor(e Elevator) Elevator {
	//Tanken: Alle går på heisen som stopper, så ordre må cleares uansett fordeling
	e.Requests[e.Floor][elevio.BT_Cab].OrderState = SO_NoOrder
	e.Requests[e.Floor][elevio.BT_Cab].ElevatorID = ""
	switch e.Dirn {
	case elevio.MD_Up:
		if !Requests_above(e) && e.Requests[e.Floor][elevio.BT_HallUp].OrderState == SO_NoOrder {
			e.Requests[e.Floor][elevio.BT_HallDown].OrderState = SO_NoOrder
			e.Requests[e.Floor][elevio.BT_HallDown].ElevatorID = ""
		}
		e.Requests[e.Floor][elevio.BT_HallUp].OrderState = SO_NoOrder
		e.Requests[e.Floor][elevio.BT_HallUp].ElevatorID = ""
	case elevio.MD_Down:
		if !Requests_below(e) && e.Requests[e.Floor][elevio.BT_HallDown].OrderState == SO_NoOrder {
			e.Requests[e.Floor][elevio.BT_HallUp].OrderState = SO_NoOrder
			e.Requests[e.Floor][elevio.BT_HallUp].ElevatorID = ""
		}
		e.Requests[e.Floor][elevio.BT_HallDown].OrderState = SO_NoOrder
		e.Requests[e.Floor][elevio.BT_HallDown].ElevatorID = ""
	case elevio.MD_Stop:
		fallthrough
	default:
		e.Requests[e.Floor][elevio.BT_HallUp].OrderState = SO_NoOrder
		e.Requests[e.Floor][elevio.BT_HallUp].ElevatorID = ""
		e.Requests[e.Floor][elevio.BT_HallDown].OrderState = SO_NoOrder
		e.Requests[e.Floor][elevio.BT_HallDown].ElevatorID = ""
	}
	return e
}

func Requests_clearOnFloor(arrivedElevatorID string, floor int) {
	//Trenger vel egt ikke å sjekke om det er en ordre her fordi hvis den er fordelt,
	//så er det jo en ordre der.
	//OBS! Må sjekke state til heis fordi det kan skje at den ikke skal cleare. Litt mer kopi av Req_clearAtCurrFloor(). Eks: hente ut state fra database

	if elevator.Requests[floor][elevio.BT_HallDown].OrderState == SO_Confirmed &&
		(arrivedElevatorID == elevator.Requests[floor][elevio.BT_HallDown].ElevatorID) {
		elevator.Requests[floor][elevio.BT_HallDown].OrderState = SO_NoOrder
		elevator.Requests[floor][elevio.BT_HallDown].ElevatorID = ""
		elevio.SetButtonLamp(elevio.BT_HallDown, floor, false) // La til denne men vet ikke hvorfor denne må være her siden setAlllights egentlig skal cleare lyset nederst
	} else if elevator.Requests[floor][elevio.BT_HallUp].OrderState == SO_Confirmed &&
		(arrivedElevatorID == elevator.Requests[floor][elevio.BT_HallUp].ElevatorID) {
		elevator.Requests[floor][elevio.BT_HallUp].OrderState = SO_NoOrder
		elevator.Requests[floor][elevio.BT_HallUp].ElevatorID = ""
		elevio.SetButtonLamp(elevio.BT_HallDown, floor, false) //HER OGSÅ
	}

	SetAllLights(elevator)
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
