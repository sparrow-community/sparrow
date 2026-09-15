# Data Model

| Concept | Representation |
|---------|----------------|
| CatchKindCompensate | `deploy.CatchKind = "compensate"` on `EventSubProcess` |
| Compensation link | `Compensation{BoundaryID: startEventID, ActivityID: enclosingSubProcessID, HandlerID: eventSubProcessID}` |
| Subscription | Existing `CompensationSubs[startEventID]` after SubProcess COMPLETED |
| PendingCompensation.Queue | May contain event sub-process element id |
| Nested throw scope | If `ScopeOf(throw)` is compensation event sub-process → use its ParentScopeID |

No proto changes.
