package distributor

import (
	"realtimeProject/project-gruppe64/fsm"
	"realtimeProject/project-gruppe64/hardwareIO"
)



const (
	maxPlaceOrderTries = 10
)

func initiateElevatorsOnNet() [3]elevatorIdentity {
	var elevs [3]elevatorIdentity
	elevs[0] = elevatorIdentity{19001, 19998, 19999, fsm.Elevator{}}
	elevs[1] = elevatorIdentity{19002, 19999, 20002, fsm.Elevator{}}
	elevs[2] = elevatorIdentity{19003, 20000,20003, fsm.Elevator{}}
	return elevs
}

type elevatorIdentity struct {
	addressPort int
	broadcastPort int
	designatedPort int
	elevatorInfo fsm.Elevator

}

type orderTimer struct{
	order hardwareIO.ButtonEvent
}

type confirmationTimer struct{
	elevatorID int
	count int
}

func OrderDistributor(hallOrderCh <-chan hardwareIO.ButtonEvent, elevatorInformation <-chan fsm.Elevator){
	elevators := initiateElevatorsOnNet()
	for {
		select {
		case hallOrd := <-hallOrderCh:
			//ISELIN
			//Handle the hall order with designator module. Info om heisene ligger i elevators.
			//HG
			//Etter det er håndtert, så skal det sendes over network
		case elevInfo := <-elevatorInformation: //Kan være både
			testId := 0//Somehow check where this elevator is from.
			for _, elev := range elevators { //checks for the id through the elevators list
				if elev.elevatorID == testId {
					elev.elevatorInfo = elevInfo //updates the elevator info when found the right one
				}
			}
		}
	}
}

