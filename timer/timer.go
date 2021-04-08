package timer

import (
	"realtimeProject/project-gruppe64/distributor"
	"time"
)

const(
	messageTimerDur = 1 //sek
)

func RunBlockingTimer (timerDur <-chan float64, timedOut chan<- bool) {
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

//her må timer for Plassert melding startes. Timer må også ha info om ordren.
//Når plassert kommer -> timer for ordren i seg selv må startes.
//Om ikke kommer -> ordren plasseres til en selv.

//Når timeren for ordren i seg selv går ut sjekkes det om den er slettet fra structen til den heisen.
//Om ordren fortsatt er der; ta den selv.


//send ordren når message timer startes (når ny melding sendes), og
//send ordren når acceptance message er mottatt (da slettes timer fra running timers).
//
func RunMessageTimer(orderToMessageTimer <-chan distributor.SendingOrder, orderMessageTimedOut chan<- distributor.SendingOrder){
	timersRunningMap := map[distributor.SendingOrder]bool{}
	for{
		select {
		case ord := <-orderToMessageTimer:
			_, found := timersRunningMap[ord]
			if found{ //Om her så er casen at vi har mottatt accepted message
				delete(timersRunningMap, ord)
			} else {
				timersRunningMap[ord] = true //setter opp en timer her
				time.AfterFunc(time.Duration(messageTimerDur)*time.Second, func() {
					_, found = timersRunningMap[ord]
					if found {
						orderMessageTimedOut <- ord
						delete(timersRunningMap, ord)
					}
				})
			}
		}
	}
}

