package elevator

import (
	"fmt"
	"realtimeProject/project-gruppe64/io"
)
// Enum definitions

const (
	NumButtons int = 3 // Fra elevator_io_types.h
	NumFloors int = 4 // Skrevet om fra elevator_io - må kanskje endres
)

type ElevatorBehaviour int

const (
	EB_Idle     ElevatorBehaviour = 0 // Evt skrive om til camelCase!
	EB_DoorOpen                   = 1
	EB_Moving                     = 2
)

type ClearRequestVariant int

const (
	CV_All    ClearRequestVariant = 0
	CV_InDirn                     = 1
)

// Definition of structs

type Elevator struct {
	Floor          int
	MotorDirection io.MotorDirection
	Requests       [NumFloors][NumButtons] int // Hardkoder størrelse sånn halvveis
	Behaviour      ElevatorBehaviour
	Config         struct{
		ClearRequestVariant ClearRequestVariant
		DoorOpenDurationSec float64 // Tilsvarer double i C
	}
}

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

func elevioMotorDirToString(mDirection io.MotorDirection) string { // Importere? Stor forbokstav!
	switch mDirection {
	case io.MD_Stop:
		return "MD_Stop"
	case io.MD_Up:
		return "MD_Up"
	case io.MD_Down:
		return "MD_Down"
	}
	return "MD_Undefined"
}

func checkRequests(e Elevator, number1 int, number2 int) bool {
	elevator:= Elevator{}
	elevatorRequests := elevator.Requests
	for i:=0; i< NumFloors; i++ {
		for j:=0; j < NumButtons; j++ {
			if elevatorRequests[i][j] == elevatorRequests[number1][number2] {
				return true
			}
		}
	}
	return false
}

func ElevatorPrint(es Elevator) {
	fmt.Printf("  +--------------------+\n")                                           // Sjekk om dette er riktig print-funksjon!
	fmt.Printf("  |floor = %-2d          |\n",es.Floor)                                // - i %2-d betyr bare - at teksten er left-justified (kosmetisk)
	fmt.Printf("|direction  = %-12.12s|\n", elevioMotorDirToString(es.MotorDirection)) // Måtte dele opp for at det skulle bli riktig
	fmt.Print("|behaviour = %-12.12s|\n", ebToString(es.Behaviour))                    // Hvorfor feilmelding her?
	fmt.Printf("  +--------------------+\n")
	fmt.Printf("  |  | up  | dn  | cab |\n")

	for f := NumFloors -1; f >= 0; f-- {
		fmt.Printf("  | %d", f)
		for btn := 0; btn < NumButtons; btn++{
			if (f == NumFloors-1 && btn == int(io.BT_HallUp)) || (f == 0 && btn == int(io.BT_HallDown)){ // Bedre måte å gjøre dette på enn ved konvertering til int
				fmt.Printf("|     ")
			} else {
				if checkRequests(es, f,btn) == true { // Fint om vi finner en bedre måte å skrive denne på
					fmt.Printf("|  #  ")
				} else{
					fmt.Printf("|  -  ")
				}
				// fmt.Printf(es.requests[f][btn] ? "|  #  " : "|  -  "); Originale linjen, disse operatorene finnes ikke i go
			}
		}
		fmt.Printf("|\n")
	}
	fmt.Printf("  +--------------------+\n")
}



func ElevatorUnitialized() Elevator { //
	elevator := Elevator{}
	elevator.Floor = -1
	elevator.MotorDirection = io.MD_Stop
	elevator.Behaviour = EB_Idle
	elevator.Config.ClearRequestVariant = CV_All
	elevator.Config.DoorOpenDurationSec = 3.0
	return elevator
}

