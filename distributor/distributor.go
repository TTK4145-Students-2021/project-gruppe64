package distributor

import (
	"realtimeProject/project-gruppe64/hardwareIO"
	"strconv"
)


const (
	ElevatorID = 0 //Må endres for hver vi laster opp på
	NumElevators = 3
)


// GOROUTINE:
func OrderDistributor(hallOrder <-chan hardwareIO.ButtonEvent, elevatorInfo <-chan ElevatorInformation, sendingOrderThroughNet chan<- SendingOrder, orderToSelf chan<- hardwareIO.ButtonEvent, messageTimer chan<- SendingOrder, messageTimerTimedOut <-chan SendingOrder, orderTimer chan<- SendingOrder, orderTimerTimedOut <- chan SendingOrder){
	elevs := initiateElevators()
	var distributedOrders map[string][]SendingOrder
	for {
		select {
		case hallOrd := <-hallOrder:
			switch hallOrd.Button {
			case hardwareIO.BT_HallUp:
				elevs.HallOrders[hallOrd.Floor][0] = true
			case hardwareIO.BT_HallDown:
				elevs.HallOrders[hallOrd.Floor][1] = true
			default:
				break
			}
			designatedID := getDesignatedElevatorID(elevs)
			if designatedID == ElevatorID {
				orderToSelf <- hallOrd
			} else {
				sOrd := SendingOrder{designatedID, ElevatorID, hallOrd}
				sendingOrderThroughNet <- sOrd
				messageTimer <- sOrd
				orderTimer <- sOrd
				distributedOrders[strconv.Itoa(designatedID)] = append(distributedOrders[strconv.Itoa(designatedID)], sOrd)
			}
			switch hallOrd.Button { //sletter ordren fra her og nå
			case hardwareIO.BT_HallUp:
				elevs.HallOrders[hallOrd.Floor][0] = false
			case hardwareIO.BT_HallDown:
				elevs.HallOrders[hallOrd.Floor][1] = false
			default:
				break
			}

		case msgTimedOut := <- messageTimerTimedOut:
			orderToSelf <- msgTimedOut.order

		case ordTimedOut := <- orderTimerTimedOut:
			for key, dOrds := range distributedOrders{
				if key == strconv.Itoa(ordTimedOut.receivingElevatorID){
					for _, dOrd := range dOrds{
						if dOrd == ordTimedOut{
							orderToSelf <- ordTimedOut.order
						}
					}
				}
			}

		case elevInfo := <-elevatorInfo:
			distributedOrders[strconv.Itoa(elevInfo.ID)] = removeExecutedOrders(elevInfo, distributedOrders[strconv.Itoa(elevInfo.ID)])
			elevs.States[strconv.Itoa(elevInfo.ID)] = getUpdatedElevatorTagged(elevInfo)
		}
	}
}


