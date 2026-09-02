# Data Model: Call Activity boundary and compensation parity

## Entities

### Boundary host (deploy-time)

Call Activity element id is a valid `attachedToRef` for:

| Boundary kind | Indexed in deploy | Notes |
|---------------|-------------------|--------|
| Timer | `timerCatch` / `TimerBoundary` | interrupting or non-interrupting |
| Message | `messageCatch` / `MessageBoundary` | same |
| Signal | `signalCatch` / `SignalBoundary` | same |
| Error | `errorBoundaries` / `MatchErrorBoundary` | `cancelActivity=true` (interrupting) |
| Compensation | `compensations` | association to handler activity |

Validation: extend `validateBoundaryHost` to accept `callActivity` ids; one boundary per **kind** per host (unchanged engine rule).

### Call Activity host token (runtime projection)

On ACTIVATED (waiting), token carries (from `ActivityPayload`):

| Field | Purpose |
|-------|---------|
| `CalledProcessInstanceID` | Child instance (existing) |
| `BoundaryId` | Primary armed boundary (often timer) |
| `MessageBoundaryId` | Message boundary if distinct |
| `SignalBoundaryId` | Signal boundary if distinct |
| `DueUnixMs` / `Duration` | Timer due |
| `MessageName` / `SignalName` | Catch correlation |

Status: `waiting` until child completes, interrupting boundary fires, or terminate.

### Child instance (runtime)

Separate instance record; `ParentProcessInstanceID` / `ParentElementID` link to caller. Terminated when interrupting boundary or error boundary fires on host (`PublicationTerminateChild`).

### Compensation subscription (runtime)

After Call Activity COMPLETED:

| Field | Source |
|-------|--------|
| `CompensationSubs[boundaryID]` | BOUNDARY_EVENT ACTIVATED with `compensation_handler_id` |
| `HandlerID` | From compensation association compile |
| `HostTokenID` | Call Activity host token |
| `Seq` | Order in scope for compensate throw |

Not created if Call Activity ended via interrupting boundary (host TERMINATED, not COMPLETED).

## State transitions

### Waiting call with armed timer

```text
CALL_ACTIVITY ACTIVATED (waiting, boundary armed, child id set)
  → FireDue → boundary Complete (interrupting)
  → host TERMINATED + TerminateChild + boundary path
```

### Non-interrupting message while child active

```text
CALL_ACTIVITY ACTIVATED (waiting, message boundary armed)
  → PublishMessage → boundary SpawnOutgoing (host still waiting, child active)
  → child completes → CALL_ACTIVITY COMPLETED → cancel boundaries + subscribe compensation
```

### Compensation after successful call

```text
CALL_ACTIVITY COMPLETED
  → compensation boundary ACTIVATED (subscription in CompensationSubs)
  → later compensate throw → handler activity Enter
```

### Error at Call Activity host

```text
ThrowError on CALL_ACTIVITY token
  → fireActivityErrorBoundary → TerminateChild + host TERMINATED + boundary COMPLETED
```

## Deploy rejection (unchanged codes)

| Code | When |
|------|------|
| `UNSUPPORTED_ELEMENT` | Second boundary of same kind on one Call Activity |
| `UNSUPPORTED_ELEMENT` | Invalid boundary attachment (existing rules) |

No new rejection codes required.

## Recover invariants

- Replay CALL_ACTIVITY ACTIVATED restores boundary fields and `called_process_instance_id`.
- Replay compensation BOUNDARY_EVENT ACTIVATED restores `CompensationSubs`.
- Redrive does not duplicate child terminate or boundary COMPLETED (existing `alreadySeen`).
