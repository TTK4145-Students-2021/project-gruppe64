package oversetting_fra_c
//import(
//	".\elevio"
//)


func requests_above(Elevator e) int{
	for level:= e.floor+1; level < N_FLOORS; level++ //level = floor, but floor was used in Elevator
	{
		for int button := 0; button < N_BUTTONS; button++
		{
			if e.requests[level][button] {
				return 1
			}
		}

	}
	return 0
}

func requests_below(Elevator e) int{
	for level:= 0; level < e.floor; level++ {
		for button := 0; button < N_BUTTONS; button++{
			if e.requests[level][button] {
				return 1
			}
		}
	}
	return 0
}

func requests_chooseDirection(Elevator e) Dirn {
	switch e.Dirn {
	case D_Up:
		if requests_above(e){
			return D_Up
		} else if requests_below(e) {
			return D_Down
		} else {
			return D_Stop
		}
	case D_Down:
	case D_Stop: // there should only be one request in this case. Checking up or down first is arbitrary.
		if requests_below(e) {
			return D_Down
		} else if requests_above(e) {
			return D_Up
		} else {
			return D_Stop
		}
	default:
		return D_Stop

	}
	return D_Stop
}

func requests_shouldStop(Elevator e) int {
	switch e.dirn {
	case D_Down:
		eRequests := e.requests[e.floor][B_HallDown] || e.requests[e.floor][B_Cab]
		eNotRequests := requests_below(e)
		if eRequests == 1 || eNotRequests == 0{
			return 1
		} else{
			return 0
		}
		//return eRequests || !eNotRequests
		//return e.requests[e.floor][B_HallDown] || e.requests[e.floor][B_Cab] || (!requests_below(e))
	case D_Up:
		eRequests := e.requests[e.floor][B_HallUp] || e.requests[e.floor][B_Cab]
		eNotRequests := requests_above(e)
		if eRequests == 1 || eNotRequests == 0{
			return 1
		} else{
			return 0
		}
		//return eRequests || !eNotRequests
		//return e.requests[e.floor][B_HallUp] || e.requests[e.floor][B_Cab] || !requests_above(e)
	case D_Stop:
	default:
		return 1

	}
	return 1
}

func requests_clearAtCurrentFloor(Elevator e) Elevator {
	switch (e.config.clearRequestVariant) {
	case CV_All:
		for button := 0; button < N_BUTTONS; button++ {
		e.requests[e.floor][button] = 0
	}
		break

	case CV_InDirn:
		e.requests[e.floor][B_Cab] = 0;
		switch e.dirn {
		case D_Up:
			e.requests[e.floor][B_HallUp] = 0
			if requests_above(e) == 0 {
				e.requests[e.floor][B_HallDown] = 0
			}
			break
		case D_Stop:
		default:
			e.requests[e.floor][B_HallUp] = 0
			e.requests[e.floor][B_HallDown] = 0
		}
		break
	}

	return e
}



