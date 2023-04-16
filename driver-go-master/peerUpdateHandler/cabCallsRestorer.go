package peerUpdateHandler

import (
	"Driver-go/databaseHandler"
	"Driver-go/elevatorHardware"
	"Driver-go/singleElevator"
	"time"
)

func SendCabCalls(cabsToBeSent []databaseHandler.OrderStruct, cabsChannelTx chan databaseHandler.OrderStruct) {
	for k := 0; k < len(cabsToBeSent); k++ {
		cabsChannelTx <- cabsToBeSent[k]
		time.Sleep(time.Duration(25) * time.Millisecond)
	}
}

func FindCabCallsForElevator(database databaseHandler.ElevatorDatabase, newPeer string) []databaseHandler.OrderStruct {
	var cabsToBeSent []databaseHandler.OrderStruct
	for i := 0; i < len(database.ElevatorList); i++ {
		if database.ElevatorList[i].ElevatorID == newPeer && newPeer != singleElevator.MyID {
			for floor := 0; floor < singleElevator.NumFloors; floor++ {
				if database.ElevatorList[i].Requests[floor][elevatorHardware.BT_Cab].ElevatorID == newPeer {
					var button elevatorHardware.ButtonEvent
					button.Floor = floor
					button.Button = elevatorHardware.BT_Cab
					panelPair := singleElevator.OrderpanelPair{ElevatorID: newPeer, OrderState: singleElevator.ConfirmedOrder}
					cabsToBeSent = append(cabsToBeSent, databaseHandler.MakeOrder(panelPair, button))
				}
			}
		}
	}
	return cabsToBeSent
}

func HandleRestoredCabs(newCabs databaseHandler.OrderStruct, doorTimer *time.Timer, immobilityTimer *time.Timer) singleElevator.Elevator {
	var newElevatorUpdate singleElevator.Elevator
	if databaseHandler.MessageIDEqualsMyID(newCabs.PanelPair.ElevatorID) {
		newElevatorUpdate = singleElevator.ExecuteAssignedOrder(newCabs.OrderedButton, singleElevator.MyID, doorTimer, immobilityTimer)
	}
	return newElevatorUpdate
}
