# Feature Specification: Call Activity boundary and compensation parity

**Feature Branch**: `005-callactivity-boundary-compensation`

**Created**: 2026-09-02

**Status**: Draft

**Input**: User description: "Engine increment 005 — CallActivity boundary & compensation parity with SubProcess: deploy and run timer/message/signal/error/compensation boundaries on Call Activity; interrupting paths terminate child instance; compensation subscribes when call completes; Recover preserved. Next P2 per AGENTS.md Remaining list."

## User Scenarios & Testing *(mandatory)*

Call Activity already spawns a separate child instance (including cross-deployment), supports IO mapping, and the runtime can terminate a child when an interrupting boundary fires — but **deploy rejects** boundaries attached to a Call Activity today (`validateBoundaryHost` allows only userTask, serviceTask, subProcess). SubProcess hosts get timer/message/signal boundaries armed on enter, compensation subscribed on complete, and error boundaries per existing rules. Authors cannot model the same control-flow around a called process as around an embedded SubProcess.

### User Story 1 - Interrupting boundaries on Call Activity (Priority: P1)

An author attaches an **interrupting** timer, message, signal, or error boundary to a Call Activity. When the caller instance is waiting at the Call Activity and the boundary fires, the engine terminates the active child instance, completes the Call Activity on the boundary path, and the caller continues along the boundary outgoing flow.

**Why this priority**: Without interrupting boundaries, long-running or cancellable called flows cannot be aborted from the caller — a common orchestration pattern (timeout, cancel message, error escalation).

**Independent Test**: Deploy caller with Call Activity + interrupting timer boundary; called process waits on user task. Fire timer → child instance terminated, caller token leaves Call Activity on boundary path, audit shows boundary COMPLETED and child PROCESS TERMINATED.

**Acceptance Scenarios**:

1. **Given** an interrupting timer boundary on a Call Activity with an active child, **When** the timer fires, **Then** the child instance is terminated and the caller follows the boundary outgoing flow.
2. **Given** an interrupting message boundary on a Call Activity, **When** a matching message is published to the caller instance, **Then** the child is terminated and the caller follows the boundary path.
3. **Given** an interrupting signal boundary on a Call Activity, **When** a matching signal is published, **Then** the child is terminated and the caller follows the boundary path.
4. **Given** an interrupting error boundary on a Call Activity, **When** a matching error is thrown at the Call Activity host token, **Then** the child is terminated and the caller follows the boundary path.
5. **Given** a cross-deployment Call Activity with an interrupting boundary, **When** the boundary fires, **Then** behavior matches same-deployment (child terminated via callee deployment context).

---

### User Story 2 - Non-interrupting boundaries while the call waits (Priority: P1)

An author attaches a **non-interrupting** timer or message boundary to a Call Activity. While the caller waits and the child is still active, the boundary may fire without terminating the child or the Call Activity host token; the spawned boundary path runs in parallel until the call completes or is interrupted.

**Why this priority**: Parity with SubProcess and existing non-interrupting boundary semantics on tasks; needed for side-effect paths (notifications, SLA logging) during a call.

**Independent Test**: Call Activity with non-interrupting message boundary + active child. Publish message → boundary path token appears; child still active; complete child → Call Activity completes and boundary path tokens are cleaned up per existing non-interrupting rules.

**Acceptance Scenarios**:

1. **Given** a non-interrupting timer on Call Activity, **When** timer fires while child is active, **Then** boundary path proceeds and child remains active.
2. **Given** a non-interrupting message on Call Activity, **When** message delivered while child is active, **Then** boundary path proceeds and child remains active.
3. **Given** non-interrupting boundary path still open, **When** child completes normally, **Then** Call Activity completes and non-interrupting boundary tokens are disarmed or completed per existing engine rules.

---

### User Story 3 - Compensation boundary on Call Activity (Priority: P2)

An author attaches a compensation boundary to a Call Activity with a compensation handler activity. When the Call Activity **completes successfully** (child finished, caller received output mapping), the engine arms a compensation subscription for that call — same as SubProcess. A later compensate throw in the same scope runs the handler.

**Why this priority**: Sagas and undo flows often need to compensate a completed call, not only inline SubProcess work.

**Independent Test**: Caller with Call Activity + compensation boundary → handler user task. Run call to completion → compensation boundary ACTIVATED in projection. Throw compensate → handler runs.

**Acceptance Scenarios**:

1. **Given** a compensation boundary on Call Activity, **When** the call completes (child done, Call Activity COMPLETED), **Then** compensation subscription is recorded on the caller instance.
2. **Given** an armed compensation subscription for a completed call, **When** a compensate throw occurs in scope, **Then** the linked handler activity runs.
3. **Given** Call Activity terminated by interrupting boundary before normal completion, **When** compensate throw occurs, **Then** no compensation subscription exists for that call (same as interrupted SubProcess).

