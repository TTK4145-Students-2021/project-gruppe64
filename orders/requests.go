package oversetting_fra_c
import(
	".\elevio"
)


func requests_above(e Elevator) int{
	for level:= e.floor+1; level < _numFloors; level++//level = floor, but floor was used in Elevator
	{
		for int button := 0; button < 3; button ++ //N_BUTTONS == 3
		{
			if e.requests[level][button] {
				return 1
			}
		}

	}
	return 0
}

func requests_below(e Elevator) int{
	for level:= 0; level < e.floor; level++ {
		for button := 0; button < 3; button++{ //3 = N_buttons
			if e.requests[level][button] {
				return 1
			}
		}
	}
	return 0
}

func requests_chooseDirection(e Elevator) Dirn {
	switch e.Dirn {
	case MD_Up:
		if requests_above(e){
			return MD_UP
		} else if requests_below(e) {
			return MD_Down
		} else {
			return MD_Stop
		}
	case MD_Down:
	case MD_Stop: // there should only be one request in this case. Checking up or down first is arbitrary.
		if requests_below(e) {
			return MD_Down
		} else if requests_above(e) {
			return MD_Up
		} else {
			return MD_Stop
		}
	default:
		return MD_Stop

	}
	return MD_Stop
}

func requests_shouldStop(e Elevator) int {
	switch e.dirn {
	case MD_Down:
		eRequests := e.requests[e.floor][BT_HallDown] || e.requests[e.floor][BT_Cab]
		eNotRequests := requests_below(e)
		if eRequests == 1 || eNotRequests == 0{
			return 1
		} else{
			return 0
		}
		//return eRequests || !eNotRequests
		//return e.requests[e.floor][B_HallDown] || e.requests[e.floor][B_Cab] || (!requests_below(e))
	case MD_Up:
		eRequests := e.requests[e.floor][BT_HallUp] || e.requests[e.floor][BT_Cab]
		eNotRequests := requests_above(e)
		if eRequests == 1 || eNotRequests == 0{
			return 1
		} else{
			return 0
		}
		//return eRequests || !eNotRequests
		//return e.requests[e.floor][B_HallUp] || e.requests[e.floor][B_Cab] || !requests_above(e)
	case MD_Stop:
	default:
		return 1

	}
	return 1
}

func requests_clearAtCurrentFloor(e Elevator) Elevator {
	switch (e.config.clearRequestVariant) {
	case CV_All:
		for button := 0; button < N_BUTTONS; button++ {
		e.requests[e.floor][button] = 0
	}
		break

	case CV_InDirn:
		e.requests[e.floor][BT_Cab] = 0;
		switch e.dirn {
		case MD_Up:
			e.requests[e.floor][BT_HallUp] = 0
			if requests_above(e) == 0 {
				e.requests[e.floor][BT_HallDown] = 0
			}
			break
		case MD_Stop:
		default:
			e.requests[e.floor][BT_HallUp] = 0
			e.requests[e.floor][BT_HallDown] = 0
		}
		break
	}

	return e
}



