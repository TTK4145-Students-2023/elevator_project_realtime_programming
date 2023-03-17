package manager

import (
	"fmt"
	"time"
)

func TimeAcknowledgementAndResend(chanTx chan MessageStruct,
	broadcastMessage MessageStruct, ackRx chan AckMessageStruct, numOfElevatorsInNetwork int) {

	resendCap := 20
	numOfExpectedAcks := numOfElevatorsInNetwork

	acknowledgementCount := 0
	resendsCount := 0
	timer := time.NewTimer(30 * time.Millisecond)

	select {
	case <-ackRx:
		//Hvis det er en ack som er til meg
		acknowledgementCount++ //må man vite hvem som har sendt acks?
		if acknowledgementCount == numOfExpectedAcks {
			timer.Stop()
			//break
		}
	case <-timer.C:
		chanTx <- broadcastMessage
		fmt.Println("REsent message")
		resendsCount++
		if resendsCount == resendCap {
			timer.Stop()
			//break
			//HUSK! Vi har resenda mange ganger. Sannsynlighet og alt det der... Vi forventer en endring i peers.
		}
		timer.Reset(30 * time.Millisecond)

	}

}

func SendAcknowledge(ackTx chan AckMessageStruct, recievedMessage MessageStruct) {
	ackMessage := AckMessageStruct{
		IDOfAckReciever: recievedMessage.SenderID,
		MessageNumber:   5} //recievedMessage.MyElevator.OrderCounter
	ackTx <- ackMessage
}
