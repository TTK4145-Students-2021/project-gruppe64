package timer

import (
	"../system"
	"fmt"
	"time"
)

// GO-ROUTINE, main initiated
// Timer for elevator doors. Times out after door timer duration if not activated again.
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

// GO-ROUTINE, main initiated
// Timer for sending orders. Start timer when order is sent, and stop timer when "order placed"-message is
// received. Times out if placed message is not received within system-given time.
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

// GO-ROUTINE, main initiated
// Timer for order execution. Start timer when order is sent (also if to the elevator itself), and stop
// timer when order is executed. Times out if order not executed within system-given time.
func RunOrderTimer(orderTimerCh <-chan system.NetOrder, orderTimerTimedOutCh chan<- system.NetOrder){
	timersRunningMap := make(map[system.NetOrder]bool)
	for{
		select{
		case orderTimer := <-orderTimerCh:
			if !timersRunningMap[orderTimer] {
				fmt.Println("Order Timer started")
				timersRunningMap[orderTimer] = true
				go time.AfterFunc(time.Duration(system.OrderTimerDuration)*time.Second, func() {
					orderTimerTimedOutCh <- orderTimer
					delete(timersRunningMap, orderTimer)
				})
			}
		}
	}
}