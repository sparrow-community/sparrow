# Data model

| Concept | Representation |
|---------|----------------|
| Intermediate escalation catch | `escalationCatch[catchID] = code` (empty = catch-all) |
| Wait payload | `EventPayload.EscalationCode` |
| Delivery | `Executor.Complete` on matching waiting tokens |

No proto changes.
