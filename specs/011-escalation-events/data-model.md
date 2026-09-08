# Data Model: Escalation Events

## Deploy-time

| Map / field | Meaning |
|-------------|---------|
| `escalationCatch[boundaryID]` | Escalation code (empty = catch-all) |
| `escalationBoundaries[activityID]` | Boundary ids |
| `escalationEnds[endID]` | Code thrown by escalation end |
| `throwEvents[id].Kind = escalation` | Intermediate escalation throw; Name = code |
| `EventSubProcess.Kind = escalation` | `EscalationCode` on spec |

## Runtime / ledger

| Field | Meaning |
|-------|---------|
| `Element.INTENT_ESCALATION_THROWN` | Throw recorded |
| `EventPayload.escalation_code` | Code on throw |

No new projection token fields required (catch is synchronous on throw).
