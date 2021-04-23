package fsm

import (
	"../hardwareIO"
	"../system"
)

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
	case system.MDUp:
		if orderAbove(e){
			returnDir = system.MDUp
		} else if orderBelow(e) {
			returnDir = system.MDDown
		} else {
			returnDir = system.MDStop
		}
	case system.MDDown:

		if orderBelow(e){
			returnDir = system.MDDown
		} else if orderAbove(e){
			returnDir = system.MDUp
		} else {
			returnDir = system.MDStop
		}
	case system.MDStop:
		if orderBelow(e){
			returnDir = system.MDDown
		} else if orderAbove(e){
			returnDir = system.MDUp
		}
	default:
		break
	}
	return returnDir
}

func elevatorShouldStop(e system.Elevator) bool {
	switch e.MotorDirection {
	case system.MDDown:
		if e.Orders[e.Floor][system.BTHallDown] != 0 ||
			e.Orders[e.Floor][system.BTCab] != 0 || !orderBelow(e){
			return true
		} else{
			return false
		}

	case system.MDUp:
		if e.Orders[e.Floor][system.BTHallUp] != 0 ||
			e.Orders[e.Floor][system.BTCab] != 0 || !orderAbove(e){
			return true
		} else{
			return false
		}

	case system.MDStop:
		break

	default:
		break
	}
	return true
}

func clearOrdersAtCurrentFloor(e system.Elevator) system.Elevator{
	switch e.Config.ClearOrdersVariant {
	case system.COAll:
		for button := 0; button < system.NumButtons; button++ {
			e.Orders[e.Floor][button] = 0
		}
	case system.COInMotorDirection:
		e.Orders[e.Floor][system.BTCab] = 0
		switch e.MotorDirection {
		case system.MDUp:
			e.Orders[e.Floor][system.BTHallUp] = 0
			if orderAbove(e) == false {
				e.Orders[e.Floor][system.BTHallDown] = 0
			}

		case system.MDDown:
			e.Orders[e.Floor][system.BTHallDown] = 0
			if orderBelow(e) == false {
				e.Orders[e.Floor][system.BTHallUp] = 0
			}

		case system.MDStop:
		default:
			e.Orders[e.Floor][system.BTHallUp] = 0
			e.Orders[e.Floor][system.BTHallDown] = 0
		}
	}
	return e
}

// Sets the cab- and hall lights according to the elevator's own orders
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
