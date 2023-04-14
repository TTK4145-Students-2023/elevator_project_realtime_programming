package manager

import (
	"Driver-go/elevator"
	"Driver-go/elevio"
)


func AssignOrderToElevator(database ElevatorDatabase, order elevio.ButtonEvent) string {

	lowCost := 100000.0
	lowestCostElevator := ""

	elevatorList := database.ElevatorList

	if shouldITakeTheOrder(order){
		lowestCostElevator = elevator.MyID
	}else {
		for i := 0; i < len(elevatorList); i++ {
			c := calculateCost(elevatorList[i], order)

			if c < lowCost && elevatorList[i].Operating == elevator.WS_Connected {
				lowCost = c
				lowestCostElevator = elevatorList[i].ElevatorID
			} else if c == lowCost && elevatorList[i].Operating == elevator.WS_Connected {
				
				var temp = database.ElevatorList[i].ElevatorID
				if temp < lowestCostElevator {
					lowCost = c
					lowestCostElevator = elevatorList[i].ElevatorID
				}
			}

		}

	}

	return lowestCostElevator
}

