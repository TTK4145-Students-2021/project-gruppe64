package distributor

import (
	"realtimeProject/project-gruppe64/configuration"
)



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

