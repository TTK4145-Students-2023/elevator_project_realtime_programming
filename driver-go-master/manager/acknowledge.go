package manager

import (
	"Driver-go/elevator"
	"fmt"
	"time"
)

func OrderBroadcast_timeAcknowledgementAndResend(initiateResend chan OrderMessageStruct,
	broadcastMessage OrderMessageStruct, ackRx chan AckMessageStruct, numOfElevatorsInNetwork int) {

	resendCap := 20
	numOfExpectedAcks := numOfElevatorsInNetwork

	acknowledgementCount := 0
	resendsCount := 0
	timer := time.NewTimer(50 * time.Millisecond)

	select {
	case <-ackRx:
		//Hvis det er en ack som er til meg
		acknowledgementCount++ //må man vite hvem som har sendt acks?
		fmt.Println("Order_Number of acknowledgements: ", acknowledgementCount)
		if acknowledgementCount == numOfExpectedAcks {
			timer.Stop()
			//break
		}

	case <-timer.C:
		initiateResend <- broadcastMessage
		fmt.Println("Resent ORDER-message")
		resendsCount++
		if resendsCount == resendCap {
			timer.Stop()
			//break
			//HUSK! Vi har resenda mange ganger. Sannsynlighet og alt det der... Vi forventer en endring i peers.
		}
		timer.Reset(50 * time.Millisecond)

	default:

	}

}

func FloorArrivalBroadcast_timeAcknowledgementAndResend(chanTx chan FloorArrivalMessageStruct,
	broadcastMessage FloorArrivalMessageStruct, ackRx chan AckMessageStruct, numOfElevatorsInNetwork int) {

	resendCap := 20
	numOfExpectedAcks := numOfElevatorsInNetwork

	acknowledgementCount := 0
	resendsCount := 0
	timer := time.NewTimer(50 * time.Millisecond)

	select {
	case ack := <-ackRx:
		//Hvis det er en ack som er til meg
		if ack.IDOfAckReciever == elevator.MyID {
			acknowledgementCount++ //må man vite hvem som har sendt acks?
			fmt.Println("Floor_Number of acknowledgements: ", acknowledgementCount)
			if acknowledgementCount == numOfExpectedAcks {
				timer.Stop()
				//break
			}
		} else {
			fmt.Println("I have received an ack that was none of my concern")
		}
	case <-timer.C:
		chanTx <- broadcastMessage
		fmt.Println("Resent FLOORARRIVAL-message")
		resendsCount++
		if resendsCount == resendCap {
			timer.Stop()
			//break
			//HUSK! Vi har resenda mange ganger. Sannsynlighet og alt det der... Vi forventer en endring i peers.
		}
		timer.Reset(50 * time.Millisecond)

	}

}

func SendAcknowledge(senderID string) AckMessageStruct {
	fmt.Println("An acknowledgement has been sent to: ", senderID)
	ackMessage := AckMessageStruct{
		IDOfAckReciever: senderID,
		MessageNumber:   5} //recievedMessage.MyElevator.OrderCounter
	return ackMessage
}
