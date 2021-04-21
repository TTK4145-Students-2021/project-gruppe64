package distributor

import (
	"fmt"
	"realtimeProject/project-gruppe64/system"
)

/*
import (
	"../system"
)
*/

// GOROUTINE:

//ENDRE DISTRIBUTOR TIL Å HA GLOBAL VARIABEL ELEVATORS SOM ER ELEVATOR INFORMATION IKKE ELEVATOR TAGGED
func OrderDistributor(hallOrder <-chan system.ButtonEvent, otherElevator <-chan system.Elevator,
	ownElevator <-chan system.Elevator, shareOwnElevatorCh chan<- system.Elevator, orderThroughNet chan<- system.NetOrder,
	orderToSelf chan<- system.ButtonEvent, messageTimer chan<- system.NetOrder,
	messageTimerTimedOut <-chan system.NetOrder, orderTimer chan<- system.NetOrder,
	orderTimerTimedOut <- chan system.NetOrder, elevatorIDConnected <-chan int, elevatorIDDisconnected <-chan int){

	elevators := initiateElevators()
	elevatorsOnline := make(map[int]bool)
	distributedOrders := make(map[int][]system.NetOrder)

	for {
		select {
		case hallOrd := <-hallOrder:
			designatedID := getDesignatedElevatorID(hallOrd, elevators, elevatorsOnline)
			if designatedID == system.ElevatorID {
				orderToSelf <- hallOrd
			} else {
				sOrd := system.NetOrder{ReceivingElevatorID: designatedID, SendingElevatorID: system.ElevatorID, Order: hallOrd}
				orderThroughNet <- sOrd
				messageTimer <- sOrd
				orderTimer <- sOrd
				if distributedOrders[designatedID] != nil {
					distributedOrders[designatedID] = append(distributedOrders[designatedID], sOrd)
				} else {
					distributedOrders[designatedID] = []system.NetOrder{sOrd}
				}
			}

		case msgTimedOut := <-messageTimerTimedOut:
			fmt.Println("reassigning")
			if distributedOrders[msgTimedOut.ReceivingElevatorID] != nil {
				distributedOrders[msgTimedOut.ReceivingElevatorID] = removeOrderFromOrders(msgTimedOut, distributedOrders[msgTimedOut.ReceivingElevatorID])
			}
			designatedID := getDesignatedElevatorID(msgTimedOut.Order, elevators, elevatorsOnline)
			if designatedID == system.ElevatorID {
				orderToSelf <- msgTimedOut.Order
			} else {
				sOrd := system.NetOrder{ReceivingElevatorID: designatedID, SendingElevatorID: system.ElevatorID, Order: msgTimedOut.Order}
				orderThroughNet <- sOrd
				messageTimer <- sOrd
				if distributedOrders[designatedID] != nil {
					distributedOrders[designatedID] = append(distributedOrders[designatedID], sOrd)
				} else {
					distributedOrders[designatedID] = []system.NetOrder{sOrd}
				}
			}

		case ordTimedOut := <-orderTimerTimedOut:
			for key, dOrds := range distributedOrders {
				if key == ordTimedOut.ReceivingElevatorID {
					for _, dOrd := range dOrds {
						if dOrd == ordTimedOut {
							distributedOrders[key] = removeOrderFromOrders(ordTimedOut, distributedOrders[key])
							orderToSelf <- ordTimedOut.Order
						}
					}
				}
			}

		case e := <-ownElevator:
			shareOwnElevatorCh <- e
			system.LogElevator(e)
			elevators[system.ElevatorID] = e
			setAllHallLights(elevators)

		case e := <-otherElevator:
			if e.ID != system.ElevatorID {
				elevators[e.ID] = e
				setAllHallLights(elevators)
				if removeExecutedOrders(e, distributedOrders[e.ID]) != nil {
					distributedOrders[e.ID] = removeExecutedOrders(e, distributedOrders[e.ID])
				}
			}
		case eID := <-elevatorIDConnected:
			elevatorsOnline[eID] = true

		case eID := <-elevatorIDDisconnected:
			elevatorsOnline[eID] = false

		}
	}
}

/*
func checkForMotorStop(ordersForMotorCheck <-chan [system.NumFloors][system.NumButtons]int, motorStop chan<- bool) {
	var ordersCheck [system.NumFloors][system.NumButtons]int
	for {
		select{
		case ordForMChck := <-ordersForMotorCheck:
			ordersCheck = ordForMChck
			time.AfterFunc(5*time.Second)
		}
	}
}
*/
