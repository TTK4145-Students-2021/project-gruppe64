package distributor

import (
	"realtimeProject/project-gruppe64/fsm"
	"realtimeProject/project-gruppe64/hardwareIO"
)

//////////////////////////////////SKAL LIGGE I NETWORK/////////////////////////////////////////

type OrderToSend struct{
	receivingElevatorID int
	sendingElevatorID int
	order hardwareIO.ButtonEvent
}

type ElevatorInformation struct{
	ID     int
	Floor int
	MotorDirection hardwareIO.MotorDirection
	Orders [hardwareIO.NumFloors][hardwareIO.NumButtons]int
	Behaviour fsm.ElevatorBehaviour
}


///////////////////////////////////////////////////////////////////////////////////////////////

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
type elevatorCostCalculatedOrders struct {
	elevatorCostOrders [hardwareIO.NumFloors][2]bool //Need json tags maybe (?)
}
type costCalculatedOrders struct {
	allCostOrders [NumElevators]elevatorCostCalculatedOrders
}
