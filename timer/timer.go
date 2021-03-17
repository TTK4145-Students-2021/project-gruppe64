package timer

import "time"

//duration in sec, here bug is. Will not start again since it is in sleep

func RunTimer (timerDur <-chan float64, timedOut chan<- bool) {
	timerRunning := false
	stopTimerFromTimeOut := false
	for {
		select {
		case tD := <-timerDur:
			if timerRunning {
				stopTimerFromTimeOut = true
				time.AfterFunc(time.Duration(tD)*time.Second, func() {
					if stopTimerFromTimeOut {
						stopTimerFromTimeOut = false
					} else {
						timerRunning = false
						timedOut <- true
					}
				})
			} else {
				timerRunning = true
				time.AfterFunc(time.Duration(tD)*time.Second, func() {
					if stopTimerFromTimeOut {
						stopTimerFromTimeOut = false
					} else {
						timerRunning = false
						timedOut <- true
					}
				})
			}
		default:
			break
		}
	}
}