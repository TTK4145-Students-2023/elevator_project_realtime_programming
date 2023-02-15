package main

import (
	"time"
)

var timerEndTime time.Time
var timerActive bool

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
