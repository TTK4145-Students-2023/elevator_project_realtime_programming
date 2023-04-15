package singleElevator

import (
	"Driver-go/elevatorHardware"
)

func Requests_above(e Elevator) bool {
	for floor := e.Floor + 1; floor < NumFloors; floor++ {
		for btn := 0; btn < NumButtons; btn++ {
			if e.Requests[floor][btn].OrderState == Confirmed &&
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
			if e.Requests[floor][btn].OrderState == Confirmed &&
				e.Requests[floor][btn].ElevatorID == e.ElevatorID {
				return true
			}
		}
	}

	return false
}

func Requests_here(e Elevator) bool {

	for btn := 0; btn < NumButtons; btn++ {
		if e.Requests[e.Floor][btn].OrderState == Confirmed &&
			e.Requests[e.Floor][btn].ElevatorID == e.ElevatorID {
			return true
		}
	}

	return false
}

func Requests_chooseDirection(e Elevator) DirectionBehaviourPair {
	switch e.Direction {
	case elevatorHardware.MD_Up:
		if Requests_above(e) {
			return DirectionBehaviourPair{elevatorHardware.MD_Up, Moving}
		} else if Requests_here(e) {
			return DirectionBehaviourPair{elevatorHardware.MD_Down, DoorOpen}
		} else if Requests_below(e) {
			return DirectionBehaviourPair{elevatorHardware.MD_Down, Moving}
		} else {
			return DirectionBehaviourPair{elevatorHardware.MD_Stop, Idle}
		}
	case elevatorHardware.MD_Down:
		if Requests_below(e) {
			return DirectionBehaviourPair{elevatorHardware.MD_Down, Moving}
		} else if Requests_here(e) {
			return DirectionBehaviourPair{elevatorHardware.MD_Up, DoorOpen}
		} else if Requests_above(e) {
			return DirectionBehaviourPair{elevatorHardware.MD_Up, Moving}
		} else {
			return DirectionBehaviourPair{elevatorHardware.MD_Stop, Idle}
		}
	case elevatorHardware.MD_Stop:
		if Requests_here(e) {
			return DirectionBehaviourPair{elevatorHardware.MD_Stop, DoorOpen}
		} else if Requests_above(e) {
			return DirectionBehaviourPair{elevatorHardware.MD_Up, Moving}
		} else if Requests_below(e) {
			return DirectionBehaviourPair{elevatorHardware.MD_Down, Moving}
		} else {
			return DirectionBehaviourPair{elevatorHardware.MD_Stop, Idle}
		}
	default:
		return DirectionBehaviourPair{elevatorHardware.MD_Stop, Idle}
	}
}

func Requests_shouldStop(e Elevator) bool {
	switch e.Direction {
	case elevatorHardware.MD_Down:
		return (e.Requests[e.Floor][elevatorHardware.BT_HallDown].OrderState == Confirmed &&
			e.Requests[e.Floor][elevatorHardware.BT_HallDown].ElevatorID == e.ElevatorID) ||
			e.Requests[e.Floor][elevatorHardware.BT_Cab].OrderState == Confirmed ||
			!Requests_below(e)

	case elevatorHardware.MD_Up:
		return (e.Requests[e.Floor][elevatorHardware.BT_HallUp].OrderState == Confirmed &&
			e.Requests[e.Floor][elevatorHardware.BT_HallUp].ElevatorID == e.ElevatorID) ||
			e.Requests[e.Floor][elevatorHardware.BT_Cab].OrderState == Confirmed ||
			!Requests_above(e)

	default:
		return true
	}
}

func Requests_shouldClearImmediately(e Elevator, btn_floor int, btn_type elevatorHardware.ButtonType) bool {
	return e.Floor == btn_floor &&
		((e.Direction == elevatorHardware.MD_Up && btn_type == elevatorHardware.BT_HallUp) ||
			(e.Direction == elevatorHardware.MD_Down && btn_type == elevatorHardware.BT_HallDown) ||
			e.Direction == elevatorHardware.MD_Stop || btn_type == elevatorHardware.BT_Cab)
}

func Requests_clearAtCurrentFloor(e Elevator) Elevator {
	e = setNoOrder(e, e.Floor, elevatorHardware.BT_Cab)

	switch e.Direction {
	case elevatorHardware.MD_Up:
		if !Requests_above(e) && checkNoOrder(e, elevatorHardware.BT_HallUp) {
			e = setNoOrder(e, e.Floor, elevatorHardware.BT_HallDown)
		}
		e = setNoOrder(e, e.Floor, elevatorHardware.BT_HallUp)
	case elevatorHardware.MD_Down:
		if !Requests_below(e) && checkNoOrder(e, elevatorHardware.BT_HallDown) {
			e = setNoOrder(e, e.Floor, elevatorHardware.BT_HallUp)
		}
		e = setNoOrder(e, e.Floor, elevatorHardware.BT_HallDown)
	case elevatorHardware.MD_Stop:
		fallthrough
	default:
		e = setNoOrder(e, e.Floor, elevatorHardware.BT_HallUp)
		e = setNoOrder(e, e.Floor, elevatorHardware.BT_HallDown)
	}
	return e
}

func Requests_clearOnFloor(arrivedElevatorID string, floor int) Elevator {

	if elevatorObject.Requests[floor][elevatorHardware.BT_HallDown].OrderState != NoOrder &&
		(arrivedElevatorID == elevatorObject.Requests[floor][elevatorHardware.BT_HallDown].ElevatorID) {
		elevatorObject = setNoOrder(elevatorObject, floor, elevatorHardware.BT_HallDown)

	} else if elevatorObject.Requests[floor][elevatorHardware.BT_HallUp].OrderState != NoOrder &&
		(arrivedElevatorID == elevatorObject.Requests[floor][elevatorHardware.BT_HallUp].ElevatorID) {
		elevatorObject = setNoOrder(elevatorObject, floor, elevatorHardware.BT_HallUp)
	}

	SetAllLights(elevatorObject)
	return elevatorObject
}
