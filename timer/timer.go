package timer

import "time"

//GetWallTime returns the current time in milliseconds
func GetWallTime() float64{
	timeNow := time.Now()
	//timeSec := float64(timeNow.Second())
	nano := float64(timeNow.UnixNano())
	milli := nano/1000000 //milliseconds
	return milli
}

var TimerEndTime float64
var TimerActive bool

//timerStart takes in starts a timer , and puts the endtime for that timer into a channel, timerEndTime. It also states
//that timerActive is true.
func TimerStart (duration float64){
	TimerEndTime = GetWallTime() + duration //duration must here be given in milliseconds.
	TimerActive = true
}

//timerStop is used when we want the timer to stop being active.
func TimerStop(){
	TimerActive = false
}
func TimerTimedOut() bool{
	return TimerActive  &&  GetWallTime() > TimerEndTime
}