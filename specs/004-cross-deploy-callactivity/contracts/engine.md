# Contracts: Cross-deployment Call Activity

No new `EngineService` RPCs. Behavior is deploy validation + runtime resolution + existing Call Activity ledger intents.

## Deploy contract

### Caller-only definitions

- Call Activity with `calledElement="PaymentFlow"` MUST deploy successfully when `PaymentFlow` is **not** in the same BPMN file.
- Compiled `CallActivity` has `external_callee=true` and empty `start_event_id`.
- Input/output `dataInputAssociation` / `dataOutputAssociation` with `sourceRef`/`targetRef` only (unchanged).

### Bundled definitions (regression)

- When `calledElement` references a process in the same file, behavior identical to pre-004: embed callee, compile start event, child uses **same** deployment id as caller.

### Deploy failures (unchanged)

- Empty `calledElement` → `UNSUPPORTED_ELEMENT`
- Embedded callee fails M1 validation → deploy error
- Multi-instance Call Activity → `UNSUPPORTED_ELEMENT`

## Runtime contract

### Call Activity activation (external)

1. Engine resolves `calledElement` to latest deployed revision for that process id.
2. On success: CALL_ACTIVITY ACTIVATING → ACTIVATED; caller token `waiting` with `called_process_instance_id`.
3. Child PROCESS COMMAND appended with callee `deployment_id` and `process_version`.
4. Child instance `GetInstance`: `deployment_id` = callee; `parent_process_instance_id` = caller.

### Call Activity activation (callee missing)

1. Resolution fails before ACTIVATED.
2. REJECTION with code `NOT_FOUND` on the entering command chain.
3. Caller projection MUST NOT show waiting on Call Activity with a child id.

### IO mapping

- Input mapping: caller variables → child start variables (same as same-deployment).
- Output mapping: child variables → caller on `Complete` of Call Activity (same as same-deployment).

### Boundaries

- Interrupting boundary on Call Activity terminates child via existing `PublicationTerminateChild`; child loaded by `inst.DeploymentID` (callee).

## event.v1 (optional additive)

### `ActivityPayload` (optional extend)

| Field | Type | When |
|-------|------|------|
| `called_process_id` | string | CALL_ACTIVITY ACTIVATED; callee process id for audit |

Skip if implementation relies on child instance `process_id` + linkage only. If added: regenerate via `protocol/proto/build.sh`; FILE-compatible additive field only.

No new `Element.Intent` or `Element.Type`.

## engine.v1

No changes required. Existing fields:

- `Instance.deployment_id` — distinguishes caller vs child
- `Instance.parent_process_instance_id` — child → caller
- `Token.called_process_instance_id` — caller → child

## processing Go API (non-wire)

```go
// deploy.CallActivity gains ExternalCallee bool

// handlers — wired from processing init:
func SetCalleeResolver(fn func(processID string) (deploymentID string, err error))

// handlers.Publication gains:
//   CalledDeploymentID string
```

`startCalledInstance` MUST use `CalledDeploymentID` for child COMMAND/projection/Enter.

## Compatibility

- Same-deployment fixtures: no caller-visible change.
- Clients ignoring optional `called_process_id` remain compatible.
- Recover: requires both caller and callee deployment artifacts in `deploy.Store` (documented in quickstart).
