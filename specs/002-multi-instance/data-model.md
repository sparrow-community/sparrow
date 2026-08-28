# Data Model: Multi-Instance Activities

## Entities

### Multi-instance spec (compile-time, in `deploy`)

| Field | Notes |
|-------|--------|
| `host_element_id` | User Task, Service Task, or Sub-Process id |
| `sequential` | from `isSequential` |
| `cardinality_expr` | `loopCardinality` text; empty when collection-driven |
| `input_collection_var` | `loopDataInputRef` |
| `input_element_var` | `inputDataItem` name |
| `output_collection_var` | `loopDataOutputRef` |
| `output_element_var` | `outputDataItem` name |
| `completion_expr` | `completionCondition` text; empty → all-complete default |
| `behavior` | `All` or `One` (maps to default completion when no custom expr) |

Validation: reject MI on unsupported types; reject `complexBehaviorDefinition` and behavior event refs.

### Multi-instance loop (runtime projection)

| Field | Notes |
|-------|--------|
| `host_element_id` | loop host |
| `host_token_id` | token that arrived from incoming flow (parked while loop runs) |
| `total_instances` | `numberOfInstances` |
| `active_instances` | waiting or running inner work |
| `completed_instances` | finished successfully |
| `sequential` | mode flag |
| `output_items` | slice built per inner complete |
| `cancelled` | true after early completion triggered cancel of stragglers |

Keyed in `projection.Instance` (e.g. `MultiInstanceLoops map[string]*MultiInstanceLoop` by host element id).

### Token (existing, extended)

| Field | Notes |
|-------|--------|
| `LoopInstanceIndex` | 0..N-1 for inner instance tokens; -1 or omitted for non-MI |
| `ScopeHost` | true for MI Sub-Process host tokens per inner instance |
| existing fields | JobType, waiting boundaries, etc. per inner instance |

Inner User Task / Service Task tokens share the host `ElementID` but distinct `Token.ID` and `LoopInstanceIndex`.

### ActivityPayload (ledger, extended)

| Field | Notes |
|-------|--------|
| `loop_instance_index` | int32; set on inner ACTIVATED/COMPLETED/TERMINATED |
| `loop_counters` (optional) | snapshot on host COMPLETED for audit: total/active/completed |

### Loop counter variables (runtime, for expressions)

Injected into `expr.Eval` context during completion checks (names align with BPMN/Camunda conventions):

| Name | Meaning |
|------|---------|
| `nrOfInstances` | total started |
| `nrOfActiveInstances` | still active |
| `nrOfCompletedInstances` | completed successfully |
| `loopCounter` | index of the inner instance that just completed (when evaluating per-complete) |

Per-inner-instance element variable (collection input) merged into variable map for that iteration only.

## State transitions

### Parallel MI User Task (cardinality N)

```text
flow arrives → host ACTIVATING
           → N × USER_TASK ACTIVATED (loop_instance_index 0..N-1, waiting)
each Complete(inner) → USER_TASK COMPLETED
when completion satisfied → cancel remaining inner tokens (if any)
                         → USER_TASK COMPLETED (host, single outgoing)
```

### Sequential MI

```text
flow arrives → USER_TASK ACTIVATED (index 0)
complete 0   → USER_TASK COMPLETED (index 0) → ACTIVATED (index 1)
...
complete N-1 → host join → single outgoing
```

### MI Sub-Process

```text
per index i: SUB_PROCESS ACTIVATED (host ScopeHost, index i)
             → inner tokens in scope
             → scope COMPLETED (index i)
loop join → SUB_PROCESS COMPLETED (host once)
```

### Interrupting boundary on MI host

```text
boundary fires → TERMINATE all inner tokens + loop host
              → boundary outgoing (once)
```

## Validation rules

- Deploy: unsupported MI shape → `UNSUPPORTED_ELEMENT`.
- Runtime: `Complete` without resolvable token → `NOT_FOUND`.
- Collection missing or not array → zero instances, immediate host completion.
- Cardinality evaluates to non-integer or negative → treat as zero instances.
- Recover: loop state must match pre-crash counts and waiting tokens.
