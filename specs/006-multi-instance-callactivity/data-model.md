# Data Model: Multi-Instance Call Activity

## Entities

### Multi-instance spec (compile-time)

Same as `002` `MultiInstanceSpec`, now also keyed by **Call Activity** element id after deploy indexes Call Activity loop characteristics.

| Field | Notes |
|-------|--------|
| host element id | Call Activity id |
| sequential / cardinality / collection / completion / behavior | unchanged from `002` |

### Call Activity spec (compile-time)

Existing `deploy.CallActivity` (called process, IO mappings, `ExternalCallee`). MI no longer forbidden at validate time.

### Multi-instance loop (runtime projection)

Existing `projection.MultiInstanceLoop` on the **parent** instance, keyed by Call Activity id:

| Field | Notes |
|-------|--------|
| host_token_id | parked incoming token (`MultiInstanceHost`) |
| total / active / completed | join counters |
| sequential / output_items / cancelled | unchanged |

### Tokens (parent instance)

| Role | Fields |
|------|--------|
| Loop host | `MultiInstanceHost=true`, `LoopInstanceIndex=-1`, boundaries armed, waiting while loop runs |
| Inner call | `LoopInstanceIndex=0..N-1`, `CalledProcessInstanceID` set after ACTIVATED, waiting until child completes or cancelled |

### Called process instance (child)

| Field | Notes |
|-------|--------|
| own process_instance_id | one per inner iteration |
| parent_process_instance_id | parent |
| parent_element_id | Call Activity id |
| parent_token_id | **inner** host token id (not the loop host token) |
| deployment_id | callee deployment (may differ when cross-deploy) |

### Publication (engine-internal, extended)

`PublicationStartChild` gains optional snapped input variables for the child (JSON name→value), computed at inner enter from parent vars + Call Activity input mappings + MI collection element for that index.

## State transitions

### Parallel MI Call Activity (cardinality N)

```text
parent arrives → CALL_ACTIVITY host ACTIVATING/ACTIVATED (loop host, boundaries)
              → MultiInstanceStart
              → N × inner CALL_ACTIVITY ACTIVATED (each with child id) + StartChild pubs
              → N child PROCESS starts
each child completes → ResumeParent(inner token) → inner CALL_ACTIVITY COMPLETED
when join satisfied → cancel remaining inners + TerminateChild for stragglers
                   → host CALL_ACTIVITY COMPLETED once → single outgoing
                   → cancel host boundaries + subscribeCompensation if any
```

### Sequential

```text
host ACTIVATED → inner[0] + one child
child done → inner[0] COMPLETED → inner[1] + next child …
last done → host COMPLETED once
```

### Early completion / interrupting boundary

```text
completion condition true OR interrupting boundary on host
  → cancelMultiInstance: TERMINATE all active inner CALL_ACTIVITY tokens
  → TerminateChild for each CalledProcessInstanceID
  → host completes (completion path) OR host TERMINATED + boundary path
```

### Empty / zero cardinality

```text
host enter with total≤0 → InstantLifecycle + TakeOutgoing (no children)
```

## Validation rules

- Deploy accepts MI on Call Activity when loop characteristics compile (same rules as other MI hosts).
- Reject complex behavior / none-one behavior event refs (unchanged `002`).
- Reject Call Activity without `calledElement` / invalid IO (unchanged).
- At most one boundary per kind on the Call Activity host (unchanged `005`).

## Recover

Rebuild parent `MultiInstanceLoops` and tokens (including `CalledProcessInstanceID` per inner) from EVENTs; rebuild child instances from their logs. Unfinished StartChild / ResumeParent COMMANDs redrive without re-resolving callee to a newer revision when `CalledDeploymentID` / child `deployment_id` is present.
