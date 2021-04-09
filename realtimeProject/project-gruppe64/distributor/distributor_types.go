package distributor

import (
	"realtimeProject/hardwareIO"
)

type ElevatorTagged struct  {
	Behaviour string `json:"behaviour"`
	Floor int `json:"floor"`
	MotorDirection string  `json:"direction"`
	CabOrders [hardwareIO.NumFloors]bool `json:"cabRequests"`
}

// https://mholt.github.io/json-to-go/
type Elevators struct{
	HallOrders [hardwareIO.NumFloors][2]bool `json:"hallRequests"`
	States map[string]ElevatorTagged `json:"states"`
}