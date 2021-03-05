package timer
import "time"

func get_wall_time() float64{
	time_now := time.Now()
	time_sec := float64(time_now.Second())
	nano := float64(time_now.UnixNano())
	nano = nano*0.000000001
	time_out := time_sec+nano

	return time_out //returning current time in seconds and nanoseconds I think, it is probably wrong since it is in int64.
}

var timerEndtime = 0.00
var timerActive = false

func timerStart (duration float64){
	timerEndtime = get_wall_time() + duration
	timerActive = true
}

func timerStop(){
	timerActive = false
}

func timer_timedOut() bool{
	return timerActive && get_wall_time()>timerEndtime
}

