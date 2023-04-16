package singleElevator

import (
	"Driver-go/elevatorHardware"
)

func checkNoOrder(elevator Elevator, btn elevatorHardware.ButtonType) bool {
	return elevator.Requests[elevator.Floor][btn].OrderState == NoOrder
}

func ordersAbove(e Elevator) bool {
	for floor := e.Floor + 1; floor < NumFloors; floor++ {
		for btn := 0; btn < NumButtons; btn++ {
			if e.Requests[floor][btn].OrderState == ConfirmedOrder &&
				e.Requests[floor][btn].ElevatorID == e.ElevatorID { //Antar at requests har verdi 1 om bestilling og null ellers
				return true
			}
		}
	}

	return false
}

func ordersBelow(e Elevator) bool {
	for floor := 0; floor < e.Floor; floor++ {
		for btn := 0; btn < NumButtons; btn++ {
			if e.Requests[floor][btn].OrderState == ConfirmedOrder &&
				e.Requests[floor][btn].ElevatorID == e.ElevatorID {
				return true
			}
		}
	}

	return false
}

func orderHere(e Elevator) bool {

	for btn := 0; btn < NumButtons; btn++ {
		if e.Requests[e.Floor][btn].OrderState == ConfirmedOrder &&
			e.Requests[e.Floor][btn].ElevatorID == e.ElevatorID {
			return true
		}
	}

	return false
}

func ordersChooseDirection(e Elevator) DirectionBehaviourPair {
	switch e.Direction {
	case elevatorHardware.MD_Up:
		if ordersAbove(e) {
			return DirectionBehaviourPair{elevatorHardware.MD_Up, Moving}
		} else if orderHere(e) {
			return DirectionBehaviourPair{elevatorHardware.MD_Down, DoorOpen}
		} else if ordersBelow(e) {
			return DirectionBehaviourPair{elevatorHardware.MD_Down, Moving}
		} else {
			return DirectionBehaviourPair{elevatorHardware.MD_Stop, Idle}
		}
	case elevatorHardware.MD_Down:
		if ordersBelow(e) {
			return DirectionBehaviourPair{elevatorHardware.MD_Down, Moving}
		} else if orderHere(e) {
			return DirectionBehaviourPair{elevatorHardware.MD_Up, DoorOpen}
		} else if ordersAbove(e) {
			return DirectionBehaviourPair{elevatorHardware.MD_Up, Moving}
		} else {
			return DirectionBehaviourPair{elevatorHardware.MD_Stop, Idle}
		}
	case elevatorHardware.MD_Stop:
		if orderHere(e) {
			return DirectionBehaviourPair{elevatorHardware.MD_Stop, DoorOpen}
		} else if ordersAbove(e) {
			return DirectionBehaviourPair{elevatorHardware.MD_Up, Moving}
		} else if ordersBelow(e) {
			return DirectionBehaviourPair{elevatorHardware.MD_Down, Moving}
		} else {
			return DirectionBehaviourPair{elevatorHardware.MD_Stop, Idle}
		}
	default:
		return DirectionBehaviourPair{elevatorHardware.MD_Stop, Idle}
	}
}

func elevatorShouldStop(e Elevator) bool {
	switch e.Direction {
	case elevatorHardware.MD_Down:
		return (e.Requests[e.Floor][elevatorHardware.BT_HallDown].OrderState == ConfirmedOrder &&
			e.Requests[e.Floor][elevatorHardware.BT_HallDown].ElevatorID == e.ElevatorID) ||
			e.Requests[e.Floor][elevatorHardware.BT_Cab].OrderState == ConfirmedOrder ||
			!ordersBelow(e)

	case elevatorHardware.MD_Up:
		return (e.Requests[e.Floor][elevatorHardware.BT_HallUp].OrderState == ConfirmedOrder &&
			e.Requests[e.Floor][elevatorHardware.BT_HallUp].ElevatorID == e.ElevatorID) ||
			e.Requests[e.Floor][elevatorHardware.BT_Cab].OrderState == ConfirmedOrder ||
			!ordersAbove(e)

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

func ClearOrderAtCurrentFloor(e Elevator) Elevator {
	e = setLocalNoOrder(e.Floor, elevatorHardware.BT_Cab)

	switch e.Direction {
	case elevatorHardware.MD_Up:
		if !ordersAbove(e) && checkNoOrder(e, elevatorHardware.BT_HallUp) {
			e = setLocalNoOrder(e.Floor, elevatorHardware.BT_HallDown)
		}
		e = setLocalNoOrder(e.Floor, elevatorHardware.BT_HallUp)
	case elevatorHardware.MD_Down:
		if !ordersBelow(e) && checkNoOrder(e, elevatorHardware.BT_HallDown) {
			e = setLocalNoOrder(e.Floor, elevatorHardware.BT_HallUp)
		}
		e = setLocalNoOrder(e.Floor, elevatorHardware.BT_HallDown)
	case elevatorHardware.MD_Stop:
		fallthrough
	default:
		e = setLocalNoOrder(e.Floor, elevatorHardware.BT_HallUp)
		e = setLocalNoOrder(e.Floor, elevatorHardware.BT_HallDown)
	}
	return e
}

func ClearOrderAtThisFloor(arrivedElevatorID string, floor int) Elevator {

	if elevatorObject.Requests[floor][elevatorHardware.BT_HallDown].OrderState != NoOrder &&
		(arrivedElevatorID == elevatorObject.Requests[floor][elevatorHardware.BT_HallDown].ElevatorID) {
		elevatorObject = setLocalNoOrder(floor, elevatorHardware.BT_HallDown)

	} else if elevatorObject.Requests[floor][elevatorHardware.BT_HallUp].OrderState != NoOrder &&
		(arrivedElevatorID == elevatorObject.Requests[floor][elevatorHardware.BT_HallUp].ElevatorID) {
		elevatorObject = setLocalNoOrder(floor, elevatorHardware.BT_HallUp)
	}

	SetAllLights(elevatorObject)
	return elevatorObject
}
