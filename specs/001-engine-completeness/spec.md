# Feature Specification: Engine Completeness (next increment)

**Feature Branch**: `001-engine-completeness`

**Created**: 2026-08-26

**Status**: Draft

**Input**: User description: "使用 spec-kit 重新规划整理项目：把已完成的执行引擎能力固化为约束，把 STATUS 中的下一步整理成可独立交付的用户故事（CallActivity 独立子实例与变量映射、Error Event Sub-Process、定义版本共存），明确暂缓项。"

## User Scenarios & Testing *(mandatory)*

This increment does not restart the engine. Callers can already deploy a process, start an instance, complete waiting work, fire timers, publish messages and signals, throw errors onto waiting activities, and pull Service Task jobs. The stories below close the next gaps on the path to executable completeness.

### User Story 1 - Called process as its own instance (Priority: P1)

An author models a Call Activity that names another process in the same definitions. When an instance reaches that Call Activity, a **distinct** process instance of the called process starts. The calling instance waits at the Call Activity until the called instance completes. Operators can look up the called instance on its own (status, variables, tokens, audit trail) without reading the caller's token map as if the called work were an embedded sub-process.

The same called process MAY be referenced by more than one Call Activity. A called process MAY itself contain Call Activities (including calling back into a process already on the stack), as long as each invocation still creates a new instance rather than silently reusing the caller.

**Why this priority**: Today's Call Activity inlines the called process into the caller's instance. That blocks independent query, independent variables, and later cross-deployment calls. It is the documented design gap for Call Activity.

**Independent Test**: Deploy one definitions file with a caller process and a called process that waits on a user task. Start the caller. Confirm a second instance id exists for the called process, the caller is waiting at the Call Activity, completing the called user task finishes the called instance, and the caller then continues to its end.

**Acceptance Scenarios**:

1. **Given** a deployment that contains process `Caller` (Call Activity → end) and process `Child` (start → user task → end), **When** an instance of `Caller` is started, **Then** a second process instance of `Child` is created and `Caller` remains active waiting at the Call Activity.
2. **Given** that called instance is waiting on the user task, **When** that user task is completed, **Then** the called instance completes and the calling instance completes the Call Activity and continues.
3. **Given** two Call Activities in `Caller` that both name `Child`, **When** the caller runs, **Then** each Call Activity starts its own `Child` instance (they MUST NOT be rejected because the called process was already claimed).
4. **Given** an interrupting boundary on the Call Activity, **When** that boundary fires while the called instance is still active, **Then** the called instance is terminated and the caller continues on the boundary path.

---

### User Story 2 - Map variables into and out of the called instance (Priority: P1)

Authors declare which caller variables are copied into the called instance at start, and which called-instance variables are copied back onto the caller when the Call Activity completes. Variables that are not mapped MUST NOT leak across the instance boundary.

If the Call Activity has no mapping, the called instance starts with an empty variable set, and completing it writes no extra variables onto the caller.

**Why this priority**: Independent instances without mapping still share nothing useful. Mapping is the contract authors expect between caller and called process.

**Independent Test**: Start a caller with variables `{a:1, b:2}`, map only `a` in and `result` out. The called process writes `result`. After the Call Activity completes, the caller has `result` and still has `b`; the called instance never had `b`.

**Acceptance Scenarios**:

1. **Given** a Call Activity with an input mapping from caller `orderId` to called `id`, **When** the called instance starts, **Then** it has `id` and does not have unmapped caller variables.
2. **Given** a Call Activity with an output mapping from called `total` to caller `amount`, **When** the called instance completes, **Then** the caller variable `amount` equals the called instance's `total`.
3. **Given** a Call Activity with no mappings, **When** the called instance starts and later completes, **Then** the called instance started with no variables copied from the caller, and the caller's variables are unchanged by the call.

---

### User Story 3 - Error Event Sub-Process (Priority: P2)

Authors attach an Event Sub-Process started by an error to a process or to an embedded sub-process. When a matching error is thrown inside that scope (error end, throw error, or an uncaught error bubbling to the scope), the Event Sub-Process starts.

