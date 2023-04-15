package elevatorInterface

import (
	"Driver-go/elevio"
	"Driver-go/manager"
	"Driver-go/singleElevator"
	"time"
)



func HandleNewFloorAndUpdateDatabase(floor int, database manager.ElevatorDatabase, doorTimer *time.Timer, immobilityTimer *time.Timer) manager.ElevatorDatabase {
	newElevatorUpdate := singleElevator.Fsm_onFloorArrival(floor, doorTimer, immobilityTimer)
	database = manager.UpdateDatabase(newElevatorUpdate, database)
	return database
}

func HandleNewButtonAndUpdateDatabase(button elevio.ButtonEvent, database manager.ElevatorDatabase, doorTimer *time.Timer, immobilityTimer *time.Timer) manager.ElevatorDatabase {
	chosenElevator := manager.AssignOrderToElevator(database, button)
	newElevatorUpdate := singleElevator.HandleNewOrder(chosenElevator, button, doorTimer, immobilityTimer)
	database = manager.UpdateDatabase(newElevatorUpdate, database)
	return database
}

func HandleObstruction(obstruction bool, doorTimer *time.Timer, immobilityTimer *time.Timer) {
	if singleElevator.IsDoorOpen() && obstruction {
		doorTimer.Stop()
		immobilityTimer.Reset(3 * time.Second)
	} else if !obstruction && singleElevator.IsDoorOpen() {
		immobilityTimer.Stop()
		singleElevator.SetWorkingState(singleElevator.WS_Connected)
		doorTimer.Reset(3 * time.Second)
	}
}

func HandleStopButton(database manager.ElevatorDatabase) {
	singleElevator.ElevatorPrint(singleElevator.GetSingleEleavtorObject())
	for i := 0; i < len(database.ElevatorList); i++ {
		singleElevator.ElevatorPrint(database.ElevatorList[i])
	}
}





