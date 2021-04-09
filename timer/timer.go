package timer

import (
	"realtimeProject/project-gruppe64/distributor"
	"time"
)

const(
	messageTimerDuration = 1 //sek
	orderTimerDuration = 60 //sek
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
		default:
			break
		}
	}
}




//send ordren når message timer startes (når ny melding sendes), og
//send ordren når acceptance message er mottatt (da slettes timer fra running timers).
//
func RunMessageTimer(messageTimer <-chan distributor.SendingOrder, messageTimerTimedOut chan<- distributor.SendingOrder){
	timersRunningMap := map[distributor.SendingOrder]bool{}
	for{
		select {
		case ord := <-messageTimer:
			_, found := timersRunningMap[ord]
			if found{ //Om her så er casen at vi har mottatt accepted message
				delete(timersRunningMap, ord)
			} else {
				timersRunningMap[ord] = true //setter opp en timer her
				time.AfterFunc(time.Duration(messageTimerDuration)*time.Second, func() {
					_, found = timersRunningMap[ord]
					if found {
						messageTimerTimedOut <- ord
						delete(timersRunningMap, ord)
					}
				})
			}
		}
	}
}

func RunOrderTimer(orderTimer <-chan distributor.SendingOrder, orderTimerTimedOut chan<- distributor.SendingOrder){
	for{
		select {
		case ord := <-orderTimer:
			time.AfterFunc(time.Duration(orderTimerDuration)*time.Second, func() {
				orderTimerTimedOut <- ord
			})
		}
	}
}

