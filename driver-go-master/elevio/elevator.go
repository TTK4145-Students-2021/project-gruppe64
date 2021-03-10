package elevio

import (
	"fmt"
)
 // Enum definitions

var _numButtons int = 3 // Fra elevator_io_types.h

type ElevatorBehaviour int

const (
	EB_Idle ElevatorBehaviour 	= 0		// Evt skrive om til camelCase!
	EB_DoorOpen 				= 1
	EB_Moving					= 2
)

type ClearRequestVariant int

const (
	CV_All 		ClearRequestVariant	= 0
	CV_InDirn 						= 1
)

// Definition of structs

type Elevator struct {
	floor int
	motorDirection MotorDirection
	requests[_numFloors][_numButtons] // Hvor kommer denne fra? Får feil
	behaviour ElevatorBehaviour
	config struct{
		clearRequestVariant ClearRequestVariant
		doorOpenDurationSec float64 // Tilsvarer double i C
	}
}


// Implementering av funksjoner (C-fil)

func ebToString(elevBehav chan<- ElevatorBehaviour) string { // Trenger vi å importere? Stor forbokstav!
	behav := ""
	select {
	case elevBehav<-EB_Idle:
		behav = "ED_Idle"
	case elevBehav<-EB_DoorOpen:
		behav = "EB_DoorOpen"
	case elevBehav<-EB_Moving:
		behav = "EB_Moving"
	default:
		behav = "EB_Undefined"
	}
	return behav // Returne her eller i hver case?
}

/* Uten channels
func ebToString(elevBehav ElevatorBehaviour) string {
	switch elevBehav {
	case EB_Idle:
		return "ED_Idle"
	case EB_DoorOpen:
		return "EB_DoorOpen"
	case EB_Moving:
		return "EB_Moving"
	}
	return "EB_Undefined"
}*/

func elevioMotorDirToString(mDirection MotorDirection) string { // Importere? Stor forbokstav!
	switch mDirection {
	case MD_Stop:
		return "MD_Stop"
	case MD_Up:
		return "MD_Up"
	case MD_Down:
		return "MD_Down"
	}
	return "MD_Undefined"
}

func elevatorPrint(es Elevator) {
	fmt.Printf("  +--------------------+\n") // Sjekk om dette er riktig print-funksjon!
	fmt.Printf("  |floor = %-2d          |\n",es.floor) // - i %2-d betyr bare - at teksten er left-justified (kosmetisk)
	fmt.Printf("|direction  = %-12.12s|\n", elevioMotorDirToString(es.motorDirection)) // Måtte dele opp for at det skulle bli riktig
	fmt.Print("|behaviour = %-12.12s|\n", ebToString(es.behaviour) // Hvorfor feilmelding her?
	fmt.Printf("  +--------------------+\n")
	fmt.Printf("  |  | up  | dn  | cab |\n")

	for f := _numFloors-1; f >= 0; f-- {
		fmt.Printf("  | %d", f)
		for btn := 0; btn < _numButtons; btn++{
			if (f == _numFloors-1 && btn == BT_HallUp) || (f == 0 && btn == BT_HallDown){ // Hvorfor blir BT_HallUp feil når ButtonType er int?
				fmt.Printf("|     ")
			} else {
				if es.requests[f][btn] == true { // Fint om vi finner en bedre måte å skrive denne på
					fmt.Printf("|  #  ")
				} else{
					fmt.Printf("|  -  ")
				}
				// fmt.Printf(es.requests[f][btn] ? "|  #  " : "|  -  "); Originale linjen, disse operatorene finnes ikke i go
			}
		}
		fmt.Printf("|\n");
	}
	fmt.Printf("  +--------------------+\n");
}



func elevatorUnitialized() Elevator { // Helt klart feil, må finne ut hvordan denne løses
	Elevator.floor = -1
	Elevator.motorDirection = MD_Stop
	Elevator.behaviour = EB_Idle
	Elevator.config.clerRequestVariant = CV_All
	Elevator.config.doorOpenDurationSec = 3.0
	return Elevator.floor, Elevator.motorDirection, Elevator.behaviour,
	Elevator.config.clerRequestVariant, Elevator.config.doorOpenDurationSec
}
