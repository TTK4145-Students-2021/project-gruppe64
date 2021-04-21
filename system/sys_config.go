package system

const (
	NumFloors = 4
	NumElevators = 3

	ElevatorID =0

	LocalHost = "localhost:15657"
)



//Other configurations
const (
	//Elevator
	NumButtons = 3
	ElevatorClearOrdersVariant = COInMotorDirection
	ElevatorDoorOpenDuration = 3.0 //Sec


	CheckMotorAfterDuration = 3.0 //Sec

	//Distributor
	MaxReassignNum = 3

	//Timer
	MessageTimerDuration = 4.0 //Sec

	//Network
	NetResendNum = 10

)
