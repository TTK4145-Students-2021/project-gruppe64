package distributor

import (
	"realtimeProject/project-gruppe64/system"
)

/*
import (
	"../system"
)
*/

func OrderDistributor(hallOrderCh <-chan system.ButtonEvent, otherElevatorCh <-chan system.Elevator,
	ownElevatorCh <-chan system.Elevator, shareOwnElevatorCh chan<- system.Elevator,
	orderThroughNetCh chan<- system.NetOrder, orderToSelfCh chan<- system.ButtonEvent,
	messageTimerCh chan<- system.NetOrder, messageTimerTimedOutCh <-chan system.NetOrder,
	orderTimerCh chan<- system.NetOrder, orderTimerTimedOutCh <- chan system.NetOrder,
	elevatorConnectedCh <-chan int, elevatorDisconnectedCh <-chan int){

	elevators := initiateElevators()
	elevatorsOnline := make(map[int]bool)
	distributedOrders := make(map[int][]system.NetOrder)

	for {
		select {
		case hallOrder := <-hallOrderCh:
			designatedID := getDesignatedElevatorID(hallOrder, elevators, elevatorsOnline)
			if designatedID == system.ElevatorID {
				orderToSelfCh <- hallOrder
			} else {
				sOrd := system.NetOrder{ReceivingElevatorID: designatedID, SendingElevatorID: system.ElevatorID, Order: hallOrder}
				orderThroughNetCh <- sOrd
				messageTimerCh <- sOrd
				orderTimerCh <- sOrd
				if distributedOrders[designatedID] != nil {
					distributedOrders[designatedID] = append(distributedOrders[designatedID], sOrd)
				} else {
					distributedOrders[designatedID] = []system.NetOrder{sOrd}
				}
			}

		case otherElevator:= <-otherElevatorCh:
			elevators[otherElevator.ID] = otherElevator
			setAllHallLights(elevators)
			if removeExecutedOrders(otherElevator, distributedOrders[otherElevator.ID]) != nil {
				distributedOrders[otherElevator.ID] = removeExecutedOrders(otherElevator,
					distributedOrders[otherElevator.ID])
			}

		case ownElevator := <-ownElevatorCh:
			shareOwnElevatorCh <- ownElevator
			system.LogElevator(ownElevator)
			elevators[system.ElevatorID] = ownElevator
			setAllHallLights(elevators)

		case messageTimerTimedOut := <-messageTimerTimedOutCh:
			if distributedOrders[messageTimerTimedOut.ReceivingElevatorID] != nil {
				distributedOrders[messageTimerTimedOut.ReceivingElevatorID] = removeOrderFromOrders(
					messageTimerTimedOut, distributedOrders[messageTimerTimedOut.ReceivingElevatorID])
			}
			designatedID := getDesignatedElevatorID(messageTimerTimedOut.Order, elevators, elevatorsOnline)
			if designatedID == system.ElevatorID {
				orderToSelfCh <- messageTimerTimedOut.Order
			} else {
				sOrd := system.NetOrder{ReceivingElevatorID: designatedID, SendingElevatorID: system.ElevatorID,
					Order: messageTimerTimedOut.Order}
				orderThroughNetCh <- sOrd
				messageTimerCh <- sOrd
				if distributedOrders[designatedID] != nil {
					distributedOrders[designatedID] = append(distributedOrders[designatedID], sOrd)
				} else {
					distributedOrders[designatedID] = []system.NetOrder{sOrd}
				}
			}

		case orderTimerTimedOut := <-orderTimerTimedOutCh:
			for key, dOrds := range distributedOrders {
				if key == orderTimerTimedOut.ReceivingElevatorID {
					for _, dOrd := range dOrds {
						if dOrd == orderTimerTimedOut {
							distributedOrders[key] = removeOrderFromOrders(orderTimerTimedOut, distributedOrders[key])
							orderToSelfCh <- orderTimerTimedOut.Order
						}
					}
				}
			}

		case elevatorConnected := <-elevatorConnectedCh:
			elevatorsOnline[elevatorConnected] = true

		case elevatorDisconnected := <-elevatorDisconnectedCh:
			elevatorsOnline[elevatorDisconnected] = false

		}
	}
}

