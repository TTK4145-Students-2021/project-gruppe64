package distributor

import (
	"realtimeProject/project-gruppe64/hardwareIO"
)


const (
	ElevatorID = 0 //Må endres for hver vi laster opp på
	NumElevators = 3
)


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

