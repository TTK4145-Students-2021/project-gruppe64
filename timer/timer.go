package main

import (
	"fmt"
	"sync"
	"time"
)


//getWallTime returns the current time in milliseconds
func getWallTime() float64{
	timeNow := time.Now()
	//timeSec := float64(timeNow.Second())
	nano := float64(timeNow.UnixNano())
	milli := nano/1000000 //milliseconds
	return milli
}

//timerStart takes in starts a timer , and puts the endtime for that timer into a channel, timerEndTime. It also states
//that timerActive is true.
func timerStart (duration float64, timerActive chan <- bool, timerEndTime chan <- float64){
	//defer wg.Done()
	timerEndTime <- getWallTime() + duration //duration must here be given in milliseconds.
	timerActive <- true
}

//timerStop is used when we want the timer to stop being active.
func timerStop(timerActive chan <- bool, timerEndTime chan <- float64){
	timerActive <- false
}


//timerTimedOut is used as to find out whether we have surpassed the duration-limit. If we have, it returns true.
func Countdown(duration float64, wg *sync.WaitGroup,timerActive chan bool, timerEndTime  chan float64) {
	defer wg.Done()
	timerStart(duration, timerActive, timerEndTime)
	endTime := <-timerEndTime
	active := <- timerActive
	for{
		if active == true {
			if getWallTime() > endTime{
				fmt.Println("Done boi")
				break
				//can also just return true here, if we don't want to have a for loop but poll instead. Thought for loop
				//was best for goroutines.
			}

		}
	}
}
func IsTimerDone(timerActive chan bool) bool{
	active := <-timerActive
	if active == true{
		return false //if it is not done return false
	}
	return true
}


func main(){

	var wg sync.WaitGroup
	timerActive := make(chan bool,1)
	timerEndTime := make(chan float64,1)
	//timerStart(5000, timerActive, timerEndTime)

	for i := 0; i < 30; i++{
		wg.Add(1)
		go Countdown(float64(i*100), &wg, timerActive, timerEndTime)
	}


	wg.Wait()
	close(timerActive)
	close(timerEndTime)
	if IsTimerDone(timerActive){
		fmt.Println("Finished process")
	}


	//wg.Add(1)
	//go timerTimedOut(1000, &wg, timerActive, timerEndTime) //is 1 sec slower than what is written in timer, e.g. if 5000 ms is written
	// in timerStart it uses about 6 secs. Much slower if we use seconds instead of milliseconds though.
	//wg.Add(1)
	//go timerTimedOut(4000,&wg, timerActive, timerEndTime)
	//wg.Wait() //used to wait for the goroutine to finish

}

