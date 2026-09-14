# Feature Specification: Process-Level Typed Start

**Feature Branch**: `015-process-typed-start`

**Created**: 2026-09-14

**Status**: Draft

**Input**: User description: "推进 Planned P1：Process-level typed start — message, timer, signal, conditional。Error start 仅 Event Sub-Process；进程级 error start 部署拒绝。钉死 executable process 范围后开本增量。"

## User Scenarios & Testing *(mandatory)*

Authors can already start processes with a none start (CreateInstance) or an instantiate exclusive event-based gateway. Process-level message, timer, signal, and conditional starts are still treated like none starts or left underspecified. This increment makes each typed start create an instance only when its trigger fires, rejects invalid process-level error starts, and keeps Recover coverage.

### User Story 1 - Message start creates on publish (Priority: P1)

An author models a process whose entry is a message start event bound to a message name. Deploy succeeds. CreateInstance against that deployment is rejected. Publishing that message name creates a new instance at the start event and the token continues.

**Why this priority**: Message start is the most common typed entry and reuses PublishMessage correlation.

**Independent Test**: Deploy message-start fixture; CreateInstance fails; PublishMessage creates instance and process reaches wait or end.

**Acceptance Scenarios**:

1. **Given** a process with only a message start (`order.created`) → User Task → end, **When** CreateInstance runs, **Then** it is rejected with a clear code.
2. **Given** that deployment, **When** PublishMessage with `order.created` runs, **Then** one new instance is created, the start event completes, and the instance waits on the User Task.
3. **Given** a waiting intermediate message catch on another instance for the same name, **When** PublishMessage runs, **Then** existing catch delivery still works; message start creation follows the same name matching rules without stealing an already-matching catch when both apply in one publish (document: prefer completing waiters in existing instances first; then create start instances for remaining unmatched delivery intent, or create one start instance per publish when no waiter matched—see Assumptions).

---

### User Story 2 - Timer start creates when due (Priority: P1)

An author models a process with a timer start (date, duration, or cycle as already supported for timer catches). Deploy arms the timer. When FireDue runs at or after the due time, a new instance is created and continues from the start event.

**Why this priority**: Completes the timed entry path already used mid-process and in Event Sub-Processes.

**Independent Test**: Deploy timer-start fixture; advance time / FireDue; instance exists and progresses.

**Acceptance Scenarios**:

1. **Given** a process with only a timer start → end, **When** CreateInstance runs, **Then** it is rejected.
2. **Given** that deployment and a due timer, **When** FireDue runs, **Then** one instance is created and the process completes (or reaches the next wait).
3. **Given** a cycle timer start, **When** successive due times are fired, **Then** each due firing can create another instance per cycle semantics already used for timer catches.

---

### User Story 3 - Signal start creates on publish (Priority: P1)

An author models a process with a signal start. PublishSignal with that signal name creates a new instance; CreateInstance is rejected.

**Why this priority**: Symmetric to message start; PublishSignal already exists.

**Independent Test**: Deploy signal-start fixture; PublishSignal creates instance and progresses.

**Acceptance Scenarios**:

1. **Given** a process with only a signal start → end, **When** CreateInstance runs, **Then** it is rejected.
2. **Given** that deployment, **When** PublishSignal with the bound signal name runs, **Then** one new instance is created and continues from the start event.
3. **Given** other instances waiting on the same signal catch, **When** PublishSignal runs, **Then** those waiters still complete; start creation follows the same preference as message (waiters first).

---

### User Story 4 - Conditional start creates when condition holds (Priority: P1)

An author models a process with a conditional start whose condition expression uses process variables. The instance is created when an evaluation supplies variables that make the condition true (same expression style as conditional sequence flows). CreateInstance without a true condition is rejected; a dedicated evaluation path or Publish-style trigger that supplies variables creates the instance when the condition holds.

**Why this priority**: Listed with other typed starts; closes the process-level conditional entry gap without requiring intermediate conditional catch in the same increment.

**Independent Test**: Deploy conditional-start fixture; supply variables where condition is false → no instance; supply where true → instance created and progresses.

**Acceptance Scenarios**:

1. **Given** a process with only a conditional start (`approved == true`) → end, **When** CreateInstance runs with no true condition path, **Then** it is rejected (CreateInstance does not start typed entries).
2. **Given** that deployment, **When** the conditional start is evaluated with variables `{approved: false}`, **Then** no instance is created.
3. **Given** that deployment, **When** the conditional start is evaluated with variables `{approved: true}`, **Then** one instance is created and the process completes (or reaches the next wait).

---

### User Story 5 - None start and mixed alternative starts (Priority: P1)

None start keeps CreateInstance behavior. A process may expose multiple alternative process-level starts (for example none + message). CreateInstance uses a none start when present. Message/timer/signal/conditional triggers create instances at their respective start events.

**Why this priority**: BPMN alternative starts must not break existing none-start fixtures.

**Independent Test**: None-only fixture unchanged; none+message fixture: CreateInstance enters via none; PublishMessage enters via message start.

**Acceptance Scenarios**:

1. **Given** a process with only a none start, **When** CreateInstance runs, **Then** behavior matches today’s suite (unchanged).
2. **Given** a process with a none start and a message start, **When** CreateInstance runs, **Then** the instance enters via the none start.
3. **Given** that same process, **When** PublishMessage for the message start name runs and no waiter consumes it, **Then** an instance is created entering via the message start.

---

### User Story 6 - Process-level error start rejected; ESP unchanged (Priority: P2)

An author models a process-level start with an error event definition. Deploy rejects it. Event Sub-Process error starts remain Supported as today.

