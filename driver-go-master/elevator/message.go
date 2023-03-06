package elevator

import "Driver-go/elevio"

type OrderMessageStruct struct{
	systemID string
	messageID string
	elevatorID string

	orderCounter int
	
	
	orderedButton elevio.ButtonEvent

	chosenElevator string
}



type IAmAliveMessageStruct struct{
	systemID string
	messageID string
	elevatorID string

	elevator Elevator

}