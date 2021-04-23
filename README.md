# TTK4145 Sanntidsprogrammering

## Elevator Project
The goal of the Elevator project is to create software for controlling `n` elevators working in parallel across `m` floors. The main requirements of the elevators can be summarized as

- **No orders are lost**
- **Multiple elevators should be more efficient than one**
- **An individual elevator should behave sensibly and efficiently**
- **The lights and buttons should function sd expected**

We have chosen to implement this software using The `Go` programming language.

### Utilized code
In the project we have used the following pre-made code:
- [`bcast.go`](https://github.com/TTK4145/Network-go/blob/master/network/bcast/bcast.go), [`conn.go`](https://github.com/TTK4145/Network-go/tree/master/network/conn), [`localip.go`](https://github.com/TTK4145/Network-go/blob/master/network/localip/localip.go) and [`peers.go`](https://github.com/TTK4145/Network-go/blob/master/network/peers/peers.go) from the [Network-go module](https://github.com/TTK4145/Network-go/tree/master/network)
- [`hall_request_assigner`](https://github.com/TTK4145/Project-resources/tree/master/cost_fns) from [Project-resources](https://github.com/TTK4145/Project-resources)

Additionally, we used slightly modified versions of the pre-made code [`main.go`](https://github.com/TTK4145/driver-go/blob/master/main.go) and [`elevator_io.go`](https://github.com/TTK4145/driver-go/blob/master/elevio/elevator_io.go) from the [driver-go module](https://github.com/TTK4145/driver-go), and we based some of our code on the `C`-code given in [`elev_algo`](https://github.com/TTK4145/Project-resources/tree/master/elev_algo) from [Project-resources](https://github.com/TTK4145/Project-resources)

### Our code
Our code is divided into a [`system`]()-folder, [`main.go`]() and the five modules [`distributor`](), [`fsm`](), [`hardwareIO`](), [`network`]() and [`timer`](), and [`main.go`]().

------ LEGG TIL LINKER OVER NÅR ALT ER PUSHET TIL MASTER --------------







** Bold **
`funksjon, variabel, osv`
[]() hyperlink

- bullet point