---

### User Story 4 - Recover with Call Activity boundaries (Priority: P2)

After a crash, Recover rebuilds a caller waiting at Call Activity with armed timer/message/signal boundaries, active child linkage, and compensation subscriptions for completed calls.

**Why this priority**: Completeness requires boundary and compensation state to be ledger-derived, not memory-only.

**Independent Test**: Reach waiting call with timer boundary armed; Recover → same boundary due time and child id; fire timer after recover → same outcome as before crash.

**Acceptance Scenarios**:

1. **Given** caller waiting at Call Activity with armed timer boundary and active child, **When** Recover runs, **Then** boundary fields and `called_process_instance_id` match pre-crash projection.
2. **Given** completed Call Activity with compensation subscription, **When** Recover runs, **Then** subscription remains; compensate throw still runs handler.

---

### Edge Cases

- Multiple boundary **kinds** on one Call Activity: one per kind (timer, message, signal, error, compensation) — same rule as today for other hosts; multiple same-kind boundaries stay a separate increment.
- Boundary fires while child start is still in flight: boundary MUST NOT leave caller in ambiguous double-active state; interrupting boundary waits for or cancels in-flight child start consistently.
- Error boundary on Call Activity vs error thrown inside child: error boundary on the **host** catches errors thrown **at** the Call Activity token; errors inside the child instance follow child-scope rules (unchanged).
- Compensation handler inside callee vs caller: compensation boundary and handler are on the **caller** deployment; handler runs in caller instance scope.
- Cross-deployment call with compensation: subscription and handler use caller deployment; child termination on interrupt uses callee deployment id from child instance.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Deploy MUST accept timer, message, signal, error, and compensation boundaries attached to a Call Activity element (parity with SubProcess host rules).
- **FR-002**: On Call Activity ACTIVATED (waiting), the engine MUST arm attached timer, message, and signal boundaries on the host token using the same rules as SubProcess.
- **FR-003**: When an interrupting boundary fires on a waiting Call Activity, the engine MUST terminate the linked child instance (if any) and MUST complete the Call Activity on the boundary path.
- **FR-004**: When a non-interrupting boundary fires on a waiting Call Activity, the engine MUST NOT terminate the child; boundary path semantics MUST match existing non-interrupting behavior on SubProcess.
- **FR-005**: On normal Call Activity completion, the engine MUST arm compensation subscription when a compensation boundary is attached (same as SubProcess COMPLETED path).
- **FR-006**: Compensate throw MUST invoke handlers linked to completed Call Activities within the same scope ordering rules as today.
- **FR-007**: Recover MUST rebuild waiting Call Activity boundary state, child linkage, and compensation subscriptions from the event log.
- **FR-008**: Existing same-deployment and cross-deployment Call Activity behavior (IO mapping, child deployment id, NOT_FOUND on missing callee) MUST NOT regress.

### Key Entities

- **Call Activity host token**: Caller token waiting at Call Activity; carries boundary ids, due times, child instance id.
- **Child instance**: Callee process instance; terminated when interrupting boundary fires on host.
- **Compensation subscription**: Boundary ACTIVATED after Call Activity COMPLETED; links to handler activity id.
- **Boundary path token**: Spawned token on interrupting or non-interrupting boundary outgoing flow.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Deploy succeeds for Call Activity with each supported boundary kind (timer, message, signal, error, compensation) in one test pass.
- **SC-002**: Interrupting timer/message/signal on waiting Call Activity terminates child and reaches boundary end event in 100% of tested cases.
- **SC-003**: Non-interrupting boundary fires while child active without terminating child in 100% of tested cases.
- **SC-004**: Completed Call Activity arms compensation; subsequent compensate throw runs handler in tested fixture.
- **SC-005**: Recover after mid-call wait restores boundary + child linkage; post-recover boundary fire matches non-recover outcome.
- **SC-006**: Existing `m4_call_activity*` and `m8_cross_call*` tests continue to pass unchanged.

## Assumptions

- One boundary per kind per Call Activity (existing engine rule); multiple same-kind boundaries deferred to separate increment.
- Error boundary on Call Activity host catches errors propagated to the host token, not arbitrary errors inside the child process unless existing error-bubble rules already apply.
- Compensation handler activities remain userTask or serviceTask on the caller definition.
- Multi-instance Call Activity boundaries deferred to the Multi-instance Call Activity increment.
- Live instance migration out of scope.