An interrupting error start cancels remaining work in that scope and runs the error handler instead. A non-interrupting error start runs alongside remaining work in that scope.

Error code matching follows the same rule already used for error boundaries: a start that names a code catches only that code; a start with no code is a catch-all for that scope.

**Why this priority**: Error boundaries and error ends already exist. Process-level (and nested) error Event Sub-Process is the remaining error-handling path listed as next work.

**Independent Test**: Deploy a process whose default path throws an error end, plus an interrupting Event Sub-Process with a matching error start that leads to a user task. Starting the process must terminate the default path and wait in the Event Sub-Process. Completing that user task completes the process successfully via the handler.

**Acceptance Scenarios**:

1. **Given** a process with an interrupting error Event Sub-Process for code `E1`, **When** an error end with `E1` is reached on the default path, **Then** the default path is interrupted and the Event Sub-Process is active.
2. **Given** a non-interrupting error Event Sub-Process, **When** a matching error is thrown while other work in the scope is still waiting, **Then** the handler instance of the Event Sub-Process runs and the waiting work remains.
3. **Given** an error Event Sub-Process that names code `E1` and a thrown error `E2` with no catch-all, **When** `E2` is thrown, **Then** the Event Sub-Process does not start (the error continues to bubble or terminate as it does today without an Event Sub-Process).
4. **Given** an embedded sub-process with its own error Event Sub-Process, **When** a matching error is thrown inside that sub-process, **Then** that inner handler runs and the parent process is not necessarily terminated.

---

### User Story 4 - Keep process versions side by side (Priority: P3)

Operators deploy a new revision of a process without destroying earlier revisions. New instances start on the latest revision by default. An operator MAY start an instance against an older revision that is still stored. An instance that is already running MUST continue against the revision it was started with; deploying a newer revision MUST NOT migrate or rewrite that instance.

**Why this priority**: Instance fields already carry a process version, but operators cannot yet choose among stored revisions. Live migration of running instances is a later, separate change.

**Independent Test**: Deploy process `P` (v1, waits on a user task), start an instance, deploy `P` again with a different waiting element (v2), start a second instance. The first instance still waits on v1's element; the second waits on v2's. Completing each uses that instance's own definition.

**Acceptance Scenarios**:

1. **Given** process `P` already deployed, **When** the same process id is deployed again with a changed definition, **Then** both revisions remain startable and the new default start uses the latest revision.
2. **Given** an instance started on an older revision, **When** a newer revision is deployed, **Then** that instance's waiting work and audit trail still refer to the older revision.
3. **Given** an operator names a specific stored revision at start, **When** the instance is created, **Then** it runs that revision rather than the latest.

---

### Edge Cases

- What happens when the called instance terminates (error, interrupting event) rather than completing? The Call Activity MUST NOT complete normally; interrupting boundaries on the Call Activity MAY catch the interruption; otherwise the caller MUST follow the same failure/termination path used when a waiting activity is terminated today.
- What happens when Recover runs while a called instance is still active? Both caller and called MUST rebuild; the caller MUST still be waiting at the Call Activity; completing the called instance after recover MUST still resume the caller.
- How does the system handle a mapping that names a variable that does not exist? Missing source variables MUST be treated as absent (no value copied), not as a crash. Invalid mapping structure at deploy MUST be rejected before any instance starts.
- What happens when an error Event Sub-Process and an error boundary both match? The innermost matching catch wins (boundary on the throwing activity before a parent Event Sub-Process), consistent with existing error-boundary precedence.
- How does the system handle starting against a revision that was never stored? The start MUST be rejected with a not-found style failure; no instance is created.
- What happens if the called process has no none start event? Deploy MUST reject the Call Activity.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Reaching a Call Activity MUST start a new process instance of the called process, with its own instance identity, variables, tokens, and audit trail.
- **FR-002**: The calling instance MUST wait at the Call Activity until the called instance completes or is terminated.
- **FR-003**: Completing the called instance MUST complete the Call Activity and allow the caller to continue on its outgoing flow.
- **FR-004**: Operators MUST be able to retrieve a called instance by its own instance identity (status, variables, tokens, audit trail).
- **FR-005**: A definitions file MUST allow more than one Call Activity to name the same called process.
- **FR-006**: A called process MUST be allowed to contain Call Activities (nested calls create further instances).
- **FR-007**: Call Activity input mappings MUST copy only the listed caller values into the called instance at start.
- **FR-008**: Call Activity output mappings MUST copy only the listed called values onto the caller when the Call Activity completes.
- **FR-009**: A Call Activity with no mappings MUST start the called instance with no copied caller variables and MUST NOT write called variables back.
- **FR-010**: An interrupting boundary on a Call Activity MUST terminate the called instance when that boundary completes.
- **FR-011**: An Event Sub-Process started by an error MUST be armable on a process and on an embedded sub-process.
- **FR-012**: A matching thrown error MUST start that Event Sub-Process; interrupting vs non-interrupting MUST follow the start event's interrupting flag.
- **FR-013**: Error code matching for those starts MUST follow the same code vs catch-all rule as error boundaries.
- **FR-014**: Deploying a process again MUST keep prior revisions startable; new starts default to the latest revision.
- **FR-015**: An already-running instance MUST stay bound to the revision it was started with.
- **FR-016**: Operators MUST be able to start an instance against a specific stored revision.
- **FR-017**: Existing waiting-work completion, job pull, timer, message, signal, compensation, and error-boundary behavior MUST remain available on both caller and called instances.

