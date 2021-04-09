package distributor

import (
	"fmt"
	"realtimeProject/Network-go/network/sendandreceive"
	"realtimeProject/fsm"
	"realtimeProject/hardwareIO"
	"strconv"
)





// GOROUTINE:
func OrderDistributor(hallOrder <-chan hardwareIO.ButtonEvent, elevatorInfo <-chan sendandreceive.ElevatorInformation, ownElevator <-chan fsm.Elevator, sendingOrderThroughNet chan<- sendandreceive.SendingOrder, orderToSelf chan<- hardwareIO.ButtonEvent, messageTimer chan<- sendandreceive.SendingOrder, messageTimerTimedOut <-chan sendandreceive.SendingOrder, orderTimer chan<- sendandreceive.SendingOrder, orderTimerTimedOut <- chan sendandreceive.SendingOrder){
	elevs := initiateElevators()
	var distributedOrders map[string][]sendandreceive.SendingOrder
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
			if designatedID == fsm.ElevatorID {
				fmt.Println("ORDER TO SELF")
				orderToSelf <- hallOrd
			} else {
				sOrd := sendandreceive.SendingOrder{ReceivingElevatorID: designatedID, SendingElevatorID: fsm.ElevatorID, Order: hallOrd}
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
			elevs.States[strconv.Itoa(fsm.ElevatorID)] = getUpdatedElevatorTagged(sendandreceive.ElevatorInformation{ID: fsm.ElevatorID, Floor: elev.Floor, MotorDirection: elev.MotorDirection, Orders: elev.Orders, Behaviour: elev.Behaviour})

		case elevInfo := <-elevatorInfo:
			if removeExecutedOrders(elevInfo, distributedOrders[strconv.Itoa(elevInfo.ID)]) != nil{ //måtte legge til dette for å ikke få feilmelding
				distributedOrders[strconv.Itoa(elevInfo.ID)] = removeExecutedOrders(elevInfo, distributedOrders[strconv.Itoa(elevInfo.ID)])
			}
			elevs.States[strconv.Itoa(elevInfo.ID)] = getUpdatedElevatorTagged(elevInfo)
		}
	}
}