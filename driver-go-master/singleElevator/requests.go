package singleElevator

import (
	"Driver-go/elevio"
)

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

func Requests_chooseDirection(e Elevator) DirectionBehaviourPair {
	switch e.Direction {
	case elevio.MD_Up:
		if Requests_above(e) {
			return DirectionBehaviourPair{elevio.MD_Up, EB_Moving}
		} else if Requests_here(e) {
			return DirectionBehaviourPair{elevio.MD_Down, EB_DoorOpen}
		} else if Requests_below(e) {
			return DirectionBehaviourPair{elevio.MD_Down, EB_Moving}
		} else {
			return DirectionBehaviourPair{elevio.MD_Stop, EB_Idle}
		}
	case elevio.MD_Down:
		if Requests_below(e) {
			return DirectionBehaviourPair{elevio.MD_Down, EB_Moving}
		} else if Requests_here(e) {
			return DirectionBehaviourPair{elevio.MD_Up, EB_DoorOpen}
		} else if Requests_above(e) {
			return DirectionBehaviourPair{elevio.MD_Up, EB_Moving}
		} else {
			return DirectionBehaviourPair{elevio.MD_Stop, EB_Idle}
		}
	case elevio.MD_Stop:
		if Requests_here(e) {
			return DirectionBehaviourPair{elevio.MD_Stop, EB_DoorOpen}
		} else if Requests_above(e) {
			return DirectionBehaviourPair{elevio.MD_Up, EB_Moving}
		} else if Requests_below(e) {
			return DirectionBehaviourPair{elevio.MD_Down, EB_Moving}
		} else {
			return DirectionBehaviourPair{elevio.MD_Stop, EB_Idle}
		}
	default:
		return DirectionBehaviourPair{elevio.MD_Stop, EB_Idle}
	}
}

func Requests_shouldStop(e Elevator) bool {
	switch e.Direction {
	case elevio.MD_Down:
		return (e.Requests[e.Floor][elevio.BT_HallDown].OrderState == SO_Confirmed &&
			e.Requests[e.Floor][elevio.BT_HallDown].ElevatorID == e.ElevatorID) ||
			e.Requests[e.Floor][elevio.BT_Cab].OrderState == SO_Confirmed ||
			!Requests_below(e)

	case elevio.MD_Up:
		return (e.Requests[e.Floor][elevio.BT_HallUp].OrderState == SO_Confirmed &&
			e.Requests[e.Floor][elevio.BT_HallUp].ElevatorID == e.ElevatorID) ||
			e.Requests[e.Floor][elevio.BT_Cab].OrderState == SO_Confirmed ||
			!Requests_above(e)

	default:
		return true
	}
}

func Requests_shouldClearImmediately(e Elevator, btn_floor int, btn_type elevio.ButtonType) bool {
	return e.Floor == btn_floor &&
		((e.Direction == elevio.MD_Up && btn_type == elevio.BT_HallUp) ||
			(e.Direction == elevio.MD_Down && btn_type == elevio.BT_HallDown) ||
			e.Direction == elevio.MD_Stop || btn_type == elevio.BT_Cab)
}

func Requests_clearAtCurrentFloor(e Elevator) Elevator {
	e = setNoOrder(e, e.Floor, elevio.BT_Cab)

	switch e.Direction {
	case elevio.MD_Up:
		if !Requests_above(e) && checkNoOrder(e, elevio.BT_HallUp) {
			e = setNoOrder(e, e.Floor, elevio.BT_HallDown)
		}
		e = setNoOrder(e, e.Floor, elevio.BT_HallUp)
	case elevio.MD_Down:
		if !Requests_below(e) && checkNoOrder(e, elevio.BT_HallDown) {
			e = setNoOrder(e, e.Floor, elevio.BT_HallUp)
		}
		e = setNoOrder(e, e.Floor, elevio.BT_HallDown)
	case elevio.MD_Stop:
		fallthrough
	default:
		e = setNoOrder(e, e.Floor, elevio.BT_HallUp)
		e = setNoOrder(e, e.Floor, elevio.BT_HallDown)
	}
	return e
}

func Requests_clearOnFloor(arrivedElevatorID string, floor int) Elevator {

	if elevatorObject.Requests[floor][elevio.BT_HallDown].OrderState != SO_NoOrder &&
		(arrivedElevatorID == elevatorObject.Requests[floor][elevio.BT_HallDown].ElevatorID) {
			elevatorObject = setNoOrder(elevatorObject, floor, elevio.BT_HallDown)

	} else if elevatorObject.Requests[floor][elevio.BT_HallUp].OrderState != SO_NoOrder &&
		(arrivedElevatorID == elevatorObject.Requests[floor][elevio.BT_HallUp].ElevatorID) {
			elevatorObject = setNoOrder(elevatorObject, floor, elevio.BT_HallUp)
	}

	SetAllLights(elevatorObject)
	return elevatorObject
}
