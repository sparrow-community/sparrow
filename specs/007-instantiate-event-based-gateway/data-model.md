# Data Model: Instantiate Event-Based Gateway

## Entities

### Instantiate event-based gateway (compile-time)

Process-level exclusive `eventBasedGateway` with `instantiate=true`.

| Field / rule | Notes |
|--------------|--------|
| id | Gateway element id; becomes process entry when no startEvent |
| instantiate | Must be true for this role |
| eventGatewayType | Exclusive (default); Parallel instantiate rejected |
| incoming | Must be empty (no sequence flow targets this gateway) |
| outgoing | ≥2 flows to intermediateCatchEvent |

### Process entry (compile-time resolve)

| Mode | Entry element |
|------|----------------|
| Classic | First process `startEvent` with outgoing (unchanged) |
| Instantiate | Sole valid process-level instantiate exclusive EBG when `StartEvents` empty |

Mutually exclusive: startEvent(s) and instantiate EBG must not coexist on the same process.

### Tokens / events (runtime)

No new types. CreateInstance:

1. PROCESS ACTIVATING/ACTIVATED (existing)
2. Enter instantiate EBG → EVENT_BASED_GATEWAY instant complete + fork tokens to catches
3. Intermediate catch ACTIVATED (timer/message/signal waiting) as today
4. First catch complete → sibling TERMINATED → continue on winning path

### Deployment indexes

Reuse existing catch indexes (`timerCatch`, `messageCatch`, `signalCatch`) and `EventBasedSiblings`. Optional: cache instantiate entry id during compile; not required if `StartEventID` scans.

## State transitions

```text
CreateInstance (instantiate-only process)
  → PROCESS started
  → Enter EBG_instantiate
  → fork → Catch_A waiting, Catch_B waiting
  → (event A) Catch_A COMPLETED, Catch_B TERMINATED → path A → End
```

Recover while both waiting restores the same waiting catches; winning event yields the same path.
