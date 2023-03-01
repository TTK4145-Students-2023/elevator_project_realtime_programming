package elevator

import (
	"fmt"
	"time"
)

var TimerEndTime time.Time
var TimerActive bool = false

func Timer_runTimer(receiver chan<- bool, elev Elevator) {
	for {
		if elev.doorOpen {
			fmt.Println("--------Timer_runTimer-----------------")
			timer := time.NewTimer(3 * time.Second)
			<-timer.C
			receiver <- true
			elev.doorOpen = false
		}
	}
	//OBS! Mangler håndtering av obstruksjon
}

/*
func Get_wall_time() float64 {
	return float64(time.Now().UnixNano()) / 1e9
}

func Timer_start(duration float64, TimerActive *bool) {
	TimerEndTime = time.Now().Add(time.Duration(duration * float64(time.Second)))
	*TimerActive = true
	fmt.Printf(" ---------timer startet-------")
}

func Timer_stop(TimerActive *bool) {
	*TimerActive = false
	fmt.Printf("-----timer stoppet----------------")
}

func Timer_timedOut(receiver chan<- bool, TimerActive *bool) {
	v := *TimerActive && time.Now().After(TimerEndTime)
	fmt.Print("------------goer timedOut-------------------")
	if v {
		receiver <- v
		fmt.Print("-------timedOut er true-------------")
	}

	//return timerActive && time.Now().After(timerEndTime)
}*/

//goroutine som får inn chan som enten skal være true hvis skal starte, false hvis skal stoppe. teller ned og returner en chan true for finished
/*func Timer(receiver chan<-start_timer) chan bool{
	timerEndTime = time.Now().Add(time.Duration(duration * float64(time.Second)))
	v := timerActive && time.Now().After(timerEndTime)
	reutn v;
}
*/
