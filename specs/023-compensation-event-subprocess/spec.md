# Feature Specification: Compensation Event Sub-Process

**Feature Branch**: `023-compensation-event-subprocess`

**Created**: 2026-09-15

**Status**: Draft

**Input**: User description: "继续推进 Planned：Compensation Event Sub-Process。说明与文档中不要使用 MI、EBG、ESB 等简写。"

## User Scenarios & Testing *(mandatory)*

Compensation today runs only via compensation boundary + association to an `isForCompensation` userTask or serviceTask. BPMN also allows a compensation **event sub-process** nested inside an embedded SubProcess: after that SubProcess completes, a parent-scope compensate throw that targets it (or broadcasts) starts the event sub-process instead of a boundary handler. The classic booking pattern (C.6.0) uses that event sub-process to throw nested compensate events for inner activities.

### User Story 1 - Deploy and run compensation event sub-process as SubProcess handler (Priority: P1)

An author nests a `triggeredByEvent` SubProcess with a compensate start event inside an embedded SubProcess (no sequence flows on the event sub-process). The enclosing SubProcess has no compensation boundary. After the SubProcess completes, a parent compensate throw (broadcast or `activityRef` to the SubProcess) starts the compensation event sub-process. Completing the event sub-process lets the throw path continue and the process finish.

**Independent Test**: Fixture → complete inner work → SubProcess completes → compensate throw → event sub-process runs (waiting user task inside) → Complete → process completed.

**Acceptance Scenarios**:

1. **Given** an embedded SubProcess containing a compensation event sub-process, **When** Deploy runs, **Then** Deploy succeeds.
2. **Given** that process after the SubProcess completed, **When** the parent compensate throw fires, **Then** the compensation event sub-process becomes active (its start path runs) and the throw waits until the event sub-process completes.
3. **Given** the event sub-process waiting on an inner user task, **When** that task is Completed, **Then** the event sub-process completes, the compensate throw completes, and the process reaches completed.

---

### User Story 2 - Nested compensate throws inside the compensation event sub-process (Priority: P1)

Inside the compensation event sub-process, intermediate compensate throws collect subscriptions for activities whose scope is the **enclosing compensated SubProcess** (not the event sub-process id), so inner boundary handlers (flight/hotel undo) run in reverse completion order.

**Independent Test**: SubProcess with two completed compensatable inners + compensation event sub-process that parallel-throws compensate for both → parent throw starts event sub-process → both undo handlers run → process completes.

**Acceptance Scenarios**:

1. **Given** completed inners with compensation handlers and a compensation event sub-process that throws compensate, **When** the parent throw starts the event sub-process, **Then** nested throws invoke the inner undo handlers.
2. **Given** those nested throws, **When** handlers complete, **Then** the event sub-process can finish and the outer throw path continues.

---

### User Story 3 - Recover mid compensation event sub-process (Priority: P2)

An operator recovers while the compensation event sub-process is waiting on a handler activity. After Recover, completing that activity finishes the compensation chain with the same outcome as a continuous run.

**Acceptance Scenarios**:

1. **Given** PendingCompensation with the compensation event sub-process waiting, **When** Recover finishes, **Then** the wait and pending throw are restored.
2. **Given** Recover restored that state, **When** the waiting activity is Completed, **Then** the process reaches the same completed outcome as without a crash.

---

### Edge Cases

- Compensation event sub-process MUST be nested inside an embedded SubProcess (not at process root in this increment).
- Exactly one compensation event sub-process per enclosing SubProcess; Deploy rejects multiples.
- Enclosing SubProcess MUST NOT also have a compensation boundary handler (association); Deploy rejects both.
- Compensate start MUST have empty `activityRef` in this increment (targeted start deferred); non-empty → `UNSUPPORTED_ELEMENT`.
- Compensate start is treated as non-interrupting for subscription purposes; it is not armed as a live message/timer/signal/error/escalation wait while the enclosing SubProcess runs.
- Unfinished SubProcess compensation (spec 008) remains unchanged: compensation event sub-process only becomes available after the enclosing SubProcess has completed and subscribed.
- Call Activity unfinished child compensation stays out of scope.
- Documentation spells out “event sub-process” and “multi-instance” without abbreviation.

## Requirements *(mandatory)*

- **FR-001**: Deploy MUST accept an event sub-process whose sole start event carries exactly one `compensateEventDefinition` when nested in an embedded SubProcess.
- **FR-002**: Deploy MUST reject compensation event sub-process at process root, multiple compensation event sub-processes in one SubProcess, coexisting compensation boundary on the same SubProcess, and non-empty `activityRef` on the compensate start.
- **FR-003**: When the enclosing SubProcess completes, the engine MUST subscribe compensation with the compensation event sub-process as handler (ledger-backed like boundary subscriptions).
- **FR-004**: Parent-scope compensate throw (broadcast or `activityRef` = that SubProcess) MUST queue and Enter the compensation event sub-process as the handler.
- **FR-005**: Completing the compensation event sub-process MUST AdvanceCompensation for the pending throw.
- **FR-006**: Compensate throws whose scope is a compensation event sub-process MUST collect same-scope subscriptions using the enclosing compensated SubProcess id (ParentScopeID).
- **FR-007**: Compensation event sub-process MUST NOT be armed via ordinary event sub-process start arms for message/timer/signal/error/escalation delivery.
- **FR-008**: Recover MUST restore pending compensation and a waiting compensation event sub-process (or its inner wait) so completion yields the same outcome.
- **FR-009**: Existing compensation suites (`m3_compensation`, `m4_compensate_*`, unfinished SubProcess) MUST remain green.

## Success Criteria *(mandatory)*

- **SC-001**: Simple compensation event sub-process fixture: parent throw runs event sub-process; process completes after inner Complete.
- **SC-002**: Nested-throw fixture: inner undo handlers run from throws inside the compensation event sub-process.
- **SC-003**: Recover mid-handler then Complete matches continuous-run final status.
- **SC-004**: Reject fixtures fail Deploy with `UNSUPPORTED_ELEMENT`.
- **SC-005**: Processing regression suite remains green.

## Assumptions

- Handler activities inside the compensation event sub-process may be ordinary waiting tasks (not only `isForCompensation`).
- Synthetic subscription key may reuse the compensate start event id as the compensation subscription id on the ledger.
- Condition and message languages unchanged.
