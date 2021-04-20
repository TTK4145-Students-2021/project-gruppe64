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
	elevators := make(map[int]system.Elevator)
	for elevID := 0; elevID < system.NumElevators; elevID ++ {
		elevators[elevID] = system.Elevator{}
	}
	return elevators
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
	return system.ElevatorTagged{Behaviour: behaviourString, Floor: e.Floor, MotorDirection: motorDirString, CabOrders: cabOrds}
}


func getDesignatedElevatorID(hallOrder system.ButtonEvent, elevators map[int]system.Elevator, elevatorsOnline map[int]bool) int {
	onlineElevatorsTagged := make(map[string]system.ElevatorTagged)
	for ID, elevator := range elevators {
		if ID == system.ElevatorID{
			onlineElevatorsTagged[strconv.Itoa(ID)] = getElevatorTagged(elevator)
		} else {
			if elevatorsOnline[ID]{
				onlineElevatorsTagged[strconv.Itoa(ID)] = getElevatorTagged(elevator)
			}
		}
	}
	elevatorsTagged := system.ElevatorsTagged{}
	elevatorsTagged.States = onlineElevatorsTagged
	switch hallOrder.Button {
	case system.BT_HallUp:
		elevatorsTagged.HallOrders[hallOrder.Floor][0] = true
	case system.BT_HallDown:
		elevatorsTagged.HallOrders[hallOrder.Floor][1] = true
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





func checkIfOrderExecuted(elev system.Elevator, ord system.NetOrder) bool {
	if elev.Orders[ord.Order.Floor][ord.Order.Button] == 1 {
		return false
	} else {
		return true
	}

}

func removeExecutedOrders(elev system.Elevator, distributedOrds []system.NetOrder) []system.NetOrder{
	var updatedDistributedOrds []system.NetOrder
	for _, dOrds := range distributedOrds{
		if !checkIfOrderExecuted(elev, dOrds){
			updatedDistributedOrds = append(updatedDistributedOrds, dOrds)
		} else {

		}
	}
	return updatedDistributedOrds
}

func removeOrderFromOrders(orderToRemove system.NetOrder, orders []system.NetOrder) []system.NetOrder {
	retOrders := orders
	for index := 0; index < len(retOrders); index++ {
		if retOrders[index] == orderToRemove {
			retOrders[index] = retOrders[len(retOrders) - 1]
			retOrders[len(retOrders) - 1] = system.NetOrder{}
			return retOrders[:len(retOrders) - 1]
		}
	}
	return retOrders
}



func setAllHallLights(elevators map[int]system.Elevator){
	lightsToSet := [system.NumFloors][system.NumButtons]int{}
	for _, elevator := range elevators {
		for f := 0; f < system.NumFloors; f++ {
			for b := 0; b < system.NumButtons - 1; b++ {
				if elevator.Orders[f][b] == 0 {
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


