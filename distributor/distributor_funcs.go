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

func initiateElevators() elevators{
	elevs := elevators{}
	tempElevs := [NumElevators]elevatorTagged{}
	for count := 0; count < NumElevators; count ++ {
		tempElevs[count-1] = elevatorTagged{} //Trenger json flag!
	}
	elevs.states = tempElevs
	return elevators{}
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