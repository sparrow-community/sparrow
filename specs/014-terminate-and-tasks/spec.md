# Feature Specification: Terminate End and Remaining Tasks

**Feature Branch**: `014-terminate-and-tasks`

**Created**: 2026-09-11

**Status**: Draft

**Input**: User description: "引擎完整性 Remaining P3 五项一起实现：Terminate End；Send Task；Receive Task；Manual Task；Business Rule Task。每项独立可测；Recover 覆盖等待类任务。"

## User Scenarios & Testing *(mandatory)*

Authors still cannot deploy terminate ends or the four remaining task types; Deploy rejects them or treats terminate as a none end. This increment closes that completeness gap.

### User Story 1 - Terminate end cancels remaining work (Priority: P1)

An author models a parallel split: one branch waits on a User Task, the other reaches a terminate end. When the terminate end is reached, remaining waiting work is cancelled and the process finishes successfully (not as an unhandled error).

**Why this priority**: Terminate is a distinct end semantics already parsed but not executed.

**Independent Test**: Parallel fixture User Task vs terminate end; start instance; User Task is terminated; process completed.

**Acceptance Scenarios**:

1. **Given** a process with a parallel fork to a User Task and a terminate end, **When** an instance is started, **Then** the User Task is terminated and the process is completed.
2. **Given** a terminate end inside an embedded SubProcess with a sibling waiting User Task, **When** the instance is started, **Then** the inner User Task is terminated, the SubProcess completes, and the parent process completes.
3. **Given** a none end (no terminate definition), **When** a parallel User Task is still waiting, **Then** the process does not complete (existing join/end behavior unchanged).

---

### User Story 2 - Manual Task waits for Complete (Priority: P1)

An author models a Manual Task. The instance waits until an operator completes that element, then continues to the end.

**Why this priority**: Same wait/complete story as User Task with a distinct element type.

**Independent Test**: start → wait at Manual Task → Complete → process completed.

**Acceptance Scenarios**:

1. **Given** start → Manual Task → end, **When** the instance is created, **Then** it waits on the Manual Task (type manual task).
2. **Given** that wait, **When** Complete runs, **Then** the process completes.

---

### User Story 3 - Receive Task waits for a message (Priority: P1)

An author models a Receive Task bound to a message name. The instance waits until that message is published, then continues.

**Why this priority**: Completes the receive side of message tasks without collapsing to a generic catch event.

**Independent Test**: start → wait at Receive Task → PublishMessage → process completed.

**Acceptance Scenarios**:

1. **Given** a Receive Task with message name `order.confirmed`, **When** the instance is created, **Then** it waits with that message name.
2. **Given** that wait, **When** PublishMessage with `order.confirmed` runs, **Then** the Receive Task completes and the process completes.
3. **Given** a Receive Task with `instantiate=true`, **When** Deploy runs, **Then** deploy is rejected.

---

### User Story 4 - Send Task publishes then continues (Priority: P1)

An author models a Send Task that names a message. Reaching it publishes that message (same delivery as a message throw) and the token continues without an external Complete.

**Why this priority**: Completes the send side of message tasks.

**Independent Test**: Process A send + process B receive (or catch) of the same name; starting A delivers to B.

**Acceptance Scenarios**:

1. **Given** start → Send Task (`order.confirmed`) → end, **When** the instance is started, **Then** the Send Task completes in the same command and the process completes.
2. **Given** a second instance waiting on that message name, **When** the Send Task runs, **Then** the waiting instance is completed via the published message.

---

### User Story 5 - Business Rule Task is a job (Priority: P1)

An author models a Business Rule Task. The engine waits for a job worker (same Activate / Complete / Fail / incident path as Service Task). Decision logic stays outside the engine.

**Why this priority**: No embedded DMN; job worker is the existing extension point.

**Independent Test**: start → wait job → Activate by job type → Complete → process completed.

**Acceptance Scenarios**:

