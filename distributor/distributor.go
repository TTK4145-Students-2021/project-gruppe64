package distributor

import (
	"encoding/json"
	"fmt"
	"net/mail"
	"realtimeProject/project-gruppe64/fsm"
	"realtimeProject/project-gruppe64/hardwareIO"
)



const (
	maxPlaceOrderTries = 10
)

func initiateElevatorsOnNet() [3]elevatorIdentity {
	elevOneNums := []int{19001, 19998, 20001} // nr. 1 for IP-adresses, 2 for elevstructs to this, 3 for orders to this
	elevTwoNums := []int{19002, 19999, 20002}
	elevThreeNums := []int{19003, 20000, 20003}

	var elevs [3]elevatorIdentity
	elevs[0] = elevatorIdentity{19001, 19998, 19999, fsm.Elevator{}}
	elevs[1] = elevatorIdentity{19002, 19999, 20002, fsm.Elevator{}}
	elevs[2] = elevatorIdentity{19003, 20000,20003, fsm.Elevator{}}
	return elevs
}

func updateElevatorStruct(elev fsm.Elevator) {
	
}

// https://mholt.github.io/json-to-go/
type elevators struct{
	hallOrders [][]bool
	states struct{
		elevatorOne struct{
			behaviour string `json:"behaviour"`
			floor int `json:"floor"`
			motorDirection string  `json:"direction"`
			orders []bool `json: "cabRequests"`
		} `json: "one"`
		elevatorTwo struct{
			behaviour string `json:"behaviour"`
			floor int `json:"floor"`
			motorDirection string  `json:"direction"`
			orders []bool `json: "cabRequests"`
		} `json: "two"`
		elevatorThree struct{
			behaviour string `json:"behaviour"`
			floor int `json:"floor"`
			motorDirection string  `json:"direction"`
			orders []bool `json: "cabRequests"`
		}`json: "three"`
	}
} //Flagg må være samme for at skal funke

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

// type orderToSend{
// receivingElevatorID int
// sendingElevatorID int
// order int }

func OrderDistributor(hallOrderCh <-chan hardwareIO.ButtonEvent, elevatorInformation <-chan fsm.Elevator){
	for {
		select {
		case hallOrd := <-hallOrderCh:

			elevOne, _ := json.Marshal(elevators[0].elevatorInfo)
			elevTwo, _ := json.Marshal(elevators[1].elevatorInfo)
			elevThree, _ := json.Marshal(elevators[2].elevatorInfo)
			fmt.Println("realtimeProject/project-gruppe64/designator/hall_request_assigner --input ")
			// Send til network chan for ordre: type orderToSend

		case elevInfo := <-elevatorInformation: //Kommer som elevator struct, gjør om til json objekt for lagring (se struct for sending til designator)
			testId := 0//Somehow check where this elevator is from.
			for _, elev := range elevators { //checks for the id through the elevators list
				if elev.elevatorID == testId {
					elev.elevatorInfo = elevInfo //updates the elevator info when found the right one
				}
			}
		}
	}
}

