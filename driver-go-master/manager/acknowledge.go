package manager

import (
	"fmt"
	"time"
)

func TimeAcknowledgementAndResend(chanTx chan MessageStruct, broadcastMessage MessageStruct,
	numOfElevatorsInNetwork int, receivedAckChan chan AckMessageStruct) {

	resendCap := 20
	numOfExpectedAcks := numOfElevatorsInNetwork

	acknowledgementCount := 0
	resendsCount := 0
	timer := time.NewTimer(50 * time.Millisecond)

	select {
	case <-receivedAckChan:
		//Hvis det er en ack som er til meg
		acknowledgementCount++ //må man vite hvem som har sendt acks?
		fmt.Println("TimeAndResendFunc has received ", acknowledgementCount, " acknowledgement on 'receivedAckChan'")
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
		timer.Reset(50 * time.Millisecond)

	}

}

func SendAcknowledge(recievedMessage MessageStruct) AckMessageStruct {
	ackMessage := AckMessageStruct{
		IDOfAckReciever: recievedMessage.SenderID,
		MessageNumber:   5} //recievedMessage.MyElevator.OrderCounter
	return ackMessage
}
