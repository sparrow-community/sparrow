# Data Model: Cross-deployment Call Activity

## Entities

### CallActivity (compile-time, in `deploy`)

| Field | Embedded callee | External callee |
|-------|-----------------|-----------------|
| `id` | callActivity element id | same |
| `called_process_id` | from `calledElement` | from `calledElement` |
| `start_event_id` | compiled from embedded process | **empty** (resolved at runtime from callee deployment) |
| `inputs` / `outputs` | compiled mappings | same |
| `external_callee` | `false` | `true` |

Validation:

- `calledElement` non-empty always required.
- Embedded: callee in catalog; callee scope validated as today; cannot call root process.
- External: IO mapping compile rules unchanged; no catalog presence required.
- Multi-instance Call Activity: still unsupported (unchanged).

### Deployment indexes (unchanged semantics)

| Index | Embedded | External |
|-------|----------|----------|
| `callActivities[id]` | full spec | full spec with `external_callee=true` |
| `calledProcesses[processID]` | embedded process body | **absent** |
| `calledProcessOwner` | set | **absent** for that process |

### Publication (runtime, `handlers`)

| Field | Purpose |
|-------|---------|
| `DeploymentID` | Caller deployment (parent context) |
| `CalledDeploymentID` | **New.** Callee deployment for child start; set at OnEnter |
| `CalledProcessID` | Callee process id |
| `ChildInstanceID` | Preallocated UUIDv7 |
| `ParentInstanceID`, `CallActivityID`, `HostTokenID` | Unchanged |

### Instance (runtime projection)

| Field | Parent (caller) | Child (callee) |
|-------|-----------------|----------------|
| `DeploymentID` | caller deployment id | **callee deployment id** |
| `ProcessID` | caller process id | callee process id |
| `ParentProcessInstanceID` | empty | caller instance id |
| `ParentElementID` / `ParentTokenID` | empty | Call Activity id / host token |

Token on caller while waiting:

| Field | Value |
|-------|--------|
| `Status` | `waiting` |
| `ElementID` | Call Activity id |
| `CalledProcessInstanceID` | preallocated child id |

Invariant: child `DeploymentID` ≠ parent `DeploymentID` when `external_callee=true` and deployments differ.

### Ledger (EVENT/COMMAND)

Child PROCESS activation COMMAND/EVENT:

- `deployment_id` = callee deployment id
- `process_version` = callee version
- `ProcessPayload.parent_*` = caller linkage (unchanged)

Parent CALL_ACTIVITY ACTIVATED:

- `ActivityPayload.called_process_instance_id` = child id (unchanged)
- Optional: `called_process_id` for audit (see contracts)

## State transitions

### Successful cross-deploy call

```text
Caller token arrives at Call Activity
  → OnEnter: resolve callee dep (latest revision)
  → ACTIVATING / ACTIVATED (waiting, child id set)
  → Publication StartChild (caller dep, callee dep, child id)
  → startCalledInstance: child CREATE_INSTANCE COMMAND (callee dep)
  → child PROCESS ACTIVATED + Enter(startEvent on callee dep)
  → child runs ...
  → child PROCESS COMPLETED → ResumeParent → caller Complete Call Activity
```

### Callee not found (FR-006)

```text
Caller token arrives at Call Activity (external)
  → OnEnter: resolve fails NOT_FOUND
  → no ACTIVATED; command REJECTION
  → caller token unchanged (not waiting on call)
```

### Recover mid-call

```text
Log contains: parent CALL_ACTIVITY ACTIVATED + child CREATE_INSTANCE (callee dep)
  → Replay restores parent waiting + child active with correct deployment ids
  → Redrive unfinished StartChild uses same CalledDeploymentID + ChildInstanceID
```

## Rejection codes

| Code | When |
|------|------|
| `NOT_FOUND` | `calledElement` not deployed at external OnEnter resolve |
| `UNSUPPORTED_ELEMENT` | Deploy: empty `calledElement`; invalid IO mapping (unchanged) |

No new stable codes required beyond existing `NOT_FOUND` on enter failure.
