package timer

import (
	"fmt"
	"realtimeProject/project-gruppe64/system"
	"time"
)

const(
	messageTimerDuration = 2 //sek
	orderTimerDuration = 40 //sek
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
func RunMessageTimer(messageTimer <-chan system.SendingOrder, messageTimerTimedOut chan<- system.SendingOrder){
	timersRunningMap := make(map[system.SendingOrder]bool)
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
						fmt.Println("Message timer timed out")
						messageTimerTimedOut <- ord
						delete(timersRunningMap, ord)
					}
				})
			}
		default:
			break
		}
	}
}

func RunOrderTimer(orderTimer <-chan system.SendingOrder, orderTimerTimedOut chan<- system.SendingOrder){
	for{
		select {
		case ord := <-orderTimer:
			time.AfterFunc(time.Duration(orderTimerDuration)*time.Second, func() {
				fmt.Println("Order timer timed out")
				orderTimerTimedOut <- ord
			})
		default:
			break
		}
	}
}


