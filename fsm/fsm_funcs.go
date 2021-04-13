package fsm

import (
	"fmt"
	"realtimeProject/project-gruppe64/hardwareIO"
	"realtimeProject/project-gruppe64/system"
)

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
	default:
		for button := 0; button < system.NumButtons; button++ { //_numButtons= 3
			e.Orders[e.Floor][button] = 0
		}
	}
	return e
}


/////////////////////////////////////////////////////////////////////////////////////
////////////////////////////Maybe not necessary/////////////////////////////////////

func elevatorBehaviourToString(eB system.ElevatorBehaviour) string {
	switch eB {
	case system.EB_Idle:
		return "ED_Idle"
	case system.EB_DoorOpen:
		return "EB_DoorOpen"
	case system.EB_Moving:
		return "EB_Moving"
	}
	return "EB_Undefined"
}

func motorDirectionToString(mD system.MotorDirection) string { // Importere? Stor forbokstav!
	switch mD {
	case system.MD_Stop:
		return "MD_Stop"
	case system.MD_Up:
		return "MD_Up"
	case system.MD_Down:
		return "MD_Down"
	}
	return "MD_Undefined"
}

func printElevator(e system.Elevator) {
	fmt.Printf("  +--------------------+\n")                                             // Sjekk om dette er riktig print-funksjon!
	fmt.Printf("  |floor = %-2d          |\n", e.Floor)                               // - i %2-d betyr bare - at teksten er left-justified (kosmetisk)
	fmt.Printf("|direction  = %-12.12s|\n", motorDirectionToString(e.MotorDirection)) // Måtte dele opp for at det skulle bli riktig
	fmt.Printf("|behaviour = %s|\n", elevatorBehaviourToString(e.Behaviour))                         // Hvorfor feilmelding her?
	fmt.Printf("  +--------------------+\n")
	fmt.Printf("  |  | up  | dn  | cab |\n")
	for f := system.NumFloors - 1; f >= 0; f-- {
		fmt.Printf("\n  | %d", f)
		for b := 0; b < system.NumButtons; b++ {
			if (f == system.NumFloors-1 && b == int(system.BT_HallUp)) || (f == 0 && b == system.BT_HallDown) {
				fmt.Printf("|     ")
			} else {
				if e.Orders[f][b] != 0 {
					fmt.Printf("|  #  ")
				} else {
					fmt.Printf("|  -  ")
				}
			}
		}
	}
}