# Research: Multi-Instance Activities

## 1. Ledger subject for inner instances

**Decision**: Keep the loop host's BPMN element id as `Element.id`. Each inner iteration appends ACTIVITY (or PROCESS for MI sub-process) EVENTs with an additive `loop_instance_index` on `ActivityPayload` (and equivalent on sub-process scope entry). No new `Element.Type`.

**Rationale**: Constitution III — Job/timer/message are payloads of the activity, not peer subjects; loop iterations are lifecycle stages of the same activity.

**Alternatives considered**:

- Mint synthetic element ids per inner instance (`Task_1`, `Task_2`) — breaks deploy lookup and BPMN audit alignment.
- New `LOOP_INSTANCE` subject — second ontology; rejected.

## 2. Token model for parallel vs sequential

**Decision**:

- On loop entry, compute `numberOfInstances` from cardinality expression or collection size.
- **Parallel**: activate all inner instances immediately — one token per inner instance at the activity (each tagged with `loop_instance_index`).
- **Sequential**: activate only index `0`; on inner COMPLETED, activate index `n+1` until all started instances complete or completion condition fires.
- Maintain projection `MultiInstanceLoop` per host element: `{host_element_id, host_token_id, total, active, completed, sequential, output_collection[], cancelled}` rebuilt from EVENTs.

**Rationale**: Matches BPMN semantics and reuses existing per-token waiting/Complete model. `Complete` must target a specific inner instance (token id or loop index).

**Alternatives considered**:

- Single token with internal counter only — cannot represent concurrent waiting user tasks / jobs.
- Sub-process per inner instance for all MI — overkill for flat tasks; reserved for MI Sub-Process story.

## 3. Cardinality and collection input

**Decision**:

- Compile `loopCardinality` text and `completionCondition` text at deploy; evaluate at runtime with `processing/expr` against instance variables plus injected loop counters (`nrOfInstances`, `nrOfActiveInstances`, `nrOfCompletedInstances`, `loopCounter`).
- Collection: read `loopDataInputRef` variable as JSON array; `numberOfInstances = len(array)`; bind `inputDataItem` name to `array[i]` as a variable visible during inner instance `i`.
- Missing/non-array collection → zero instances, immediate host completion (spec edge case).
- Cardinality ≤ 0 → zero instances.

**Rationale**: Reuses existing expression and JSON variable encoding; no Camunda FEEL engine in scope.

**Alternatives considered**:

- Require only integer literal cardinality in v1 — too limiting for Story 3 collection path.
- Fail instance on bad collection — spec says treat as empty.

## 4. Completion condition and early cancel

**Decision**: After each inner COMPLETED (or inner TERMINATED when error path ends an inner scope), evaluate `completionCondition`. If true (or default all-complete), cancel remaining active inner tokens (TERMINATING → TERMINATED), assemble output collection, complete host once, take single outgoing flow.

Default when `completionCondition` absent: `nrOfCompletedInstances == nrOfInstances` (All behavior).

`behavior=One` on loop characteristics maps to a deploy-time default completion expression equivalent to `nrOfCompletedInstances >= 1` unless a custom `completionCondition` overrides.

**Rationale**: Spec FR-006–FR-008; Camunda-compatible counter names.

**Alternatives considered**:

- Poll completion only at end — cannot support early exit.

## 5. Output collection

**Decision**: When `loopDataOutputRef` and `outputDataItem` are set, on each inner COMPLETED read the named element variable from that iteration and append to an in-memory slice on `MultiInstanceLoop`; on host completion write the JSON array to `loopDataOutputRef`.

**Rationale**: Spec FR-010; standard BPMN data association pattern.

**Alternatives considered**:

- Defer output collection to a later story — spec includes it in FR-010; keep in increment but test in Story 3 fixture.

## 6. MI embedded Sub-Process

**Decision**: Loop host is the Sub-Process element. Each inner instance: host token with `ScopeHost=true` + `loop_instance_index`, child token(s) inside scope per existing SubProcess rules. Inner scope COMPLETED increments loop counters; same completion/cancel machinery as flat activities.

**Rationale**: Reuses `ScopeHost` pattern from embedded Sub-Process and Call Activity parking.

**Alternatives considered**:

- Flatten sub-process body N times at deploy — violates no parallel graph / duplicate element ids.

## 7. Recover

**Decision**: Replay ACTIVITY/PROCESS EVENTs with `loop_instance_index` and rebuild `MultiInstanceLoop` counters and per-index token presence. Waiting inner tokens restore `Complete` targets. Partial output collection rebuilt from inner COMPLETED variable deltas if persisted on payload (or re-derived from instance variables if stored per-index — prefer payload audit).

**Rationale**: Constitution II; same pattern as Call Activity parent/child rebuild.

## 8. Protocol and Complete API

**Decision**: Add optional fields:

- `ActivityPayload.loop_instance_index` (int32, 0-based; unset/0 for non-MI)
- `Token.loop_instance_index` on projection and `engine.v1.Token` for `GetInstance`
- `CompleteRequest` optional `loop_instance_index` or require `token_id` (preferred: **token_id** already unique — document that MI completes use the inner token id; loop index is denormalized for audit)

Deploy rejects `complexBehaviorDefinition` and none/one behavior event refs.

**Rationale**: Wire compatibility; clients already pass `token_id` to `Complete`.

**Alternatives considered**:

- New `CompleteLoopInstance` RPC — unnecessary surface area.

## 9. Deploy validation

**Decision**: `deploy` extracts `multiInstanceLoopCharacteristics` from `element.Activity` / SubProcess. Supported: `isSequential`, `loopCardinality`, `loopDataInputRef`, `inputDataItem`, `loopDataOutputRef`, `outputDataItem`, `completionCondition`, `behavior` One/All only. Unsupported element types or complex behavior → `UNSUPPORTED_ELEMENT` at deploy.

**Rationale**: Fail fast; matches spec FR-014.
