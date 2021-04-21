package system

const (
	NumFloors = 4
	NumElevators = 3

	ElevatorID =1

	LocalHost = "localhost:15661"
)



//Other configurations
const (
	//Elevator
	NumButtons = 3
	ElevatorClearOrdersVariant = CO_InMotorDirection
	ElevatorDoorOpenDuration = 3.0 //Sec


	CheckMotorAfterDuration = 3.0 //Sec

	//Distributor
	MaxReassignNum = 3

	//Timer
	MessageTimerDuration = 4.0 //Sec

	//Network
	NetResendNum = 10

)
