package fsm

import (
	"realtimeProject/project-gruppe64/hardwareIO"
	"realtimeProject/project-gruppe64/system"
)

/*
import (
	"../hardwareIO"
	"../system"
	"fmt"
)
*/

func setAllButtonLights(e system.Elevator){
	for f := 0; f < system.NumFloors; f++ {
		for b := 0; b < system.NumButtons; b++  {
			if e.Orders[f][b] != 0 {
				hardwareIO.SetButtonLamp(system.ButtonType(b), f, true)
			} else {
				hardwareIO.SetButtonLamp(system.ButtonType(b), f, false)
			}
		}
	}
}


func orderAbove(e system.Elevator) bool {
	for f := e.Floor+1; f < system.NumFloors; f++ {
		for b := 0; b < system.NumButtons; b++ {
			if e.Orders[f][b] != 0 {
				return true
			}
		}
	}
	return false
}

func orderBelow(e system.Elevator) bool {
	for f := 0; f < e.Floor; f++ {
		for b := 0; b < system.NumButtons; b++ {
			if e.Orders[f][b] != 0 {
				return true
			}
		}
	}
	return false
}

func chooseDirection(e system.Elevator) system.MotorDirection{
	var returnDir system.MotorDirection
	switch e.MotorDirection {
	case system.MD_Up:
		if orderAbove(e){
			returnDir = system.MD_Up
		} else if orderBelow(e) {
			returnDir = system.MD_Down
		} else {
			returnDir = system.MD_Stop
		}
	case system.MD_Down:

		if orderBelow(e){
			returnDir = system.MD_Down
		} else if orderAbove(e){
			returnDir = system.MD_Up
		} else {
			returnDir = system.MD_Stop
		}
	case system.MD_Stop:
		if orderBelow(e){
			returnDir = system.MD_Down
		} else if orderAbove(e){
			returnDir = system.MD_Up
		}
	default:
		break
	}
	return returnDir
}

func elevatorShouldStop(e system.Elevator) bool {
	switch e.MotorDirection {
	case system.MD_Down:
		if e.Orders[e.Floor][system.BT_HallDown] != 0 || e.Orders[e.Floor][system.BT_Cab] != 0 || !orderBelow(e){ //This gives me if eRequests == true, right?
			return true
		} else{
			return false
		}
	case system.MD_Up:
		if e.Orders[e.Floor][system.BT_HallUp] != 0 || e.Orders[e.Floor][system.BT_Cab] != 0 || !orderAbove(e){ //This gives me if eRequests == true, right?
			return true
		} else{
			return false
		}

	case system.MD_Stop:
		break
	default:
		break
	}
	return true
}

func clearOrdersAtCurrentFloor(e system.Elevator) system.Elevator{
	switch e.Config.ClearOrdersVariant {
	case system.CO_All: //CV:clear request variant
		for button := 0; button < system.NumButtons; button++ { //_numButtons= 3
			e.Orders[e.Floor][button] = 0
		}

	case system.CO_InMotorDirection:
		e.Orders[e.Floor][system.BT_Cab] = 0
		switch e.MotorDirection {
		case system.MD_Up:
			e.Orders[e.Floor][system.BT_HallUp] = 0
			if orderAbove(e) == false {
				e.Orders[e.Floor][system.BT_HallDown] = 0
			}

		case system.MD_Down:
			e.Orders[e.Floor][system.BT_HallDown] = 0
			if orderBelow(e) == false {
				e.Orders[e.Floor][system.BT_HallUp] = 0
			}

		case system.MD_Stop:
		default:
			e.Orders[e.Floor][system.BT_HallUp] = 0
			e.Orders[e.Floor][system.BT_HallDown] = 0
		}
	}
	return e
}

