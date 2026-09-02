# Contracts: Call Activity boundary and compensation parity

No new `EngineService` or `JobService` RPCs. Behavior is deploy validation + existing boundary/compensation COMMAND/EVENT semantics on `CALL_ACTIVITY` and `BOUNDARY_EVENT`.

## Deploy contract

### Boundary attachment

- `boundaryEvent.attachedToRef` MAY reference a `callActivity` id in the same process (or caller process for cross-deploy call).
- Supported boundary kinds: timer, message, signal, error (`errorEventDefinition`), compensation (`compensateEventDefinition` + association).
- Same per-kind uniqueness rule as User Task / SubProcess: at most one timer, one message, one signal, one error, one compensation boundary per Call Activity.

### Deploy failures (unchanged)

- Boundary on unknown element id → `UNSUPPORTED_ELEMENT`
- Duplicate kind on same host → `UNSUPPORTED_ELEMENT`
- Invalid compensation handler → `UNSUPPORTED_ELEMENT`

## Runtime contract

### Call Activity enter (waiting)

1. CALL_ACTIVITY ACTIVATING → ACTIVATED with `called_process_instance_id` and boundary payload fields when boundaries attached.
2. Child instance started via existing `PublicationStartChild` (unchanged from 004).

### Interrupting boundary (timer / message / signal / error)

1. Boundary completes on host token while Call Activity waiting.
2. Child instance terminated if `called_process_instance_id` set.
3. Host CALL_ACTIVITY TERMINATED (or error-boundary equivalent sequence).
4. Caller continues on boundary outgoing flow.

### Non-interrupting boundary

1. Boundary completes without terminating host or child.
2. Spawned token on boundary outgoing; host remains waiting at Call Activity.

### Compensation

1. On CALL_ACTIVITY COMPLETED (normal child completion path): compensation boundary ACTIVATED → subscription recorded.
2. Compensate throw in scope runs handler per existing ordering (reverse completion order).
3. No subscription if Call Activity never reached COMPLETED (interrupted/terminated).

### Cross-deployment

- Boundary definitions and handlers live on **caller** deployment.
- Child terminate uses child instance `deployment_id` (callee) — unchanged from 004.

## event.v1

No schema changes. Relevant existing fields:

- `ActivityPayload`: `boundary_id`, `message_boundary_id`, `signal_boundary_id`, `due_unix_ms`, `duration`, `message_name`, `signal_name`, `called_process_instance_id`
- `EventPayload.compensation_handler_id` on compensation boundary ACTIVATED

## engine.v1

No changes. `GetInstance` token fields expose boundary and child linkage.

## processing Go API (non-wire)

```go
// deploy.validateBoundaryHost — add callActivity to allowed hosts

// CallActivityHandler.OnEnter — attachBoundary on ACTIVATED payload
// CallActivityHandler.OnComplete — cancelAttachedBoundary + subscribeCompensation
```

## Compatibility

- Existing Call Activity tests without boundaries: unchanged.
- Clients relying on deploy rejection of boundaries on Call Activity will need updated BPMN (intended behavior change).
