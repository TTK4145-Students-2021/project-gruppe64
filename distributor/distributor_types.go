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

type elevatorTagged struct  {
	behaviour string `json:"behaviour"`
	floor int `json:"floor"`
	motorDirection string  `json:"direction"`
	cabOrders [hardwareIO.NumFloors]bool `json:"cabRequests"`
}

// https://mholt.github.io/json-to-go/
type elevators struct{
	hallOrders [hardwareIO.NumFloors][2]bool `json:"hallRequests"`
	states [NumElevators]elevatorTagged `json:"states"`
}

type elevatorCostCalculatedOrders struct {
	elevatorCostOrders [hardwareIO.NumFloors][2]bool //Need json tags maybe (?)
}
type costCalculatedOrders struct {
	allCostOrders [NumElevators]elevatorCostCalculatedOrders
}
