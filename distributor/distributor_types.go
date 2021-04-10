package distributor

import (
	"realtimeProject/project-gruppe64/configuration"
	"realtimeProject/project-gruppe64/fsm"
	"realtimeProject/project-gruppe64/hardwareIO"
)

//////////////////////////////////SKAL LIGGE I NETWORK/////////////////////////////////////////

type SendingOrder struct{
	receivingElevatorID int
	sendingElevatorID int
	order hardwareIO.ButtonEvent
}

type ElevatorInformation struct{
	ID     int
	Floor int
	MotorDirection hardwareIO.MotorDirection
	Orders [configuration.NumFloors][configuration.NumButtons]int
	Behaviour fsm.ElevatorBehaviour
}

///////////////////////////////////////////////////////////////////////////////////////////////

type ElevatorTagged struct  {
	Behaviour string `json:"behaviour"`
	Floor int `json:"floor"`
	MotorDirection string  `json:"direction"`
	CabOrders [configuration.NumFloors]bool `json:"cabRequests"`
}

// https://mholt.github.io/json-to-go/
type ElevatorsTagged struct{
	HallOrders [configuration.NumFloors][2]bool `json:"hallRequests"`
	States map[string]ElevatorTagged `json:"states"`
}

