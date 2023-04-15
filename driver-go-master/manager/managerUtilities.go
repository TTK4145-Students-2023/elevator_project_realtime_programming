package manager

import (
	"Driver-go/elevio"
	"Driver-go/singleElevator"
	"time"
)

func IsElevatorInDatabase(elevatorID string, database ElevatorDatabase) bool {
	for i := 0; i < len(database.ElevatorList); i++ {
		if database.ElevatorList[i].ElevatorID == elevatorID { //Sjekker at valgt heis ikke er unconnected
			return true
		}
	}
	return false
}

func MessageIDEqualsMyID(messageUpdateID string) bool {
	if messageUpdateID == singleElevator.MyID {
		return true
	} else {
		return false
	}
}

func shouldITakeTheOrder(order elevio.ButtonEvent) bool {
	if order.Button == elevio.BT_Cab || singleElevator.GetIsAlone() || singleElevator.AvailableAtCurrFloor(order.Floor) {
		return true
	} else {
		return false
	}
}

func GetElevatorFromID(database ElevatorDatabase, elevatorID string) singleElevator.Elevator {
	var e singleElevator.Elevator
	for i := 0; i < len(database.ElevatorList); i++ {
		if database.ElevatorList[i].ElevatorID == elevatorID {
			return database.ElevatorList[i]
		}
	}
	return e
}

func SendCabCalls(cabsToBeSent []OrderStruct, cabsChannelTx chan OrderStruct) {
	for k := 0; k < len(cabsToBeSent); k++ {
		cabsChannelTx <- cabsToBeSent[k]
		time.Sleep(time.Duration(25) * time.Millisecond)
	}
}
