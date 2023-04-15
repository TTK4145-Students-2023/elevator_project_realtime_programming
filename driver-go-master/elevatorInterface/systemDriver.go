package elevatorInterface

import (
	"Driver-go/elevatorHardware"
	"Driver-go/databaseHandler"
	"Driver-go/singleElevator"
	"time"
)

func HandleNewFloorAndUpdateDatabase(floor int, database databaseHandler.ElevatorDatabase, doorTimer *time.Timer, immobilityTimer *time.Timer) databaseHandler.ElevatorDatabase {
	newElevatorUpdate := singleElevator.FloorArrival(floor, doorTimer, immobilityTimer)
	database = databaseHandler.UpdateDatabase(newElevatorUpdate, database)
	return database
}

func HandleNewButtonAndUpdateDatabase(button elevatorHardware.ButtonEvent, database databaseHandler.ElevatorDatabase, doorTimer *time.Timer, immobilityTimer *time.Timer) databaseHandler.ElevatorDatabase {
	chosenElevator := databaseHandler.AssignOrderToElevator(database, button)
	newElevatorUpdate := singleElevator.HandleNewOrder(chosenElevator, button, doorTimer, immobilityTimer)
	database = databaseHandler.UpdateDatabase(newElevatorUpdate, database)
	return database
}

func HandleObstruction(obstruction bool, doorTimer *time.Timer, immobilityTimer *time.Timer) {
	if singleElevator.IsDoorOpen() && obstruction {
		doorTimer.Stop()
		immobilityTimer.Reset(3 * time.Second)
	} else if !obstruction && singleElevator.IsDoorOpen() {
		immobilityTimer.Stop()
		singleElevator.SetWorkingState(singleElevator.Connected)
		doorTimer.Reset(3 * time.Second)
	}
}

func HandleStopButton(database databaseHandler.ElevatorDatabase) {
	singleElevator.ElevatorPrint(singleElevator.GetSingleEleavtorObject())
	for i := 0; i < len(database.ElevatorList); i++ {
		singleElevator.ElevatorPrint(database.ElevatorList[i])
	}
}
