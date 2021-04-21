package timer

import (
	"realtimeProject/project-gruppe64/system"
	"time"
)

/*
import (
	"../system"
	"fmt"
	"time"
)
 */

func RunDoorTimer (doorTimerDurationCh <-chan float64, doorTimerTimedOutCh chan<- bool) {
	timerRunning := false
	stopTimerFromTimeOut := false
	for {
		select {
		case doorTimerDuration := <-doorTimerDurationCh:
			if timerRunning {
				stopTimerFromTimeOut = true
				time.AfterFunc(time.Duration(doorTimerDuration)*time.Second, func() {
					if stopTimerFromTimeOut {
						stopTimerFromTimeOut = false
					} else {
						timerRunning = false
						doorTimerTimedOutCh <- true
					}
				})
			} else {
				timerRunning = true
				time.AfterFunc(time.Duration(doorTimerDuration)*time.Second, func() {
					if stopTimerFromTimeOut {
						stopTimerFromTimeOut = false
					} else {
						timerRunning = false
						doorTimerTimedOutCh <- true
					}
				})
			}
		}
	}
}

func RunMessageTimer(messageTimerCh <-chan system.NetOrder, placedMessageReceivedCh <-chan system.NetOrder,
	messageTimerTimedOutCh chan<- system.NetOrder){
	timersRunningMap := make(map[system.NetOrder]bool)

	for{
		select {
		case messageTimer := <-messageTimerCh:
			timersRunningMap[messageTimer] = true
			go time.AfterFunc(time.Duration(system.MessageTimerDuration)*time.Second, func(){
				if timersRunningMap[messageTimer]{
					messageTimerTimedOutCh <- messageTimer
					delete(timersRunningMap, messageTimer)
				} else {
					delete(timersRunningMap, messageTimer)
				}
			})
		case placedMessageReceived := <-placedMessageReceivedCh:
			timersRunningMap[placedMessageReceived] = false
		}
	}
}

