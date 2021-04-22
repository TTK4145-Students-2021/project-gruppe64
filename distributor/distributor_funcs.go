package distributor

import (
	"encoding/json"
	"log"
	"os/exec"
	"realtimeProject/project-gruppe64/hardwareIO"
	"realtimeProject/project-gruppe64/system"
	"strconv"
)
/*
import (
	"encoding/json"
	"log"
	"os/exec"
	"realtimeProject/project-gruppe64/hardwareIO"
	"realtimeProject/project-gruppe64/system"
	"strconv"
)
*/

func initiateElevators() map[int] system.Elevator {
	elevs := make(map[int]system.Elevator)
	for elevID := 0; elevID < system.NumElevators; elevID ++ {
		elevs[elevID] = system.Elevator{}
	}
	return elevs
}

func getElevatorTagged(e system.Elevator) system.ElevatorTagged{
	var behaviourString string
	switch e.Behaviour {
	case system.EBIdle:
		behaviourString = "idle"
	case system.EBDoorOpen:
		behaviourString = "doorOpen"
	case system.EBMoving:
		behaviourString = "moving"
	default:
		behaviourString = ""
	}
	var motorDirString string
	switch e.MotorDirection {
	case system.MDUp:
		motorDirString = "up"
	case system.MDDown:
		motorDirString = "down"
	case system.MDStop:
		motorDirString = "stop"
	default:
		motorDirString = ""
	}
	cabOrds := [system.NumFloors]bool{}
	indexCount := 0
	for _, f := range e.Orders{
		if f[2] == 0{
			cabOrds[indexCount] = false
		} else {
			cabOrds[indexCount] = true
		}
		indexCount += 1
	}
	return system.ElevatorTagged{Behaviour: behaviourString, Floor: e.Floor, MotorDirection: motorDirString,
		CabOrders: cabOrds}
}

func getDesignatedElevatorID(ord system.ButtonEvent, elevs map[int]system.Elevator, elevsOnline map[int]bool) int {
	availableElevsTagged := make(map[string]system.ElevatorTagged)
	for ID, e := range elevs {
		if ID == system.ElevatorID && !e.MotorError{
			availableElevsTagged[strconv.Itoa(ID)] = getElevatorTagged(e)
		} else {
			if elevsOnline[ID] && !e.MotorError{
				availableElevsTagged[strconv.Itoa(ID)] = getElevatorTagged(e)
			}
		}
	}
	if len(availableElevsTagged) == 0 {
		availableElevsTagged[strconv.Itoa(system.ElevatorID)] = getElevatorTagged(elevs[system.ElevatorID])
	}
	elevatorsTagged := system.ElevatorsTagged{}
	elevatorsTagged.States = availableElevsTagged
	switch ord.Button {
	case system.BTHallUp:
		elevatorsTagged.HallOrders[ord.Floor][0] = true
	case system.BTHallDown:
		elevatorsTagged.HallOrders[ord.Floor][1] = true
	default:
		break
	}
	elevsEncoded, errM := json.Marshal(elevatorsTagged)
	if errM != nil {
		log.Fatal(errM)
	}
	costCmd := exec.Command("./distributor/hall_request_assigner.exe", "--input",  string(elevsEncoded))
	// LINUX:
	// costCmd := exec.Command("./distributor/hall_request_assigner", "--input",  string(elevsEncoded))
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
			if flr[0] == true || flr[1] == true {
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

func setAllHallLights(elevs map[int]system.Elevator){
	lightsToSet := [system.NumFloors][system.NumButtons]int{}
	for _, e := range elevs {
		for f := 0; f < system.NumFloors; f++ {
			for b := 0; b < system.NumButtons - 1; b++ {
				if e.Orders[f][b] == 0 {
					if lightsToSet[f][b] != 1 {
						lightsToSet[f][b] = 0
					}
				} else {
					lightsToSet[f][b] = 1
				}
			}
		}
	}
	for f := 0; f < system.NumFloors; f++ {
		for b := 0; b < system.NumButtons -1; b++ {
			if lightsToSet[f][b] == 0 {
				hardwareIO.SetButtonLamp(system.ButtonType(b), f, false)
			} else {
				hardwareIO.SetButtonLamp(system.ButtonType(b), f, true)
			}
		}
	}
}

func removeOrder(orders []system.NetOrder, i int)  []system.NetOrder {
	orders[i] = orders[len(orders)-1]
	return orders[:len(orders) - 1]
}

func checkIfOnlyOneOnline(elevsOnline map[int]bool) bool {
	for ID, online := range elevsOnline{
		if ID != system.ElevatorID{
			if online {
				return false
			}
		}
	}
	return true
}