package requests

import (
	"realtimeProject/project-gruppe64/elevator"
	"realtimeProject/project-gruppe64/io"
)


//Requests_above checks if there are any requests above by first iterating through the floors above the current floor
//and checking if any of the buttons are pushed. Returns 1 if there are any requests, 0 if there are none.
func RequestsAbove(e elevator.Elevator) bool{
	for level := e.Floor + 1; level < elevator.NumFloors; level++ { //level = floor, but floor was used in Elevator
		for button := 0; button < elevator.NumButtons; button ++{
			if e.Requests[level][button] != 0{
				return true
			}
		}
	}
	return false
}

//Requests_below checks if there are any requests below by first iterating through the floors below the current floor
//and checking if any of the buttons are pushed. Returns 1 if there are any requests, 0 if there are none.
func RequestsBelow(e elevator.Elevator) bool{
	for level:= 0; level < e.Floor; level++ {
		for button := 0; button < elevator.NumButtons; button++{
			if e.Requests[level][button] != 0 {
				return true
			}
		}
	}
	return false
}


//Requests_chooseDirection takes in an elevator and gives out in which direction the elevator should move based on what
//its direction is and where there are requests.
func RequestsChooseDirection(e elevator.Elevator) io.MotorDirection {
	switch e.MotorDirection  {
	case io.MD_Up:
		if RequestsAbove(e){
			return io.MD_Up
		} else if RequestsBelow(e) {
			return io.MD_Down
		} else {
			return io.MD_Stop
		}
	case io.MD_Down:
		if RequestsBelow(e){
			return io.MD_Down
		} else if RequestsAbove(e) {
			return io.MD_Up
		} else {
			return io.MD_Stop
		}
	case io.MD_Stop: // there should only be one request in this case. Checking up or down first is arbitrary.
		if RequestsBelow(e) {
			return io.MD_Down
		} else if RequestsAbove(e) {
			return io.MD_Up
		} else {
			return io.MD_Stop
		}
	default:
		return io.MD_Stop
	}
	return io.MD_Stop
}


//Requests_shouldStop takes in an elevator-object, and switches on its directions.
//if it goes down, and has any requests going down or stopping on the floor it is on, it must stop. If it goes down and
//has no requests below it (Requests_below), it should also stop. It is the same, but opposite way around if the
//motor direction is MD_Up.
func RequestsShouldStop(e elevator.Elevator) bool {
	switch e.MotorDirection {
	case io.MD_Down:
		if e.Requests[e.Floor][io.BT_HallDown] != 0 || e.Requests[e.Floor][io.BT_Cab] != 0 || RequestsBelow(e) == false{ //This gives me if eRequests == true, right?
			return true
		} else{
			return false
		}

	case io.MD_Up:
		if e.Requests[e.Floor][io.BT_HallUp] != 0 || e.Requests[e.Floor][io.BT_Cab] != 0 || RequestsAbove(e) == false{ //This gives me if eRequests == true, right?
			return true
		} else{
			return false
		}

	case io.MD_Stop:
	default:
		return true

	}
	return true
}



// Requests_clearAtCurrentFloor takes in an elevator and checks what clearRequestsVariant  it has (CV_All or CV_InDirn).
//Based on what clearRequestsVariant, it removes different requests that have been taken care of.
func RequestsClearAtCurrentFloor(e elevator.Elevator) elevator.Elevator {
	switch e.Config.ClearRequestVariant {
	case elevator.CV_All: //CV:clear request variant
		for button := 0; button < elevator.NumButtons; button++ { //_numButtons= 3
			e.Requests[e.Floor][button] = 0
		}

	case elevator.CV_InDirn:
		e.Requests[e.Floor][io.BT_Cab] = 0
		switch e.MotorDirection {
		case io.MD_Up:
			e.Requests[e.Floor][io.BT_HallUp] = 0
			if RequestsAbove(e) == false {
				e.Requests[e.Floor][io.BT_HallDown] = 0
			}

		case io.MD_Down:
			e.Requests[e.Floor][io.BT_HallDown] = 0
			if RequestsBelow(e) == false {
				e.Requests[e.Floor][io.BT_HallUp] = 0
			}

		case io.MD_Stop:
		default:
			e.Requests[e.Floor][io.BT_HallUp] = 0
			e.Requests[e.Floor][io.BT_HallDown] = 0
		}
	}

	return e
}