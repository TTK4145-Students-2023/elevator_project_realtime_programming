package cost

import(
	"Driver-go/elevator"
	"math"
	"Driver-go/elevio"
   
)

// Constants
const (
    baseCost             = 10.0
    ratePerUnitDistance  = 0.5
    directionChangeCost  = 5.0
    waitingTime          = 10.0
    waitingTimeRate      = 0.1
)


func calculateCost(e *elevator.Elevator, floor int) float64 {
    // Determine current location of elevator and direction
    

    currFloor := e.Floor
    currDir := e.Dirn

    // Calculate distance to requested floor
    distance := math.Abs(float64(currFloor - floor))

    // Calculate cost based on distance
    cost := distance * ratePerUnitDistance

    // If elevator needs to change direction, add direction change cost
    if currDir != elevio.MD_Stop && currDir != getDirection(currFloor, floor) {
        cost += directionChangeCost
    }

    // Add any additional costs
    cost += waitingTimeCost(e)

    return cost
}


// Helper function to calculate direction to travel in
func getDirection(fromFloor, toFloor int) elevio.MotorDirection {
    if fromFloor < toFloor {
        return elevio.MD_Up
    } else if fromFloor > toFloor {
        return elevio.MD_Down
    }
    return elevio.MD_Stop
}

// Helper function to calculate waiting time cost
func waitingTimeCost(e *elevator.Elevator) float64 {
    if e.Behaviour == elevator.EB_Idle || e.Behaviour == elevator.EB_DoorOpen {
        return waitingTime * waitingTimeRate
    }
    return 0
}

