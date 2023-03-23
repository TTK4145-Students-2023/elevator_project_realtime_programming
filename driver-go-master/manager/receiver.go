package manager

import (
	"Driver-go/elevator"
	"time"
)

func AliveMessageReceiver(aliveRx chan elevator.IAmAliveMessageStruct, database ElevatorDatabase, newOrder chan elevator.Elevator, confirmedOrderChan chan elevator.OrderMessageStruct) {
	for {
		aliveMsg := <-aliveRx
		database = UpdateDatabase(aliveMsg.Elevator, database)
		if aliveMsg.ElevatorID != elevator.MyID {
			elevatorFromSearch := SearchMessageOrderUpdate(aliveMsg, database, confirmedOrderChan)
			newOrder <- elevatorFromSearch
		}
		
		time.Sleep(25 * time.Millisecond)
	}
}
