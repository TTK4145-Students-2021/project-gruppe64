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
	elevatorConnectedCh <-chan int, elevatorDisconnectedCh <-chan int){

	elevators := initiateElevators()
	elevatorsOnline := make(map[int]bool)

	for {
		select {
		case hallOrder := <-hallOrderCh:
			designatedID := getDesignatedElevatorID(hallOrder, elevators, elevatorsOnline)
			if designatedID == system.ElevatorID {
				orderToSelfCh <- hallOrder
			} else {
				sOrd := system.NetOrder{ReceivingElevatorID: designatedID, SendingElevatorID: system.ElevatorID, Order: hallOrder, ReassignNum: 0}
				orderThroughNetCh <- sOrd
				messageTimerCh <- sOrd
			}

		case otherElevator:= <-otherElevatorCh:
			elevators[otherElevator.ID] = otherElevator
			setAllHallLights(elevators)

		case ownElevator := <-ownElevatorCh:
			shareOwnElevatorCh <- ownElevator
			system.LogElevator(ownElevator)
			elevators[system.ElevatorID] = ownElevator
			setAllHallLights(elevators)

		case messageTimerTimedOut := <-messageTimerTimedOutCh:
			if messageTimerTimedOut.ReassignNum == system.MaxReassignNum{
				orderToSelfCh <- messageTimerTimedOut.Order
			} else {
				designatedID := getDesignatedElevatorID(messageTimerTimedOut.Order, elevators, elevatorsOnline)
				if designatedID == system.ElevatorID {
					orderToSelfCh <- messageTimerTimedOut.Order
				} else {
					sOrd := system.NetOrder{ReceivingElevatorID: designatedID, SendingElevatorID: system.ElevatorID,
						Order: messageTimerTimedOut.Order, ReassignNum: messageTimerTimedOut.ReassignNum + 1}
					orderThroughNetCh <- sOrd
					messageTimerCh <- sOrd
				}
			}

		case elevatorConnected := <-elevatorConnectedCh:
			elevatorsOnline[elevatorConnected] = true

		case elevatorDisconnected := <-elevatorDisconnectedCh:
			elevatorsOnline[elevatorDisconnected] = false

		}
	}
}
