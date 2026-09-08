# Feature Specification: Escalation Events

**Feature Branch**: `011-escalation-events`

**Created**: 2026-09-08

**Status**: Draft

**Input**: User description: "引擎完整性优先；下一步 Remaining P3：Escalation — 支持 escalation 抛出/捕获（中间抛出、结束、边界、事件子流程）；未捕获不失败实例；Recover 覆盖。"

## User Scenarios & Testing *(mandatory)*

Authors need BPMN escalation to signal a parent scope without treating the case as a fatal error. Escalation throw/end must be catchable by escalation boundaries and escalation Event Sub-Processes; an uncaught escalation must not terminate the process instance.

### User Story 1 - Intermediate throw + interrupting boundary (Priority: P1)

An author models an escalation intermediate throw inside an embedded SubProcess and an interrupting escalation boundary on that SubProcess. Deploy succeeds. When the throw fires, the SubProcess is cancelled, the boundary path completes, and the thrower’s inner path does not continue.

**Why this priority**: Core throw/catch pair mirrors error boundaries with escalation semantics.

**Independent Test**: Fixture with SP + throw + interrupting boundary; run to boundary end; SP host TERMINATED.

**Acceptance Scenarios**:

1. **Given** a SubProcess with an escalation intermediate throw and an interrupting escalation boundary on the SubProcess, **When** the instance reaches the throw, **Then** the boundary path completes and the SubProcess is terminated.
2. **Given** the same model with a mismatched escalation code on the boundary, **When** the throw fires, **Then** the escalation is uncaught at that boundary and the throw continues (or completes) without failing the instance.

---

### User Story 2 - Escalation end + Event Sub-Process (Priority: P1)

An author uses an escalation end event and an interrupting escalation Event Sub-Process in the same scope. The end throws; the Event Sub-Process starts and handles the escalation.

**Why this priority**: Escalation end + ESP is a standard Camunda-style pattern.

**Independent Test**: Fixture like existing error ESP but with escalation; assert ESP start and handler task waiting; process does not go to StatusTerminated solely due to the throw.

**Acceptance Scenarios**:

1. **Given** an escalation end and matching interrupting escalation ESP, **When** the end is reached, **Then** the ESP is entered and the parent scope tokens are terminated per interrupting rules.
2. **Given** a non-interrupting escalation ESP, **When** escalation is thrown on a parallel path, **Then** the main path can still complete while the ESP runs.

---

### User Story 3 - Uncaught escalation does not fail the instance (Priority: P1)

An author throws an escalation with no catcher. The throw completes and the process continues or completes normally; the instance is not terminated as a failure.

**Why this priority**: Primary semantic difference from error.

**Independent Test**: Process with only escalation intermediate throw then end; status completed.

**Acceptance Scenarios**:

1. **Given** an escalation intermediate throw with no catcher, **When** it fires, **Then** outgoing continues and the process can complete.
2. **Given** an escalation end with no catcher, **When** it fires, **Then** the end completes and the process may complete; instance is not StatusTerminated due to the escalation alone.

---

### User Story 4 - Recover (Priority: P2)

Recover mid-wait after an escalation path left a User Task waiting (e.g. ESP handler task); Complete after Recover matches a continuous run.

**Acceptance Scenarios**:

1. **Given** an interrupting escalation ESP waiting on a handler User Task, **When** Recover/Open runs then Complete, **Then** outcome matches continuous execution.

---

### Edge Cases

- Non-interrupting escalation boundary on SubProcess: host continues; boundary path runs on a spawned token.
- Catch-all escalation (empty code) matches any thrown code; prefer exact code over catch-all.
- Multiple escalation boundaries with different codes on one activity (same as error).
- Link, Terminate, conditional flow, and task-type Remaining items stay out of scope.
- Intermediate catch escalation (standalone) is out of scope for this increment (boundary + ESP only).

## Requirements *(mandatory)*

- **FR-001**: Deploy MUST accept escalation intermediate throw and escalation end events with resolved escalation codes.
- **FR-002**: Deploy MUST accept escalation boundaries (interrupting and non-interrupting) and escalation Event Sub-Process starts.
- **FR-003**: Throwing an escalation MUST search the current scope for an armed escalation Event Sub-Process, then an escalation boundary on the enclosing activity/scope, then bubble to parent scopes (same structural walk as error).
- **FR-004**: An uncaught escalation MUST NOT terminate the process instance; the thrower MUST continue (intermediate throw) or complete normally (escalation end).
- **FR-005**: Interrupting catch MUST terminate the affected scope/host per existing interrupting boundary / ESP rules; non-interrupting catch MUST leave the host running and start the catch path concurrently.
- **FR-006**: Recover MUST restore waiting tokens after an escalation catch so subsequent Complete matches continuous outcome.
- **FR-007**: Existing error suites MUST remain unchanged.

## Success Criteria *(mandatory)*

- **SC-001**: Intermediate throw inside SubProcess caught by interrupting escalation boundary completes the boundary path.
- **SC-002**: Escalation end starts a matching interrupting Event Sub-Process.
- **SC-003**: Uncaught intermediate throw completes the process without StatusTerminated.
- **SC-004**: Recover then Complete on an escalation ESP handler matches continuous outcome.
- **SC-005**: Existing `./processing` error and boundary suites still pass.

## Assumptions

- Escalation codes resolve from `<escalation escalationCode="…">` (or name/id fallback), parallel to error codes.
- Ledger records an escalation-thrown intent with an escalation code payload (distinct from error).
- No new public ThrowEscalation API on waiting activities in this increment; throws come from throw/end events.
- BPMN: uncaught escalation is ignored for instance failure (unlike uncaught error).
