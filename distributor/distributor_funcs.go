package distributor

import (
	"encoding/json"
	"log"
	"os/exec"
	"realtimeProject/project-gruppe64/fsm"
	"realtimeProject/project-gruppe64/hardwareIO"
	"strconv"
)

func initiateElevators() Elevators{
	elevs := Elevators{}
	elevs.HallOrders = [hardwareIO.NumFloors][2]bool{}
	var statesMap map[string]ElevatorTagged
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


func getDesignatedElevatorID(elevs Elevators) int {
	elevsEncoded, err := json.Marshal(elevs)
	if err != nil {
		log.Fatal(err)
	}
	costCmd := exec.Command("./designatorTest/hall_request_assigner.exe", "--input",  string(elevsEncoded))
	out, err := costCmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	var costMap map[string][][]bool
	err = json.Unmarshal(out, &costMap)
	if err != nil {
		log.Fatal(err)
	}
	for key, data := range costMap {
		for _, flr := range data {
			if flr[0] == true || flr[1] == true { //Hvis kalkulasjonen sier at heisen har fått ordren
				retID, err := strconv.Atoi(key)
				if err != nil {
					log.Fatal(err)
				}
				return retID
			}
		}
	}
	return -1
}