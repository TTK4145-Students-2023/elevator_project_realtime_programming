package singleElevator

import (
	"Driver-go/elevatorHardware"
)

//Contains functions for saving orders locally, handling new orders from database and clearing done orders.
//Set function are helperfunctions for setting the ordersatae and assigning the chosenElevator.

func SaveLocalNewOrder(button elevatorHardware.ButtonEvent, chosenElevator string) Elevator {
	elevatorObject.Requests[button.Floor][button.Button].OrderState = NewOrder
	elevatorObject.Requests[button.Floor][button.Button].AssingedElevatorID = chosenElevator
	return elevatorObject
}

func SaveLocalNoOrder(floor int, buttonType elevatorHardware.ButtonType) Elevator {
	elevatorObject.Requests[floor][buttonType].OrderState = NoOrder
	elevatorObject.Requests[floor][buttonType].AssingedElevatorID = ""
	return elevatorObject
}

func SaveLocalConfirmedOrder(button elevatorHardware.ButtonEvent, chosenElevator string) Elevator {
	elevatorObject.Requests[button.Floor][button.Button].OrderState = ConfirmedOrder
	elevatorObject.Requests[button.Floor][button.Button].AssingedElevatorID = chosenElevator
	SetAllLights(elevatorObject)
	return elevatorObject
}

func ClearOrderAtCurrentFloor(e Elevator) Elevator {
	elevatorObject = SaveLocalNoOrder(e.Floor, elevatorHardware.BT_Cab)

	switch e.Direction {
	case elevatorHardware.MD_Up:
		if !OrdersAbove(e) && CheckNoOrder(e, elevatorHardware.BT_HallUp) {
			e = SaveLocalNoOrder(e.Floor, elevatorHardware.BT_HallDown)
		}
		e = SaveLocalNoOrder(e.Floor, elevatorHardware.BT_HallUp)
	case elevatorHardware.MD_Down:
		if !OrdersBelow(e) && CheckNoOrder(e, elevatorHardware.BT_HallDown) {
			e = SaveLocalNoOrder(e.Floor, elevatorHardware.BT_HallUp)
		}
		e = SaveLocalNoOrder(e.Floor, elevatorHardware.BT_HallDown)
	case elevatorHardware.MD_Stop:
		fallthrough
	default:
		e = SaveLocalNoOrder(e.Floor, elevatorHardware.BT_HallUp)
		e = SaveLocalNoOrder(e.Floor, elevatorHardware.BT_HallDown)
	}
	return e
}

func ClearOrderAtThisFloor(arrivedElevatorID string, floor int) Elevator {

	if elevatorObject.Requests[floor][elevatorHardware.BT_HallDown].OrderState != NoOrder &&
		(arrivedElevatorID == elevatorObject.Requests[floor][elevatorHardware.BT_HallDown].AssingedElevatorID) {
		elevatorObject = SaveLocalNoOrder(floor, elevatorHardware.BT_HallDown)

	} else if elevatorObject.Requests[floor][elevatorHardware.BT_HallUp].OrderState != NoOrder &&
		(arrivedElevatorID == elevatorObject.Requests[floor][elevatorHardware.BT_HallUp].AssingedElevatorID) {
		elevatorObject = SaveLocalNoOrder(floor, elevatorHardware.BT_HallUp)
	}

	SetAllLights(elevatorObject)
	return elevatorObject
}
