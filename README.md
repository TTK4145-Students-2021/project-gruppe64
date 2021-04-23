# TTK4145 Sanntidsprogrammering

## Elevator Project
The goal of the Elevator project is to create software for controlling `n` elevators working in parallel across `m` floors. The main requirements of the elevators can be summarized as

- **No orders are lost**
- **Multiple elevators should be more efficient than one**
- **An individual elevator should behave sensibly and efficiently**
- **The lights and buttons should function as expected**

We have chosen to implement this software using the `Go` programming language.

### Code utilized
In the project we have used the following pre-made code:
- [`bcast.go`](https://github.com/TTK4145/Network-go/blob/master/network/bcast/bcast.go), [`conn.go`](https://github.com/TTK4145/Network-go/tree/master/network/conn) and [`peers.go`](https://github.com/TTK4145/Network-go/blob/master/network/peers/peers.go) from the [Network-go module](https://github.com/TTK4145/Network-go/tree/master/network)
- [`hall_request_assigner`](https://github.com/TTK4145/Project-resources/tree/master/cost_fns) from [Project-resources](https://github.com/TTK4145/Project-resources)

Additionally, we used slightly modified versions of the pre-made code [`main.go`](https://github.com/TTK4145/driver-go/blob/master/main.go) and [`elevator_io.go`](https://github.com/TTK4145/driver-go/blob/master/elevio/elevator_io.go) from the [driver-go module](https://github.com/TTK4145/driver-go), and we based some of our code on the `C`-code given in [`elev_algo`](https://github.com/TTK4145/Project-resources/tree/master/elev_algo) from [Project-resources](https://github.com/TTK4145/Project-resources)

#### Standard libraries used
The following standard go-libraries have been used:

- `runtime`
- `time`
- `sync`
- `net`
- `os`, `os/exec`
- `encoding/json`  
- `log`  
- `io/ioutil`  
- `fmt`
- `strconv`


### Our code
Our code is divided into `main.go`, a `system`-package and the five modules `hardwareIO`, `fsm`, `distributor`, `network` and `timer`.
The code is based on peer-to-peer topology, meaning the same code can be utilized on all elevators on the network.

#### Main script
Takes care of the system's process-pair functionality, creates channels for use between modules, and initiates goroutines.

#### System-package
In the system package one will find `sys_types.go` which defines the types used in our code, `sys_funcs.go` which defines general system functionality (process-pairs, backup log), and `sys_config.go` where one can configure the system according to hardware (# elevators, # floors, elevator ID, etc.).
When running our software two files will be created in this package; a system log .json-file and a primary documentation .txt-file.

#### HardwareIO module
Consists of `hardwareIO.go` which defines the hardwareIO goroutines, and `hardwareIO_funcs.go` which defines functions used by the hardwareIO goroutines.
The hardwareIO module handles hardware input and output, as well as monitoring the motor to detect motor stop. 

#### FSM module
Consists of `fsm.go` which defines the FSM goroutines, and `fsm_funcs.go` which defines functions used by the FSM goroutines.
The FSM module controls the state of the elevator based on input from hardwareIO. It also shares the elevator state with the distributor module.

#### Distributor module
Consists of `distributor.go` which defines the distributor goroutines, and `distributor_funcs.go` which defines functions used by the distributor goroutines. 
The distributor module handles all hall orders from the elevator's hall panel, keeps track of state-information shared on the network, and times messages for order distribution as well as the order execution itself.

#### Network module
Consists of the pre-made packages `bcast.go`, `conn.go` and `peers.go`. 
It also consists of the package `sendandreceive` with `sendandreceive.go`, which defines our goroutines that handles the system's networking (sharing elevator states, sending orders, detecting peer connect/disconnect).

#### Timer module
Consists of `timer.go`, which defines three goroutine timers. They provide; timing the opening of the elevator doors, timing communication of orders over the network, and timing order execution.
