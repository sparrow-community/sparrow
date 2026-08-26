# Data Model: Engine Completeness (next increment)

## Entities

### Deployment (existing, extended)

| Field | Notes |
|-------|--------|
| `deployment_id` | UUIDv7; immutable snapshot key |
| BPMN XML | stored in `deploy.Store` |
| compiled `element.Process` set | including sibling processes for Call Activity |
| `process_id` → `process_version` | assigned at compile; first deploy of a process id is 1 |

Validation: Call Activity `calledElement` must resolve in the same definitions; called process must have a none start; mappings must name existing source/target items when present. Drop v1 uniqueness "one Call Activity per called process" and "no Call Activity inside called process".

### Process instance (existing, extended)

| Field | Notes |
|-------|--------|
| `id` | `process_instance_id` / partition key |
| `deployment_id` / `process_version` | bound at start; never rewritten |
| `status` | active / completed / terminated |
| `variables` | JSON text map |
| `parent_process_instance_id` | set on called instances |
| `parent_element_id` | Call Activity id on the parent |
| `parent_token_id` | parked host token on the parent |

### Token (existing, extended)

| Field | Notes |
|-------|--------|
| host Call Activity token | `ScopeHost=true`; `called_process_instance_id` while the child lives |
| child tokens | only on the child instance |

### Call relationship (derived from ledger)

Not a separate store. Written on:

- parent `CALL_ACTIVITY` ACTIVATED → `ActivityPayload.called_process_instance_id`
- child `PROCESS` ACTIVATED → `ProcessPayload` parent ids

Recover rebuilds parent/child fields by replaying those EVENTs.

### Variable mapping (compile-time)

| Field | Notes |
|-------|--------|
| `call_activity_id` | |
| `inputs` | list of `{source, target}` names |
| `outputs` | list of `{source, target}` names |

Applied only at child start (inputs) and parent Call Activity complete (outputs).

### Event Sub-Process arm (existing, extended)

| Field | Notes |
|-------|--------|
| `start_event_id` / scope id | as today |
| `error_code` | empty = catch-all error start |
| message / timer / signal | unchanged |

### Process revision index

In-memory (rebuild from Store on Recover): `process_id` → sorted `(version, deployment_id)`. File store needs no new file format if version is compiled from XML + existing deployments.

## State transitions

### Call Activity

```text
parent: CALL_ACTIVITY ACTIVATING → ACTIVATED (wait, ScopeHost)
child:  PROCESS ACTIVATING → ACTIVATED → … → PROCESS COMPLETED | TERMINATED
parent: CALL_ACTIVITY COMPLETING → COMPLETED (on child COMPLETED)
     or CALL_ACTIVITY TERMINATING → TERMINATED (on interrupting boundary / child TERMINATED without a catch)
```

### Error Event Sub-Process

```text
scope open → START_EVENT ACTIVATED (arm, error_code)
ERROR_THROWN matches → interrupting: terminate remaining scope work, enter Event Sub-Process
                    → non-interrupting: mint Event Sub-Process token, keep remaining work
```

### Revision

```text
Deploy(process P) → version n+1 snapshot
CreateInstance(deployment_id) → bind that snapshot
CreateInstance(process_id) → bind latest snapshot of P
CreateInstance(process_id, version) → bind that snapshot or reject NOT_FOUND
```

## Validation rules

- Missing mapping source: skip (no copy), do not reject at runtime.
- CreateInstance against unknown process id or version: `NOT_FOUND`.
- Call Activity with unresolved `calledElement`: `UNSUPPORTED_ELEMENT` at deploy.
- Error Event Sub-Process start with a definition other than error/message/timer/signal: still unsupported.
