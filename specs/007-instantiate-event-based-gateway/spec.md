# Feature Specification: Instantiate Event-Based Gateway

**Feature Branch**: `007-instantiate-event-based-gateway`

**Created**: 2026-09-07

**Status**: Draft

**Input**: User description: "引擎完整性优先；下一步 Remaining P2：支持 BPMN eventBasedGateway 的 instantiate（作为流程入口/启动竞态），在已有 exclusive/parallel catch event-based gateway 之上补齐部署与启动语义，并覆盖 Recover。"

## User Scenarios & Testing *(mandatory)*

Authors already use mid-process exclusive and parallel event-based gateways that fork to intermediate catch events; the first (exclusive) or each (parallel) event continues. Today deploy rejects `instantiate="true"`, and every executable process must have a `startEvent`. This increment allows an **instantiating exclusive event-based gateway** to serve as the process entry: CreateInstance arms the catch race without a none start event, then the first matching event continues the instance.

### User Story 1 - Deploy and start via instantiate exclusive EBG (Priority: P1)

An author models a process whose only entry is an exclusive event-based gateway with `instantiate="true"`, no incoming sequence flows, and at least two outgoing intermediate catch events (for example timer and message). Deploy succeeds without a process-level none start event. CreateInstance creates an instance that immediately arms those catches (as with a mid-process exclusive EBG). Completing one catch cancels siblings and continues on that path.

**Why this priority**: Unblocks the deferred deploy reject and proves instantiate as process start using existing catch/cancel machinery.

**Independent Test**: Deploy fixture with instantiate EBG → timer + message catches. CreateInstance → two waiting catches. FireDue or PublishMessage → one path completes; sibling terminated; process ends.

**Acceptance Scenarios**:

1. **Given** a process with instantiate exclusive EBG and two catch targets and no startEvent, **When** Deploy runs, **Then** deploy succeeds.
2. **Given** that process is deployed, **When** CreateInstance runs, **Then** an active instance waits on both catch events (no none start completion required).
3. **Given** both catches are waiting, **When** the timer fires (or the message is published), **Then** that catch completes, the sibling is terminated, and the process follows that catch’s outgoing path to completion.

---

### User Story 2 - Reject invalid instantiate configurations (Priority: P1)

Deploy rejects unsafe or unsupported instantiate shapes: instantiate with incoming sequence flows; instantiate with fewer than two outgoing catches; instantiate parallel event-based gateway (deferred); instantiate together with a conflicting none startEvent entry (unless clearly defined—default: reject none startEvent when instantiate EBG is the entry).

**Why this priority**: Fail fast at deploy; keep runtime simple.

**Independent Test**: Deploy fixtures that violate rules → `UNSUPPORTED_ELEMENT` (or equivalent clear deploy error); no instance created.

**Acceptance Scenarios**:

1. **Given** instantiate EBG has an incoming sequence flow, **When** Deploy runs, **Then** deploy fails with a clear unsupported error.
2. **Given** instantiate EBG has only one outgoing catch, **When** Deploy runs, **Then** deploy fails.
3. **Given** instantiate parallel event-based gateway (`eventGatewayType="Parallel"`), **When** Deploy runs, **Then** deploy fails in this increment (deferred).

---

### User Story 3 - Recover while instantiate race is armed (Priority: P2)

An operator recovers the engine while an instance started via instantiate EBG is waiting on multiple catches. After Recover, the same events still win correctly (timer FireDue or message publish) without duplicate sibling paths.

**Why this priority**: Completeness requires Recover coverage for unfinished waiting work.

**Independent Test**: CreateInstance on instantiate fixture; Open/Recover; PublishMessage → same outcome as non-recover message-wins test.

**Acceptance Scenarios**:

1. **Given** an instance waiting on instantiate EBG catches, **When** Recover finishes, **Then** both (or all) catches remain waiting with restored timer/message/signal state.
2. **Given** Recover restored the race, **When** one event wins, **Then** siblings terminate and the process completes on the winning path.

---

### Edge Cases

- What happens when CreateInstance is called on a process that still has a none startEvent and also an instantiate EBG? Deploy MUST reject that combination in this increment (one entry style only).
- What happens to mid-process (non-instantiate) event-based gateways? Unchanged.
- What happens if PublishMessage arrives with no waiting instance? Unchanged — does not auto-create instances in this increment (CreateInstance remains the instance-creation API).
- How does StartEventID / CreateInstance choose the entry? For instantiate-only processes, CreateInstance enters the instantiate event-based gateway element as the start.
- Parallel instantiate event-based gateway (start) and multi-gateway instantiate groups are out of scope.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Deploy MUST accept exclusive event-based gateway with `instantiate="true"` when it has no incoming sequence flows and at least two outgoing flows targeting intermediate catch events.
- **FR-002**: Deploy MAY omit process-level startEvent when the process entry is exactly one valid instantiate exclusive event-based gateway.
- **FR-003**: Deploy MUST reject instantiate event-based gateway with incoming sequence flows.
- **FR-004**: Deploy MUST reject instantiate parallel event-based gateway in this increment.
- **FR-005**: Deploy MUST reject processes that combine a none startEvent entry with an instantiate event-based gateway entry.
- **FR-006**: CreateInstance on an instantiate-only process MUST create an instance and enter the instantiate event-based gateway so all configured catches become waiting.
- **FR-007**: Exclusive instantiate races MUST cancel sibling catches when one catch completes (same as mid-process exclusive EBG).
- **FR-008**: Existing non-instantiate event-based gateway behavior MUST remain unchanged.
- **FR-009**: Recover MUST restore waiting catches for instantiate-started instances so a subsequent winning event yields the same outcome as without a crash.
- **FR-010**: PublishMessage / PublishSignal / FireDue MUST NOT create a new process instance solely because an instantiate EBG is deployed (no subscription-start registry in this increment).

### Key Entities

- **Instantiate event-based gateway**: Process-entry exclusive EBG with `instantiate=true` and no incoming flows.
- **Start catch configuration**: The set of intermediate catch events targeted by that gateway’s outgoing flows.
- **Process entry**: Either a none startEvent (existing) or an instantiate EBG (this increment), not both.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Deploy of a valid instantiate exclusive EBG process (no startEvent, two catches) succeeds.
- **SC-002**: CreateInstance arms both catches; completing via timer or message finishes exactly one path and terminates the sibling.
- **SC-003**: Invalid instantiate shapes (incoming flow, parallel instantiate, startEvent+instantiate) fail at deploy.
- **SC-004**: Recover mid-race then message (or timer) win matches the non-recover outcome.
- **SC-005**: Existing mid-process event-based gateway suite continues to pass.

## Assumptions

- Mid-process exclusive and parallel event-based gateways already work (`m2_event_based_gateway*.bpmn`).
- Sparrow continues to use **CreateInstance** to allocate process instances; `instantiate` means the EBG is the process entry (BPMN start race), not that external events alone mint instances.
- Supported catch kinds for instantiate races are the same as mid-process EBG targets already allowed (timer, message, signal as currently validated).
- Parallel instantiate event-based gateway, multi-gateway instantiate groups, and receive-task targets of EBG stay out of scope.
- Live instance migration and other Remaining/Future items stay out of scope.
