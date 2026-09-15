# Data model

| Concept | Representation |
|---------|----------------|
| Transaction | FlowElements + TYPE_TRANSACTION |
| Cancel End / Boundary | deploy maps; CatchKindCancel |
| PendingCompensation | + CancelBoundaryID, CancelTransactionID, CancelHostTokenID |

Proto: `Element.Type.TYPE_TRANSACTION = 27`.
