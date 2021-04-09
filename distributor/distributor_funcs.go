package distributor

import (
	"encoding/json"
	"log"
	"os/exec"
	"realtimeProject/project-gruppe64/fsm"
	"realtimeProject/project-gruppe64/hardwareIO"
	"strconv"
)

func initiateElevatorsTagged() ElevatorsTagged{
	elevs := ElevatorsTagged{}
	elevs.HallOrders = [hardwareIO.NumFloors][2]bool{}
	statesMap :=  make(map[string]ElevatorTagged)
	for elevNum := 0; elevNum < NumElevators; elevNum ++ {
		statesMap[strconv.Itoa(elevNum)] = ElevatorTagged{}
	}
	elevs.States = statesMap
	return elevs
}


func getUpdatedElevatorTagged(e ElevatorInformation) ElevatorTagged{
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
	return ElevatorTagged{behaviourString, e.Floor, motorDirString, cabOrds}
}


func getDesignatedElevatorID(elevs ElevatorInformation) int {
	elevsEncoded, errM := json.Marshal(elevs)
	if errM != nil {
		log.Fatal(errM)
	}
	costCmd := exec.Command("./designator/hall_request_assigner.exe", "--input",  string(elevsEncoded))
	out, errO := costCmd.Output()
	if errO != nil {
		log.Fatal(errO)
	}

	var costMap map[string][][]bool
	errU := json.Unmarshal(out, &costMap)
	if errU != nil {
		log.Fatal(errU)
	}
	for key, data := range costMap {

		for _, flr := range data {
			if flr[0] == true || flr[1] == true { //Hvis kalkulasjonen sier at heisen har fått ordren
				retID, errK := strconv.Atoi(key)
				if errK != nil {
					log.Fatal(errK)
				}
				return retID
			}
		}
	}
	return -1
}



func checkIfOrderExecuted(elev ElevatorInformation, ord SendingOrder) bool {
	if elev.Orders[ord.order.Floor][ord.order.Button] == 1 {
		return false
	} else {
		return true
	}

}

func removeExecutedOrders(elev ElevatorInformation, distributedOrds []SendingOrder) []SendingOrder{
	var updatedDistributedOrds []SendingOrder
	for _, dOrds := range distributedOrds{
		if !checkIfOrderExecuted(elev, dOrds){
			updatedDistributedOrds = append(updatedDistributedOrds, dOrds)
		}
	}
	return updatedDistributedOrds
}

