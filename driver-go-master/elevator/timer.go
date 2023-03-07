package elevator

import (
	"time"
)

var timerEndTime time.Time
var timerActive bool


func Timer_runTimer(receiver chan<- bool) {
	for {
		if elevator.DoorOpen {
			timer := time.NewTimer(3 * time.Second)
			<-timer.C
			receiver <- true
			elevator.DoorOpen = false //flytte?
		}
	}
	//OBS! Mangler håndtering av obstruksjon
 }
 


func Get_wall_time() float64 {
	return float64(time.Now().UnixNano()) / 1e9
}

func Timer_start(duration float64) {
	timerEndTime = time.Now().Add(time.Duration(duration * float64(time.Second)))
	timerActive = true
}

func Timer_stop() {
	timerActive = false
}

func Timer_timedOut() bool {
	return timerActive && time.Now().After(timerEndTime)
}
