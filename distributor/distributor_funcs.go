package distributor

import (
	"encoding/json"
	"log"
	"os/exec"
	"realtimeProject/project-gruppe64/system"
	"realtimeProject/project-gruppe64/hardwareIO"
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
	case system.EB_Idle:
		behaviourString = "idle"
	case system.EB_DoorOpen:
		behaviourString = "doorOpen"
	case system.EB_Moving:
		behaviourString = "moving"
	default:
		behaviourString = ""
	}
	var motorDirString string
	switch e.MotorDirection {
	case system.MD_Up:
		motorDirString = "up"
	case system.MD_Down:
		motorDirString = "down"
	case system.MD_Stop:
		motorDirString = "stop"
	default:
		motorDirString = ""
	}
	cabOrds := [system.NumFloors]bool{}
	indexCount := 0
	for _, f := range e.Orders{
		if f[2] == 0{ //Tror dette er cab knappen i matrisen (?) litt usikker
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
	onlineElevsTagged := make(map[string]system.ElevatorTagged)
	for ID, e := range elevs {
		if ID == system.ElevatorID{
			onlineElevsTagged[strconv.Itoa(ID)] = getElevatorTagged(e)
		} else {
			if elevsOnline[ID]{
				onlineElevsTagged[strconv.Itoa(ID)] = getElevatorTagged(e)
			}
		}
	}
	elevatorsTagged := system.ElevatorsTagged{}
	elevatorsTagged.States = onlineElevsTagged
	switch ord.Button {
	case system.BT_HallUp:
		elevatorsTagged.HallOrders[ord.Floor][0] = true
	case system.BT_HallDown:
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

func checkIfOrderExecuted(e system.Elevator, ord system.NetOrder) bool {
	if e.Orders[ord.Order.Floor][ord.Order.Button] == 1 {
		return false
	} else {
		return true
	}
}

func removeExecutedOrders(e system.Elevator, distributedOrds []system.NetOrder) []system.NetOrder{
	var updatedDistributedOrds []system.NetOrder
	for _, dOrds := range distributedOrds{
		if !checkIfOrderExecuted(e, dOrds){
			updatedDistributedOrds = append(updatedDistributedOrds, dOrds)
		}
	}
	return updatedDistributedOrds
}

func removeOrderFromOrders(removeOrd system.NetOrder, ords []system.NetOrder) []system.NetOrder {
	retOrds := ords
	for index := 0; index < len(retOrds); index++ {
		if retOrds[index] == removeOrd {
			retOrds[index] = retOrds[len(retOrds) - 1]
			retOrds[len(retOrds) - 1] = system.NetOrder{}
			return retOrds[:len(retOrds) - 1]
		}
	}
	return retOrds
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


