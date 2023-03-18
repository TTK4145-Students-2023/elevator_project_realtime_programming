package manager

import (
	"Driver-go/elevio"
	"fmt"
)

func TransmitMessages(drv_buttons chan elevio.ButtonEvent, drv_floors chan int,
	msgTx chan MessageStruct, ackTx chan AckMessageStruct, database ElevatorDatabase,
	receivedAckChan chan AckMessageStruct, initiateSendAckChan chan AckMessageStruct) {
	for {
		select {
		case floor := <-drv_floors:

			floorMessage := MakeFloorMessage(floor)
			msgTx <- floorMessage //ackRx) //Ta inn ack kanalen og håndtere den inn i funksjonen?
			fmt.Println("sent floor message")
			TimeAcknowledgementAndResend(msgTx, floorMessage, database.NumElevators, receivedAckChan)

		case button := <-drv_buttons:
			chosenElevator := AssignOrderToElevator(database, button)
			orderMessage := MakeOrderMessage(chosenElevator, button)

			fmt.Println("I have received a button press at floor: ", orderMessage.OrderedButton.Floor)
			msgTx <- orderMessage
			fmt.Println("I have sent the order")

			TimeAcknowledgementAndResend(msgTx, orderMessage, database.NumElevators, receivedAckChan)
			fmt.Println("The ugly duckling has started")
		case ackToBeSent := <-initiateSendAckChan:
			ackTx <- ackToBeSent
		}
	}
}
