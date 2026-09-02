# Research: Call Activity boundary and compensation parity

## 1. Current gap vs SubProcess

**Decision**: Parity is achieved by (a) allowing Call Activity as a boundary host at deploy time, and (b) wiring `CallActivityHandler` to the same attach/disarm/subscribe helpers SubProcess and User Task use — not by a new boundary subsystem.

**Rationale**: Code review shows:
- `validateBoundaryHost` (`deploy.go:780`) lists only `userTask`, `serviceTask`, `subProcess` — **blocks all Call Activity boundaries at deploy**.
- `CallActivityHandler.OnEnter` emits ACTIVATED with only `called_process_instance_id` — **no `attachBoundary`**.
- `CallActivityHandler.OnComplete` emits COMPLETING/COMPLETED only — **no `cancelAttachedBoundary` or `subscribeCompensation`**.
- `BoundaryEventHandler.OnComplete` interrupting path already publishes `PublicationTerminateChild` when `typ == CALL_ACTIVITY` (`boundary_event.go:213`).
- `fireActivityErrorBoundary` already calls `terminateChildPub` before terminating host (`errors.go:250`).

**Alternatives considered**:

- New CallActivity-specific boundary package — rejected; duplicates `attachBoundary` / `subscribeCompensation`.
- Scope-style error boundary (`fireScopeErrorBoundary`) for Call Activity — rejected; Call Activity is an activity host, not a scope container; use `fireActivityErrorBoundary` + `MatchErrorBoundary`.

## 2. Boundary arming on ACTIVATED

**Decision**: On Call Activity ACTIVATED (waiting), call `attachBoundary(dep, callActivityID, now, payload)` and store timer/message/signal fields on host token via existing `ActivityPayload` projection — identical to User Task / Service Task.

**Rationale**: FR-002; SubProcess uses `attachScopeBoundary` which is equivalent to `attachBoundary` for timer/message/signal. Call Activity waits on host token (not scope child token), so use `attachBoundary` directly like tasks.

**Alternatives considered**:

- Arm boundaries on child instance — wrong BPMN semantics; boundary attaches to Call Activity element on caller.

## 3. Interrupting boundary behavior

**Decision**: Reuse existing interrupting `BoundaryEventHandler` flow: TERMINATING/TERMINATED host → boundary COMPLETED → `PublicationTerminateChild` → `TakeOutgoing`. No change unless tests reveal ordering bug (child terminate must flush after host lock release as today).

**Rationale**: FR-003; code path exists; 004 cross-deploy child termination uses `inst.DeploymentID` on child record.

**Alternatives considered**:

- Terminate child before host TERMINATED events — current order matches SubProcess interrupt pattern; keep unless test fails.

## 4. Non-interrupting boundaries

**Decision**: Reuse non-interrupting branch in `BoundaryEventHandler.OnComplete`: `disarmAttachedBoundary` + `SpawnOutgoing` without terminating host; child remains active.

**Rationale**: FR-004; same handler branch as tasks/SubProcess; Call Activity host stays `waiting` with `called_process_instance_id`.

**Alternatives considered**:

- Special-case Call Activity to always interrupt — contradicts BPMN non-interrupting semantics.

## 5. Compensation subscription

**Decision**: On normal Call Activity COMPLETED (child finished, `resumeParentCall` → `Complete`), append `subscribeCompensation(dep, callActivityID, tokenID)` in `OnComplete` — same as SubProcess.

**Rationale**: FR-005/FR-006; `CompensationOf(activityID)` keyed by host id; handler runs in caller instance scope.

**Alternatives considered**:

- Subscribe on child COMPLETED in callee instance — wrong scope; compensation boundary is on caller definition.

## 6. Error boundary on Call Activity host

**Decision**: Deploy error boundary attached to Call Activity; runtime uses existing `MatchErrorBoundary` + `fireActivityErrorBoundary` when error thrown **at** the Call Activity token (not errors inside child unless bubble rules apply).

**Rationale**: Spec assumption; `fireActivityErrorBoundary` already terminates child via `terminateChildPub`.

**Alternatives considered**:

- Propagate all child errors to Call Activity host automatically — out of scope; unchanged bubble rules.

## 7. Recover

**Decision**: No new recover redrive types. Replay `ApplyEvent` on CALL_ACTIVITY ACTIVATED payload restores boundary fields; `CompensationSubs` restored from BOUNDARY_EVENT ACTIVATED after call COMPLETED; child linkage from `called_process_instance_id`.

**Rationale**: FR-007; ledger already records boundary intents on host token.

**Alternatives considered**:

- Separate recover path for call boundaries — unnecessary if EVENT replay is complete.

## 8. Proto / API

**Decision**: **No proto changes.** Reuse `ActivityPayload` boundary fields and existing `Token` projection.

**Rationale**: Same fields as User Task waiting token; GetInstance sufficient.

**Alternatives considered**:

- New `called_process_instance_id` on boundary payload — redundant with host token.