### Key Entities

- **Deployment**: An immutable snapshot of one BPMN definitions file, identified for start and recovery. After this increment it MAY represent one revision among several for the same process id.
- **Process instance**: A running or completed execution with identity, status, variables, tokens, and an audit trail. A Call Activity creates a second instance rather than borrowing the caller's.
- **Call relationship**: A link from a waiting Call Activity (caller instance + activity + token) to the called instance it started.
- **Variable mapping**: Named copies of values at Call Activity start (in) and completion (out).
- **Error Event Sub-Process**: An event-triggered sub-process whose start is an error; armed while its parent scope is active.
- **Process revision**: A stored edition of a process id; instances bind to one revision for their lifetime.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: After a Call Activity activates, an operator can identify and inspect the called work as its own process instance (not only as tokens inside the caller) in a single lookup.
- **SC-002**: In a mapped call, unmapped caller variables are absent from the called instance 100% of the time, and mapped outputs appear on the caller after the call completes.
- **SC-003**: A definitions file with two Call Activities naming the same child process deploys and runs both calls to completion without a deploy-time uniqueness rejection.
- **SC-004**: An interrupting error Event Sub-Process takes over from a matching error end on the default path; the process can then complete through the handler without remaining default-path work.
- **SC-005**: After a second revision of the same process is deployed, a previously started instance still completes against its original waiting elements; a newly started instance uses the new revision.
- **SC-006**: Recover after a crash mid-call restores both instances so that completing the called instance still resumes the caller (no stranded Call Activity).

## Assumptions

- M1–M4c capabilities (start, user/service tasks, gateways, embedded sub-process, event sub-process for message/timer/signal, timers, messages, signals, compensation, error boundary/end/throw, in-definition Call Activity as it exists today) remain the baseline and are not re-specified here except where a story changes Call Activity identity.
- This increment's Call Activity still resolves `calledElement` **inside the same definitions**. Cross-deployment called processes stay deferred.
- Live migration of running instances onto a newer revision is out of scope; versions coexist, instances do not move.
- Cluster, product-suite UI, Incident records, instantiate event-based gateway, compensation into an unfinished sub-process, Event Sub-Process nested inside another Event Sub-Process, three or more waiting boundaries of the **same** kind on one activity, and remaining element types (escalation, link, conditional, terminate, send/receive/manual/business-rule task, multi-instance) stay deferred.
- Missing mapping sources copy nothing; they do not fail the instance.
- Called instances started without input mapping have an empty variable set.
- Default start after a new deploy uses the latest revision of the targeted process; specifying a revision is optional.
- The existing single-node, one-command-at-a-time-per-instance execution model remains; a called instance is a second instance with its own command queue, and resuming the caller happens after the called instance has settled.
