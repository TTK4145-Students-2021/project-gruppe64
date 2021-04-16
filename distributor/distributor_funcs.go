package distributor

import (
	"encoding/json"
	"log"
	"os/exec"
	"realtimeProject/project-gruppe64/hardwareIO"
	"realtimeProject/project-gruppe64/system"
	"strconv"
)

func initiateElevatorsTagged() system.ElevatorsTagged{
	elevs := system.ElevatorsTagged{}
	elevs.HallOrders = [system.NumFloors][2]bool{}
	statesMap :=  make(map[string]system.ElevatorTagged)
	for elevNum := 0; elevNum < system.NumElevators; elevNum ++ {
		statesMap[strconv.Itoa(elevNum)] = system.ElevatorTagged{}
	}
	elevs.States = statesMap
	return elevs
}

func getUpdatedElevatorTagged(e system.ElevatorInformation) system.ElevatorTagged{
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

func removeOfflineElevators(elevs system.ElevatorsTagged) system.ElevatorsTagged{
	retStates := make(map[string]system.ElevatorTagged)
	for elevKey, elevTagged := range elevs.States{
		if elevTagged.Behaviour != "" {
			retStates[elevKey] = elevTagged
		}
	}
	return system.ElevatorsTagged{HallOrders: elevs.HallOrders, States: retStates}
}

func getDesignatedElevatorID(elevs system.ElevatorsTagged, elevsOffline map[string]bool) int {
	costElevs := elevs

	for ID, Offline := range elevsOffline {
		if Offline {
			costElevs.States[ID] = system.ElevatorTagged{}
		}
	}

	elevsEncoded, errM := json.Marshal(removeOfflineElevators(costElevs))
	if errM != nil {
		log.Fatal(errM)
	}
	costCmd := exec.Command("./distributor/hall_request_assigner.exe", "--input",  string(elevsEncoded))
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



func checkIfOrderExecuted(elev system.ElevatorInformation, ord system.SendingOrder) bool {
	if elev.Orders[ord.Order.Floor][ord.Order.Button] == 1 {
		return false
	} else {
		return true
	}

}

func removeExecutedOrders(elev system.ElevatorInformation, distributedOrds []system.SendingOrder) []system.SendingOrder{
	var updatedDistributedOrds []system.SendingOrder
	for _, dOrds := range distributedOrds{
		if !checkIfOrderExecuted(elev, dOrds){
			updatedDistributedOrds = append(updatedDistributedOrds, dOrds)
		} else {

		}
	}
	return updatedDistributedOrds
}

func removeOrderFromOrders(orderToRemove system.SendingOrder, orders []system.SendingOrder) []system.SendingOrder {
	retOrders := orders
	for index := 0; index < len(retOrders); index++ {
		if retOrders[index] == orderToRemove {
			retOrders[index] = retOrders[len(retOrders) - 1]
			retOrders[len(retOrders) - 1] = system.SendingOrder{}
			return retOrders[:len(retOrders) - 1]
		}
	}
	return retOrders
}

func setHallButtonLights(elevInfo system.ElevatorInformation){
	for f := 0; f < system.NumFloors; f++ {
		for b := 0; b < system.NumButtons - 1; b++  {
			if elevInfo.Orders[f][b] != 0 {
				hardwareIO.SetButtonLamp(system.ButtonType(b), f, true)
			} else {
				hardwareIO.SetButtonLamp(system.ButtonType(b), f, false)
			}
		}
	}
}

