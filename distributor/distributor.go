package distributor

import (
	"fmt"
	"realtimeProject/project-gruppe64/system"
	"strconv"
)

// GOROUTINE:

//ENDRE DISTRIBUTOR TIL Å HA GLOBAL VARIABEL ELEVATORS SOM ER ELEVATOR INFORMATION IKKE ELEVATOR TAGGED
func OrderDistributor(hallOrder chan system.ButtonEvent, elevatorInfo <-chan system.ElevatorInformation,
	ownElevator <-chan system.Elevator, sendingOrderThroughNet chan<- system.SendingOrder,
	orderToSelf chan<- system.ButtonEvent, messageTimer chan<- system.SendingOrder,
	messageTimerTimedOut <-chan system.SendingOrder, orderTimer chan<- system.SendingOrder,
	orderTimerTimedOut <- chan system.SendingOrder, elevatorIDConnected <-chan int, elevatorIDDisconnected <-chan int){

	elevators := initiateElevatorsTagged()
	disconnectedElevators := make(map[string]bool)
	distributedOrders := make(map[string][]system.SendingOrder)
	for {
		select {
		case hallOrd := <-hallOrder:
			switch hallOrd.Button {
			case system.BT_HallUp:
				elevators.HallOrders[hallOrd.Floor][0] = true
			case system.BT_HallDown:
				elevators.HallOrders[hallOrd.Floor][1] = true
			default:
				break
			}
			fmt.Println(elevators)
			fmt.Println(disconnectedElevators)
			designatedID := getDesignatedElevatorID(elevators, disconnectedElevators)
			fmt.Println(designatedID)
			if designatedID == system.ElevatorID {
				orderToSelf <- hallOrd
			} else {
				sOrd := system.SendingOrder{ReceivingElevatorID: designatedID, SendingElevatorID: system.ElevatorID, Order: hallOrd}
				sendingOrderThroughNet <- sOrd
				messageTimer <- sOrd
				orderTimer <- sOrd
				if distributedOrders[strconv.Itoa(designatedID)] != nil {
					distributedOrders[strconv.Itoa(designatedID)] = append(distributedOrders[strconv.Itoa(designatedID)], sOrd)
				} else {
					distributedOrders[strconv.Itoa(designatedID)] = []system.SendingOrder{sOrd}
				}
			}
			switch hallOrd.Button { //sletter ordren fra her og nå
			case system.BT_HallUp:
				elevators.HallOrders[hallOrd.Floor][0] = false
			case system.BT_HallDown:
				elevators.HallOrders[hallOrd.Floor][1] = false
			default:
				break
			}

		case msgTimedOut := <-messageTimerTimedOut:
			if distributedOrders[strconv.Itoa(msgTimedOut.ReceivingElevatorID)] != nil {
				distributedOrders[strconv.Itoa(msgTimedOut.ReceivingElevatorID)] = removeOrderFromOrders(msgTimedOut, distributedOrders[strconv.Itoa(msgTimedOut.ReceivingElevatorID)])
			}
			//Should not send to self, should put in to hall order again:
			hallOrder <- msgTimedOut.Order

		case ordTimedOut := <-orderTimerTimedOut:
			for key, dOrds := range distributedOrders {
				if key == strconv.Itoa(ordTimedOut.ReceivingElevatorID) {
					for _, dOrd := range dOrds {
						if dOrd == ordTimedOut {
							distributedOrders[key] = removeOrderFromOrders(ordTimedOut, distributedOrders[key])
							orderToSelf <- ordTimedOut.Order
						}
					}
				}
			}

		case elev := <-ownElevator:
			system.LogElevator(elev)
			elevators.States[strconv.Itoa(system.ElevatorID)] = getUpdatedElevatorTagged(system.ElevatorInformation{
				ID: system.ElevatorID, Floor: elev.Floor, MotorDirection: elev.MotorDirection,
				Orders: elev.Orders, Behaviour: elev.Behaviour})

		case elevInfo := <-elevatorInfo:
			if elevInfo.ID != system.ElevatorID {
				setHallButtonLights(elevInfo)
				if removeExecutedOrders(elevInfo, distributedOrders[strconv.Itoa(elevInfo.ID)]) != nil { //måtte legge til dette for å ikke få feilmelding
					distributedOrders[strconv.Itoa(elevInfo.ID)] = removeExecutedOrders(elevInfo, distributedOrders[strconv.Itoa(elevInfo.ID)])
				}
				elevators.States[strconv.Itoa(elevInfo.ID)] = getUpdatedElevatorTagged(elevInfo)
			}
		case elevIDCnct := <-elevatorIDConnected:
			disconnectedElevators[strconv.Itoa(elevIDCnct)] = false

		case elevIDDiscnct := <-elevatorIDDisconnected:
			disconnectedElevators[strconv.Itoa(elevIDDiscnct)] = true

		}
	}
}
