package singleElevator

import (
	"Driver-go/elevatorHardware"
)

//Helper functions for returning booleans based on orders and elevator position.
//Used in elevatorActions and orderHandler to increase readability.

func CheckNoOrder(e Elevator, btn elevatorHardware.ButtonType) bool {
	return e.Requests[e.Floor][btn].OrderState == NoOrder
}

func OrdersAbove(e Elevator) bool {
	for floor := e.Floor + 1; floor < NumFloors; floor++ {
		for btn := 0; btn < NumButtons; btn++ {
			if e.Requests[floor][btn].OrderState == ConfirmedOrder &&
				e.Requests[floor][btn].AssingedElevatorID == e.ElevatorID { //Antar at requests har verdi 1 om bestilling og null ellers
				return true
			}
		}
	}

	return false
}

func OrdersBelow(e Elevator) bool {
	for floor := 0; floor < e.Floor; floor++ {
		for btn := 0; btn < NumButtons; btn++ {
			if e.Requests[floor][btn].OrderState == ConfirmedOrder &&
				e.Requests[floor][btn].AssingedElevatorID == e.ElevatorID {
				return true
			}
		}
	}

	return false
}

func OrderHere(e Elevator) bool {

	for btn := 0; btn < NumButtons; btn++ {
		if e.Requests[e.Floor][btn].OrderState == ConfirmedOrder &&
			e.Requests[e.Floor][btn].AssingedElevatorID == e.ElevatorID {
			return true
		}
	}

	return false
}

func OrdersChooseDirection(e Elevator) DirectionBehaviourPair {
	switch e.Direction {
	case elevatorHardware.MD_Up:
		if OrdersAbove(e) {
			return DirectionBehaviourPair{elevatorHardware.MD_Up, Moving}
		} else if OrderHere(e) {
			return DirectionBehaviourPair{elevatorHardware.MD_Down, DoorOpen}
		} else if OrdersBelow(e) {
			return DirectionBehaviourPair{elevatorHardware.MD_Down, Moving}
		} else {
			return DirectionBehaviourPair{elevatorHardware.MD_Stop, Idle}
		}
	case elevatorHardware.MD_Down:
		if OrdersBelow(e) {
			return DirectionBehaviourPair{elevatorHardware.MD_Down, Moving}
		} else if OrderHere(e) {
			return DirectionBehaviourPair{elevatorHardware.MD_Up, DoorOpen}
		} else if OrdersAbove(e) {
			return DirectionBehaviourPair{elevatorHardware.MD_Up, Moving}
		} else {
			return DirectionBehaviourPair{elevatorHardware.MD_Stop, Idle}
		}
	case elevatorHardware.MD_Stop:
		if OrderHere(e) {
			return DirectionBehaviourPair{elevatorHardware.MD_Stop, DoorOpen}
		} else if OrdersAbove(e) {
			return DirectionBehaviourPair{elevatorHardware.MD_Up, Moving}
		} else if OrdersBelow(e) {
			return DirectionBehaviourPair{elevatorHardware.MD_Down, Moving}
		} else {
			return DirectionBehaviourPair{elevatorHardware.MD_Stop, Idle}
		}
	default:
		return DirectionBehaviourPair{elevatorHardware.MD_Stop, Idle}
	}
}

func ElevatorShouldStop(e Elevator) bool {
	switch e.Direction {
	case elevatorHardware.MD_Down:
		return (e.Requests[e.Floor][elevatorHardware.BT_HallDown].OrderState == ConfirmedOrder &&
			e.Requests[e.Floor][elevatorHardware.BT_HallDown].AssingedElevatorID == e.ElevatorID) ||
			e.Requests[e.Floor][elevatorHardware.BT_Cab].OrderState == ConfirmedOrder ||
			!OrdersBelow(e)

	case elevatorHardware.MD_Up:
		return (e.Requests[e.Floor][elevatorHardware.BT_HallUp].OrderState == ConfirmedOrder &&
			e.Requests[e.Floor][elevatorHardware.BT_HallUp].AssingedElevatorID == e.ElevatorID) ||
			e.Requests[e.Floor][elevatorHardware.BT_Cab].OrderState == ConfirmedOrder ||
			!OrdersAbove(e)

	default:
		return true
	}
}

func OrderShouldClearImmediately(e Elevator, btn_floor int, btn_type elevatorHardware.ButtonType) bool {
	return e.Floor == btn_floor &&
		((e.Direction == elevatorHardware.MD_Up && btn_type == elevatorHardware.BT_HallUp) ||
			(e.Direction == elevatorHardware.MD_Down && btn_type == elevatorHardware.BT_HallDown) ||
			e.Direction == elevatorHardware.MD_Stop || btn_type == elevatorHardware.BT_Cab)
}
