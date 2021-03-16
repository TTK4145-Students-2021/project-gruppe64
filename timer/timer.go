package timer

import (
	"fmt"
	"time"
)


//duration in sec, here bug is. Will not start again since it is in sleep

func RunTimer (timerDur <-chan float64, timedOut chan<- bool){
	localTimerCh := make(chan bool)
	timerRunning := false
	stopTimerFromTimeOut := false

	for {
		select {
		case tD :=<- timerDur:
			if timerRunning{
				stopTimerFromTimeOut = true
				go startNewTimer(tD, localTimerCh)
			} else {
				timerRunning = true
				go startNewTimer(tD, localTimerCh)
			}
		case lT := <- localTimerCh:
			if lT {
				if stopTimerFromTimeOut{
					stopTimerFromTimeOut = false
				} else {
					timerRunning = false
					timedOut <- true
				}

			}
		default:
			break
		}
	}
}

func startNewTimer(duration float64, localTimer chan<- bool){
	fmt.Print("\n Timer started \n")
	time.Sleep(time.Duration(duration)*time.Second)
	fmt.Printf("\n Timer ended \n")
	localTimer <- true
}