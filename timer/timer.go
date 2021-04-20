package timer

import (
	"fmt"
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

const(
	messageTimerDuration = 4 //sek
	orderTimerDuration = 20 //sek
)

func RunDoorTimer (doorTimerDuration <-chan float64, doorTimerTimedOut chan<- bool) {
	timerRunning := false
	stopTimerFromTimeOut := false
	for {
		select {
		case dTD := <-doorTimerDuration:
			if timerRunning {
				stopTimerFromTimeOut = true
				time.AfterFunc(time.Duration(dTD)*time.Second, func() {
					if stopTimerFromTimeOut {
						stopTimerFromTimeOut = false
					} else {
						timerRunning = false
						doorTimerTimedOut <- true
					}
				})
			} else {
				timerRunning = true
				time.AfterFunc(time.Duration(dTD)*time.Second, func() {
					if stopTimerFromTimeOut {
						stopTimerFromTimeOut = false
					} else {
						timerRunning = false
						doorTimerTimedOut <- true
					}
				})
			}
		}
	}
}




//send ordren når message timer startes (når ny melding sendes), og
//send ordren når acceptance message er mottatt (da slettes timer fra running timers).
//
func RunMessageTimer(messageTimer <-chan system.NetOrder, placedMessageReceived <-chan system.NetOrder, messageTimerTimedOut chan<- system.NetOrder){
	timersRunningMap := make(map[system.NetOrder]bool)

	for{
		select {
		case msgTmr := <-messageTimer:
			timersRunningMap[msgTmr] = true
			go time.AfterFunc(time.Duration(messageTimerDuration)*time.Second, func(){
				if timersRunningMap[msgTmr]{
					messageTimerTimedOut <- msgTmr
					delete(timersRunningMap, msgTmr)
				} else {
					delete(timersRunningMap, msgTmr)
				}
			})
		case plcdMsg := <-placedMessageReceived:
			timersRunningMap[plcdMsg] = false
		}
	}
}


func RunOrderTimer(orderTimer <-chan system.NetOrder, orderTimerTimedOut chan<- system.NetOrder){
	for{
		select {
		case ord := <-orderTimer:
			go time.AfterFunc(time.Duration(orderTimerDuration)*time.Second, func() {
				fmt.Println("Order timer timed out")
				orderTimerTimedOut <- ord
			})
		}
	}
}


