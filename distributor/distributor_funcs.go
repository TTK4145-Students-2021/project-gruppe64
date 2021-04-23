package distributor

import (
	"../hardwareIO"
	"../system"
	"encoding/json"
	"log"
	"os/exec"
	"strconv"
)

func initiateElevators() map[int] system.Elevator {
	elevs := make(map[int]system.Elevator)
	for elevID := 0; elevID < system.NumElevators; elevID ++ {
		elevs[elevID] = system.Elevator{}
	}
	return elevs
}

// Translates a standard elevator struct to the struct needed for executing hall_request_assigner
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

// Returns the ID of the elevator which got the order designated from hall_request_assigner
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
	costCmd := exec.Command("./distributor/hall_request_assigner", "--input",  string(elevsEncoded))
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

// Checks if one of the elevators on net has a hall order on a floor, if so; sets hall light for elevator
func setAllHallLights(elevs map[int]system.Elevator, elevsOnline map[int]bool){
	lightsToSet := [system.NumFloors][system.NumButtons]int{}
	for _, e := range elevs {
		if elevsOnline[e.ID] || e.ID == system.ElevatorID {
			for f := 0; f < system.NumFloors; f++ {
				for b := 0; b < system.NumButtons-1; b++ {
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