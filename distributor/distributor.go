package distributor

import (
	"fmt"
	"realtimeProject/project-gruppe64/system"
)

/*
import (
	"../system"
	"fmt"
)
*/

func OrderDistributor(hallOrderCh <-chan system.ButtonEvent, otherElevatorCh <-chan system.Elevator,
	ownElevatorCh <-chan system.Elevator, shareOwnElevatorCh chan<- system.Elevator,
	orderThroughNetCh chan<- system.NetOrder, orderToSelfCh chan<- system.ButtonEvent,
	messageTimerCh chan<- system.NetOrder, messageTimerTimedOutCh <-chan system.NetOrder,
	orderTimerCh chan<- system.NetOrder, orderTimerTimedOutCh <-chan system.NetOrder,
	elevatorConnectedCh <-chan int, elevatorDisconnectedCh <-chan int, removeOrderCh chan<- system.ButtonEvent){

	elevators := initiateElevators()
	elevatorsOnline := make(map[int]bool)

	orderToSendCh := make(chan system.NetOrder)
	go sendOrder(orderToSendCh, orderToSelfCh, orderThroughNetCh, orderTimerCh, messageTimerCh)

	for {
		select {
		case hallOrder := <-hallOrderCh: //BARE HER SOM SEND INN TIL HALL ORDER
			designatedID := getDesignatedElevatorID(hallOrder, elevators, elevatorsOnline)
			orderToSendCh <- system.NetOrder{ReceivingElevatorID: designatedID, SendingElevatorID: system.ElevatorID,
				Order: hallOrder, ReassignNum: 0}

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
				messageTimerTimedOut.ReceivingElevatorID = designatedID
				messageTimerTimedOut.ReassignNum += 1
				orderToSendCh <- messageTimerTimedOut
			}

		case elevatorConnected := <-elevatorConnectedCh:
			elevatorsOnline[elevatorConnected] = true

		case elevatorDisconnected := <-elevatorDisconnectedCh:
			elevatorsOnline[elevatorDisconnected] = false

		case orderTimerTimedOut := <-orderTimerTimedOutCh:

			elevOrds := elevators[orderTimerTimedOut.ReceivingElevatorID].Orders
			timedOutFlr := orderTimerTimedOut.Order.Floor
			timedOutBtn := int(orderTimerTimedOut.Order.Button)
			if elevOrds[timedOutFlr][timedOutBtn] != 0{
				fmt.Print("Order timer timed out. Order ")
				if orderTimerTimedOut.ReceivingElevatorID == system.ElevatorID {
					removeOrderCh <- orderTimerTimedOut.Order
					elev := elevators[system.ElevatorID]
					elev.Orders[timedOutFlr][timedOutBtn] = 0
					elevators[system.ElevatorID] = elev
					shareOwnElevatorCh <- elevators[system.ElevatorID]
					fmt.Print("is removed from self and ")

				}
				designatedID := getDesignatedElevatorID(orderTimerTimedOut.Order, elevators, elevatorsOnline)
				orderTimerTimedOut.ReceivingElevatorID = designatedID
				orderTimerTimedOut.ReassignNum = 0
				fmt.Println("given to ", designatedID)
				orderToSendCh <- orderTimerTimedOut

			}
		}
	}
}

func sendOrder(orderToSendCh <-chan system.NetOrder, orderToSelfCh chan<- system.ButtonEvent,
	orderThroughNetCh chan<- system.NetOrder, orderTimerCh chan<- system.NetOrder,
	messageTimerCh chan<- system.NetOrder){
	for{
		select{
		case orderToSend := <- orderToSendCh:
			if orderToSend.ReceivingElevatorID == system.ElevatorID {
				orderToSelfCh <- orderToSend.Order
				fmt.Println("Order",  orderToSend,"sent to self")
				if orderToSend.ReassignNum == 0 {
					orderTimerCh <- orderToSend
				}
			} else {
				orderThroughNetCh <- orderToSend
				fmt.Println("Order", orderToSend,"sent throught net")
				messageTimerCh <- orderToSend
			}
		}
	}
}