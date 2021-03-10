package oversetting_fra_c
import(
	".\elevio"
)


//Requests_above checks if there are any requests above by first iterating through the floors above the current floor
//and checking if any of the buttons are pushed. Returns 1 if there are any requests, 0 if there are none.
func RequestsAbove(e Elevator) bool{
	for level:= e.floor+1; level < _numFloors; level++{ //level = floor, but floor was used in Elevator
		for button := 0; button < _numButtons; button ++{
			if e.requests[level][button] {
				return true
			}
		}

	}
	return false
}

//Requests_below checks if there are any requests below by first iterating through the floors below the current floor
//and checking if any of the buttons are pushed. Returns 1 if there are any requests, 0 if there are none.
func RequestsBelow(e Elevator) bool{
	for level:= 0; level < e.floor; level++ {
		for button := 0; button < _numButtons; button++{
			if e.requests[level][button] {
				return true
			}
		}
	}
	return false
}


//Requests_chooseDirection takes in an elevator and gives out in which direction the elevator should move based on what
//its direction is and where there are requests.
func RequestsChooseDirection(e Elevator) MotorDirection {
	switch e.MotorDirection  {
	case MD_Up:
		if RequestsAbove(e){
			return MD_UP
		} else if RequestsBelow(e) {
			return MD_Down
		} else {
			return MD_Stop
		}
	case MD_Down:
	case MD_Stop: // there should only be one request in this case. Checking up or down first is arbitrary.
		if RequestsBelow(e) {
			return MD_Down
		} else if RequestsAbove(e) {
			return MD_Up
		} else {
			return MD_Stop
		}
	default:
		return MD_Stop

	}
	return MD_Stop
}


//Requests_shouldStop takes in an elevator-object, and switches on its directions.
//if it goes down, and has any requests going down or stopping on the floor it is on, it must stop. If it goes down and
//has no requests below it (Requests_below), it should also stop. It is the same, but opposite way around if the
//motor direction is MD_Up.
func RequestsShouldStop(e Elevator) bool {
	switch e.MotorDirection {
	case MD_Down:
		eRequests := e.requests[e.floor][BT_HallDown] || e.requests[e.floor][BT_Cab]
		eNotRequests := RequestsBelow(e)
		if eRequests || eNotRequests == false{ //This gives me if eRequests == true, right?
			return true
		} else{
			return false
		}

	case MD_Up:
		eRequests := e.requests[e.floor][BT_HallUp] || e.requests[e.floor][BT_Cab]
		eNotRequests := RequestsAbove(e)
		if eRequests || eNotRequests == false{ //This gives me if eRequests == true, right?
			return true
		} else{
			return false
		}

	case MD_Stop:
	default:
		return true

	}
	return true
}



// Requests_clearAtCurrentFloor takes in an elevator and checks what clearRequestsVariant  it has (CV_All or CV_InDirn).
//Based on what clearRequestsVariant, it removes different requests that have been taken care of.
func RequestsClearAtCurrentFloor(e Elevator) Elevator {
	switch e.config.clearRequestVariant {
	case CV_All: //CV:clear request variant
		for button := 0; button < _numButtons; button++ { //_numButtons= 3
		e.requests[e.floor][button] = 0
	}

	case CV_InDirn:
		e.requests[e.floor][BT_Cab] = 0
		switch e.MotorDirection {
		case MD_Up:
			e.requests[e.floor][BT_HallUp] = 0
			if RequestsAbove(e) == false {
				e.requests[e.floor][BT_HallDown] = 0
			}

		case MD_Down:
			e.requests[e.floor][BT_HallDown] = 0
			if RequestsBelow(e) == false {
				e.requests[e.floor][BT_HallUp] = 0
			}

		case MD_Stop:
		default:
			e.requests[e.floor][BT_HallUp] = 0
			e.requests[e.floor][BT_HallDown] = 0
		}
	}

	return e
}




