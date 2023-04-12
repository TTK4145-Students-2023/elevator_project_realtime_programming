package manager

import (
	"Driver-go/elevator"
	"time"
)

func AliveMessageReceiver(aliveRx chan elevator.IAmAliveMessageStruct, database ElevatorDatabase, newOrder chan elevator.Elevator, confirmedOrderChan chan elevator.OrderMessageStruct) {
	for {
		aliveMessage := <-aliveRx
		if aliveMessage.ElevatorID != elevator.MyID {
			//elevatorFromSearch := SearchMessageOrderUpdate(aliveMessage, database, confirmedOrderChan)
			//newOrder <- elevatorFromSearch
		}

		time.Sleep(25 * time.Millisecond)
	}
}
