package fsm

import (
	"realtimeProject/project-gruppe64/hardwareIO"
)

func orderAbove(e Elevator) bool {
	for f := e.Floor+1; f < hardwareIO.NumFloors ; f++ {
		for b := 0; b < hardwareIO.NumButtons; b++ {
			if e.Orders[f][b] != 0 {
				return true
			}
		}
	}
	return false
}

func orderBelow(e Elevator) bool {
	for f := 0; f < e.Floor; f++ {
		for b := 0; b < hardwareIO.NumButtons; b++ {
			if e.Orders[f][b] != 0 {
				return true
			}
		}
	}
	return false
}

func chooseDirection(e Elevator) hardwareIO.MotorDirection{
	var returnDir hardwareIO.MotorDirection
	switch e.MotorDirection {
	case hardwareIO.MD_Up:
		if orderAbove(e){
			returnDir = hardwareIO.MD_Up
		} else if orderBelow(e) {
			returnDir = hardwareIO.MD_Down
		} else {
			returnDir = hardwareIO.MD_Stop
		}
	case hardwareIO.MD_Down:

		if orderBelow(e){
			returnDir = hardwareIO.MD_Down
		} else if orderAbove(e){
			returnDir = hardwareIO.MD_Up
		} else {
			returnDir = hardwareIO.MD_Stop
		}
	case hardwareIO.MD_Stop:
		if orderBelow(e){
			returnDir = hardwareIO.MD_Down
		} else if orderAbove(e){
			returnDir = hardwareIO.MD_Up
		}
	default:
		break
	}
	return returnDir
}

func elevatorShouldStop(e Elevator) bool {
	switch e.MotorDirection {
	case hardwareIO.MD_Down:
		if e.Orders[e.Floor][hardwareIO.BT_HallDown] != 0 || e.Orders[e.Floor][hardwareIO.BT_Cab] != 0 || !orderBelow(e){ //This gives me if eRequests == true, right?
			return true
		} else{
			return false
		}
	case hardwareIO.MD_Up:
		if e.Orders[e.Floor][hardwareIO.BT_HallUp] != 0 || e.Orders[e.Floor][hardwareIO.BT_Cab] != 0 || !orderAbove(e){ //This gives me if eRequests == true, right?
			return true
		} else{
			return false
		}

	case hardwareIO.MD_Stop:
		break
	default:
		break
	}
	return true
}

func clearOrdersAtCurrentFloor(e Elevator) Elevator{
	switch e.Config.ClearOrdersVariant {
	case CO_All: //CV:clear request variant
		for button := 0; button < hardwareIO.NumButtons; button++ { //_numButtons= 3
			e.Orders[e.Floor][button] = 0
		}

	case CO_InMotorDirection:
		e.Orders[e.Floor][hardwareIO.BT_Cab] = 0
		switch e.MotorDirection {
		case hardwareIO.MD_Up:
			e.Orders[e.Floor][hardwareIO.BT_HallUp] = 0
			if orderAbove(e) == false {
				e.Orders[e.Floor][hardwareIO.BT_HallDown] = 0
			}

		case hardwareIO.MD_Down:
			e.Orders[e.Floor][hardwareIO.BT_HallDown] = 0
			if orderBelow(e) == false {
				e.Orders[e.Floor][hardwareIO.BT_HallUp] = 0
			}

		case hardwareIO.MD_Stop:
		default:
			e.Orders[e.Floor][hardwareIO.BT_HallUp] = 0
			e.Orders[e.Floor][hardwareIO.BT_HallDown] = 0
		}
	}

	return e
}
