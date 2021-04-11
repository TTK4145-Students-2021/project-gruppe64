package distributor

import (
	"fmt"
	"realtimeProject/project-gruppe64/system"
	"strconv"
)



// GOROUTINE:
func OrderDistributor(hallOrder <-chan system.ButtonEvent, elevatorInfo <-chan system.ElevatorInformation, ownElevator <-chan system.Elevator, sendingOrderThroughNet chan<- system.SendingOrder, orderToSelf chan<- system.ButtonEvent, messageTimer chan<- system.SendingOrder, messageTimerTimedOut <-chan system.SendingOrder, orderTimer chan<- system.SendingOrder, orderTimerTimedOut <- chan system.SendingOrder){
	elevs := initiateElevatorsTagged()
	distributedOrders := make(map[string][]system.SendingOrder)
	for {
		select {
		case hallOrd := <-hallOrder:
			switch hallOrd.Button {
			case system.BT_HallUp:
				elevs.HallOrders[hallOrd.Floor][0] = true
			case system.BT_HallDown:
				elevs.HallOrders[hallOrd.Floor][1] = true
			default:
				break
			}
			designatedID := getDesignatedElevatorID(elevs)
			if designatedID == system.ElevatorID {
				fmt.Println("ORDER TO SELF")
				orderToSelf <- hallOrd
			} else {
				sOrd := system.SendingOrder{ReceivingElevatorID: designatedID, SendingElevatorID: system.ElevatorID, Order: hallOrd}
				sendingOrderThroughNet <- sOrd
				messageTimer <- sOrd
				orderTimer <- sOrd
				distributedOrders[strconv.Itoa(designatedID)] =  []system.SendingOrder{sOrd} //NEEDS TO BE APPEND IF SOMETHING IS THERE
			}
			switch hallOrd.Button { //sletter ordren fra her og nå
			case system.BT_HallUp:
				elevs.HallOrders[hallOrd.Floor][0] = false
			case system.BT_HallDown:
				elevs.HallOrders[hallOrd.Floor][1] = false
			default:
				break
			}

		case msgTimedOut := <- messageTimerTimedOut:
			orderToSelf <- msgTimedOut.Order

		case ordTimedOut := <- orderTimerTimedOut:
			for key, dOrds := range distributedOrders{
				if key == strconv.Itoa(ordTimedOut.ReceivingElevatorID){
					for _, dOrd := range dOrds{
						if dOrd == ordTimedOut{
							orderToSelf <- ordTimedOut.Order
						}
					}
				}
			}
		case elev := <-ownElevator:
			elevs.States[strconv.Itoa(system.ElevatorID)] = getUpdatedElevatorTagged(system.ElevatorInformation{ID: system.ElevatorID, Floor: elev.Floor, MotorDirection: elev.MotorDirection, Orders: elev.Orders, Behaviour: elev.Behaviour})

		case elevInfo := <-elevatorInfo:
			if removeExecutedOrders(elevInfo, distributedOrders[strconv.Itoa(elevInfo.ID)]) != nil{ //måtte legge til dette for å ikke få feilmelding
				distributedOrders[strconv.Itoa(elevInfo.ID)] = removeExecutedOrders(elevInfo, distributedOrders[strconv.Itoa(elevInfo.ID)])
			}
			elevs.States[strconv.Itoa(elevInfo.ID)] = getUpdatedElevatorTagged(elevInfo)
		}
	}
}

/*
func elevatorOperationCheck(){
	//Sjekk ordrer som er kommet inn sist, er de like de andre for så så lenge? I så fall opererer ikke heisen. Også en timer på når akk
	//den heisen mottok sist oppdatering. On lenge siden? tilsvarende som om de er like. Heisen opererer ikke.
	//Tenk litt som timerene for order og message
}
*/

