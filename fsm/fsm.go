package fsm

import (
	"fmt"
	"realtimeProject/project-gruppe64/io"
)


var (
	elevator Elevator
	outputDevice io.ElevInputDevice
)

func fsm_init(){

}

func setAllLights(es <-chan Elevator){
	for floor := 0; floor < io.N_FLOORS; floor++ {
		for btn := 0; btn < io.N_BUTTONS; btn++ {
			outputDevice.requestButtonLight(floor, es.requests[floor][btn])
		}
	}
}

func fsm_onInitBetweenFloors(){
	outputDevice.motorDirection(D_Down)
	elevator.dirn = D_Down
	elevator.behaviour = EB_Moving
}

func fsm_onRequestButtonPress(btn_floor <-chan int, btn_type <-chan Button){
	fmt.Printf("\n\n%s(%d, %s)\n", __FUNCTION__, btn_floor, Elevio_b)
}

func fsm_onFloorArrival(newFloor <-chan int) {}

func fsm_onDoorTimeout(){}