**Why this priority**: Aligns with OMG (error start only in Event Sub-Process) and AGENTS Planned wording.

**Independent Test**: Deploy process with process-level error start → rejected; existing ESP error suite still passes.

**Acceptance Scenarios**:

1. **Given** a process whose only start is an error start, **When** Deploy runs, **Then** deploy is rejected with `UNSUPPORTED_ELEMENT` (or equivalent stable code).
2. **Given** existing Event Sub-Process error-start fixtures, **When** the suite runs, **Then** they still pass.

---

### User Story 7 - Recover for typed-start creation paths (Priority: P2)

After a typed start has created an instance (or after a timer start is armed at deploy), Recover then completing the same external trigger path yields the same outcome as a continuous run. Duplicate CreateInstance-style replay must not double-create incorrectly.

**Why this priority**: Completeness requires Recover coverage for creation triggers.

**Acceptance Scenarios**:

1. **Given** a message-start process, **When** PublishMessage creates an instance, then Recover, then Complete any wait, **Then** the outcome matches a continuous run.
2. **Given** a timer-start process with an armed due timer, **When** Recover then FireDue, **Then** instance creation and completion match a continuous run.
3. **Given** a signal-start process, **When** Recover after PublishSignal creation then finish waits, **Then** outcome matches continuous.

---

### Edge Cases

- Instantiate exclusive event-based gateway entry (no startEvent) remains Supported and unchanged; still incompatible with any process-level startEvent.
- Multiple process-level starts of the same kind (two message starts with different names) are each independently creatable by their trigger.
- Two message starts with the same message name: one PublishMessage creates one instance per matching start subscription (or document single-create if engine chooses one—default: one instance per matching start event definition on the deployment).
- Call Activity child processes that use typed starts: child creation follows the child’s start rules (none child via Call Activity enter; typed child only if the call path supplies the matching trigger—default assumption: Call Activity still requires a none or instantiate-EBG entry on the called process for this increment; typed-only called processes are rejected at Deploy of the caller or callee).
- Escalation start at process level remains rejected (ESP-only), same as today.
- Conditional start does not imply intermediate/boundary conditional catch (still Planned separately).
- Message/signal start creation does not require Collaboration or message flows.

## Requirements *(mandatory)*

- **FR-001**: Deploy MUST accept process-level message, timer, signal, and conditional start events with exactly one matching event definition and resolvable refs (message/signal/timer/condition) as for existing catch indexing.
- **FR-002**: Deploy MUST reject process-level error start events; Event Sub-Process error starts MUST remain accepted.
- **FR-003**: CreateInstance MUST start an instance only via a none start (or existing instantiate exclusive EBG entry). CreateInstance MUST reject deployments whose only entries are typed starts.
- **FR-004**: PublishMessage MUST create a new instance for each matching process-level message start subscription when the publish is not fully consumed by existing instance waiters (per Assumptions).
- **FR-005**: PublishSignal MUST create a new instance for each matching process-level signal start subscription under the same waiter-first rule.
- **FR-006**: FireDue MUST create a new instance when a process-level timer start is due, then continue from that start event.
- **FR-007**: Conditional start MUST create an instance when an evaluation supplies variables that make the condition true; false MUST create none.
- **FR-008**: When multiple alternative starts exist, CreateInstance MUST use a none start if present; typed triggers MUST enter via their own start event.
- **FR-009**: Recover MUST preserve armed timer starts and yield the same creation-and-progress outcome for message, timer, and signal start paths as a continuous run.
- **FR-010**: Existing none-start, instantiate-EBG, and Event Sub-Process suites MUST remain passing.

### Key Entities

- **Process-level start subscription**: Deployment-scoped binding of a typed start event (message name, signal name, timer due rule, or condition expression) used to create instances.
- **Alternative start**: One of several process-level start events; only one is taken per created instance.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Message-start fixture: CreateInstance rejected; PublishMessage creates instance and reaches User Task or end.
- **SC-002**: Timer-start fixture: FireDue creates instance and progresses.
- **SC-003**: Signal-start fixture: PublishSignal creates instance and progresses.
- **SC-004**: Conditional-start fixture: false evaluation creates none; true creates and progresses.
- **SC-005**: None-start and none+message fixtures behave as specified; existing none-start suite unchanged.
- **SC-006**: Process-level error start deploy rejected; ESP error suite still green.
- **SC-007**: Recover scenarios for message, timer, and signal start match continuous outcomes.
- **SC-008**: Existing processing regression suite remains green after the increment.

## Assumptions

- OMG: error (and escalation) start events are Event Sub-Process only; process-level error start is rejected rather than given a custom create path.
- PublishMessage / PublishSignal prefer completing waiters on existing instances; if any waiter or scope/ESP arm consumed the publish, start creation for that publish is skipped. If nothing consumed the publish, matching process-level start subscriptions each create one instance (buffered-message behavior for starts: unmatched messages may still buffer for future catches; start creation runs when the engine decides the publish is an instantiation trigger—default: on publish with no waiter match, also attempt start subscriptions before or instead of only buffering; if both a start subscription and buffering would apply, create start instance(s) and do not require a later catch).
- Conditional start evaluation reuses the same expression language as conditional sequence flows; the public trigger is an engine evaluation entry (may share or extend an existing command surface in planning).
- Called processes in this increment MUST still expose a none start or instantiate EBG for Call Activity; typed-only called processes are out of scope / rejected.
- Cycle timer start follows existing timer cycle semantics.
- No Collaboration, message flow, or choreography runtime.
- Intermediate/boundary conditional catch remains a separate Planned item.
