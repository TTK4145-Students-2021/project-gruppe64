package elevio

import (
	"fmt"
)
 // Enum definitions
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
	//requests[N_FLOORS][N_BUTTONS]
	behaviour ElevatorBehaviour
	config struct{
		clearRequestVariant ClearRequestVariant
		doorOpenDurationSec float64 // Tilsvarer double i C
	}
}


// Implementering av funksjoner (C-fil)

/* Med channel
func ebToString(elevBehav chan<- ElevatorBehaviour) string {
	behaviour := ""
	select {
	case elevBehav<-EB_Idle:
		behaviour = "ED_Idle"
	case elevBehav<-EB_DoorOpen:
		behaviour = "EB_DoorOpen"
	case elevBehav<-EB_Moving:
		behaviour = "EB_Moving"
	default:
		behaviour = "EB_Undefined"
	}
	return behaviour // Returne her eller i hver case?
}
*/
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
}

func elevioMotorDirToString(mDirection MotorDirection) string {
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
	fmt.Println("  +--------------------+\n") // Sjekk om dette er riktig print-funksjon!
	fmt.Println("  |floor = %-2d          |\n  |directin  = %-12.12s|\n |behaviour = %-12.12s|\n",
		es.floor, elevioMotorDirToString(es.motorDirection), ebToString(es.behaviour))
	fmt.Println("  +--------------------+\n")
	fmt.Println("  |  | up  | dn  | cab |\n")
	/*for(int f = N_FLOORS-1; f >= 0; f--){
		printf("  | %d", f);
		for(int btn = 0; btn < N_BUTTONS; btn++){
			if((f == N_FLOORS-1 && btn == B_HallUp)  ||
				(f == 0 && btn == B_HallDown)
		){
				printf("|     ");
			} else {
				printf(es.requests[f][btn] ? "|  #  " : "|  -  ");
			}
		}
		printf("|\n");
	}
	printf("  +--------------------+\n");*/
}


/*
func elevatorUnitialized() Elevator {
	...
}
*/