1. **Given** a Business Rule Task with implementation `decide.v1`, **When** the instance is created, **Then** Activate for `decide.v1` returns that job.
2. **Given** that job, **When** Complete runs, **Then** the process completes.

---

### User Story 6 - Recover waiting tasks (Priority: P2)

An operator recovers while waiting on Manual, Receive, or Business Rule Task. After Recover, the same Complete / PublishMessage / Activate+Complete yields the same outcome as a continuous run.

**Why this priority**: Completeness definition requires Recover coverage for waiting work.

**Acceptance Scenarios**:

1. **Given** a wait on Manual Task, **When** Recover then Complete, **Then** the process completes.
2. **Given** a wait on Receive Task, **When** Recover then PublishMessage, **Then** the process completes.
3. **Given** a wait on Business Rule Task, **When** Recover then Activate and Complete, **Then** the process completes.

---

### Edge Cases

- Terminate end in a SubProcess cancels only that SubProcess’s remaining work; the parent continues as after a normal SubProcess complete.
- Terminate end completes the enclosing process instance successfully (status completed), unlike unhandled error (status terminated).
- Active Call Activity children in a cancelled scope are terminated.
- Abstract `task` (no specialized type) remains unsupported.
- Script Task, Complex Gateway, Transaction/cancel, instantiate Receive Task, and process-level typed start stay out of scope.
- Multi-instance on the new task types follows the same loop rules as User Task (Manual/Receive) or Service Task (Business Rule/Send) when `multiInstanceLoopCharacteristics` is present.
- Mixed event definitions on an end (terminate plus error/message/etc.) are rejected at Deploy.

## Requirements *(mandatory)*

- **FR-001**: Deploy MUST accept terminate end events that have exactly one `terminateEventDefinition` and no other event definitions.
- **FR-002**: Reaching a terminate end MUST cancel remaining tokens in the enclosing scope (process or SubProcess), then complete that scope.
- **FR-003**: A process-level terminate end MUST complete the instance successfully after cancelling remaining work.
- **FR-004**: Deploy MUST accept Manual, Receive, Send, and Business Rule tasks; abstract `task` MUST still be rejected.
- **FR-005**: Manual Task MUST wait and complete via the existing Complete command, recording `TYPE_MANUAL_TASK`.
- **FR-006**: Receive Task MUST wait for PublishMessage using its message name (`messageRef` resolved, else task name, else id).
- **FR-007**: Receive Task with instantiate MUST be rejected at Deploy.
- **FR-008**: Send Task MUST publish its message name and continue without waiting for Complete.
- **FR-009**: Business Rule Task MUST wait as a job (job type from implementation, else name, else id), using Activate/Complete/Fail/incident like Service Task.
- **FR-010**: Recover MUST restore Manual, Receive, and Business Rule waits so subsequent work matches a continuous run.
- **FR-011**: Existing suites MUST remain passing.

## Success Criteria *(mandatory)*

- **SC-001**: Parallel User Task vs terminate end: User Task terminated; process completed.
- **SC-002**: Terminate inside SubProcess: inner wait cancelled; parent completes.
- **SC-003**: Manual Task fixture completes after Complete.
- **SC-004**: Receive Task fixture completes after PublishMessage.
- **SC-005**: Send Task fixture publishes and completes in one start; a waiter on the same name is delivered.
- **SC-006**: Business Rule Task is claimed by Activate and completes after Complete.
- **SC-007**: Recover mid-wait then finish matches continuous outcome for Manual, Receive, and Business Rule.
- **SC-008**: `go test ./processing/` passes.

## Assumptions

- No new public APIs; reuse Complete, PublishMessage, Job Activate/Fail.
- Business Rule does not execute DMN inside the engine.
- Message names resolve like existing throw/catch (`messageRef` → definitions message name).
- Compensation handlers remain User Task or Service Task unless already linked; new types MAY host interrupting boundaries.
- Live instance migration, Script Task, and Complex Gateway stay deferred.
