# Distributor module
Distributes the orders based on the calculation done by the cost function.

## Goroutines
### `OrderDistributor`
**Receiving channels**
- `hallOrderCh`: gets a ButtonEvent type from system.
- `otherElevatorCh`, `ownElevatorCh`: gets an Elevator type from system.
- `messageTimerTimedOutCh`, `orderTimerTimedOutCh`: gets a NetOrder type from system.
-`elevatorConnectedCh`, `elevatorDisconnectedCh`: gets an integer.

**Sending channels**
-`orderToSelfCh`, `removeOrderCh`: sends a ButtonEvent type from system.
- `shareOwnElevatorCh`: sends an Elevator type from system.
- `orderThroughNetCh`, `messageTimerCh`, `orderTimerCh`: sends a NetOrder type from system.

### `sendOrder`
**Receiving channels**
-`orderToSendCh`: gets a NetOrder type from system.

**Sending channels**
-`orderToSelfCh`: sends a ButtonEvent type from system.
-`orderThroughNetCh`, `orderTimerCh`, `messageTimerCh`: sends a NetOrder type form system.
