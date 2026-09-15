# Data model

| Concept | Representation |
|---------|----------------|
| ComplexGateway | FlowElements + TYPE_COMPLEX_GATEWAY |
| activationCondition | deploy map id → expr text |
| Join fire | TerminateJoinPeers + TakeOutgoing / Fork |

No proto enum change (type 22 reserved).
