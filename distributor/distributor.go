package distributor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"realtimeProject/project-gruppe64/fsm"
	"realtimeProject/project-gruppe64/hardwareIO"
	"strings"
)


const (
	ElevatorID = 0 //Må endres for hver vi laster opp på
	NumElevators = 3
)


func initiateElevators() elevators{
	elevs := elevators{}
	tempElevs := [NumElevators]elevatorTagged{}
	for count := 0; count < NumElevators; count ++ {
		tempElevs[count-1] = elevatorTagged{} //Trenger json flag!
	}
	elevs.states = tempElevs
	return elevators{}
}

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

func getUpdatedElevatorTagged(e ElevatorInformation) elevatorTagged{
	var behaviourString string
	switch e.Behaviour {
	case fsm.EB_Idle:
		behaviourString = "idle"
	case fsm.EB_DoorOpen:
		behaviourString = "doorOpen"
	case fsm.EB_Moving:
		behaviourString = "moving"
	default:
		behaviourString = ""
	}
	var motorDirString string
	switch e.MotorDirection {
	case hardwareIO.MD_Up:
		motorDirString = "up"
	case hardwareIO.MD_Down:
		motorDirString = "down"
	case hardwareIO.MD_Stop:
		motorDirString = "stop"
	default:
		motorDirString = ""
	}
	cabOrds := [hardwareIO.NumFloors]bool{}
	indexCount := 0
	for _, f := range e.Orders{
		if f[2] == 0{ //Tror dette er cab knappen i matrisen (?) litt usikker
			cabOrds[indexCount] = false
		} else {
			cabOrds[indexCount] = true
		}
		indexCount += 1
	}
	return elevatorTagged{behaviourString, e.Floor, motorDirString, cabOrds}
}


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

func getDesignatedElevatorID(elevs elevators) int { //HER BARE PRØVER JEG MEG FREM ALTSÅ, aner ikke om funker.
	// Kun tatt hensyn til at nr 1 er heis 0 osv osv., ikke om en er av nett
	elevsEncoded, _ := json.Marshal(elevs)
	costCmd := exec.Command("cmd", "/C", "start", "powershell", "realtimeProject/project-gruppe64/designator/hall_request_assigner")
	fmt.Println("realtimeProject/project-gruppe64/designator/hall_request_assigner --input '" + string(elevsEncoded) + "'") //printe det vi prøver å execute
	//https://golang.org/pkg/os/exec/#Command
	costCmd.Stdin = strings.NewReader( "--input '" + string(elevsEncoded) + "'")
	var out bytes.Buffer
	costCmd.Stdout = &out
	err := costCmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	var costOut costCalculatedOrders
	err = json.Unmarshal(out.Bytes(), &costOut)
	if err != nil {
		log.Fatal(err)
	}
	id := -1
	for _, elevOrds := range costOut.allCostOrders {
		id += 1
		for _, ord := range elevOrds.elevatorCostOrders {
			if ord[0] == true || ord[1] == true { //Hvis kalkulasjonen sier at heisen har fått ordren
				return id
			}
		}
	}
	return id
}

// GOROUTINE:
func OrderDistributor(hallOrder <-chan hardwareIO.ButtonEvent, elevatorInfo <-chan ElevatorInformation, sendOrder chan<- OrderToSend, orderToSelf chan<- hardwareIO.ButtonEvent){
	elevs := initiateElevators()
	for {
		select {
		case hallOrd := <-hallOrder:
			switch hallOrd.Button {
			case hardwareIO.BT_HallUp:
				elevs.hallOrders[hallOrd.Floor][0] = true
			case hardwareIO.BT_HallDown:
				elevs.hallOrders[hallOrd.Floor][1] = true
			default:
				break
			}
			designatedID := getDesignatedElevatorID(elevs)
			if designatedID == ElevatorID {
				orderToSelf <- hallOrd
			} else {
				sendOrder <- OrderToSend{designatedID, ElevatorID, hallOrd}
				//her må timer for Plassert melding startes. Timer må også ha info om ordren.
				//Når plassert kommer -> timer for ordren i seg selv må startes.
				//Om ikke kommer -> ordren plasseres til en selv.

				//Når timeren for ordren i seg selv går ut sjekkes det om den er slettet fra structen til den heisen.
				//Om ordren fortsatt er der; ta den selv.
			}

			switch hallOrd.Button {
			case hardwareIO.BT_HallUp:
				elevs.hallOrders[hallOrd.Floor][0] = false
			case hardwareIO.BT_HallDown:
				elevs.hallOrders[hallOrd.Floor][1] = false
			default:
				break
			}


		case elevInfo := <-elevatorInfo:
			elevs.states[elevInfo.ID] = getUpdatedElevatorTagged(elevInfo)
			}
		}
	}
}

