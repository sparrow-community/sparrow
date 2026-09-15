# Feature Specification: Instantiate Receive Task and Event-Based Gateway Entries

**Feature Branch**: `022-instantiate-receive-event-gateway`

**Created**: 2026-09-15

**Status**: Draft

**Input**: User description: "继续推进 Planned：Instantiate Receive Task；parallel instantiate event-based gateway；event-based gateway targeting Receive Task；multiple instantiate event-based gateway entries。说明与文档中不要使用 MI、EBG、ESB 等简写。"

## User Scenarios & Testing *(mandatory)*

Exclusive instantiate event-based gateway as process entry already works. This increment closes four related entry/target gaps.

### User Story 1 - Instantiate receive task as process entry (Priority: P1)

An author models a process with no start event and a receive task with `instantiate="true"` and no incoming sequence flow. Deploy succeeds. CreateInstance leaves the instance waiting on that receive task. PublishMessage for its message completes the task and the process continues.

**Independent Test**: Deploy fixture → CreateInstance → wait on receive → PublishMessage → process completed.

**Acceptance Scenarios**:

1. **Given** instantiate receive task as sole entry, **When** Deploy and CreateInstance run, **Then** the instance waits on that receive task.
2. **Given** that wait, **When** PublishMessage matches the task message, **Then** the process completes on the receive outgoing path.
3. **Given** instantiate receive combined with a start event or with incoming flows, **When** Deploy runs, **Then** Deploy fails with `UNSUPPORTED_ELEMENT`.

---

### User Story 2 - Parallel instantiate event-based gateway (Priority: P1)

An author uses `eventGatewayType="Parallel"` with `instantiate="true"` as the process entry. CreateInstance arms all catch targets. Completing one catch does not cancel siblings; each path can complete independently.

**Acceptance Scenarios**:

1. **Given** parallel instantiate event-based gateway with timer and message catches, **When** CreateInstance then FireDue and PublishMessage both succeed, **Then** both catch paths complete without sibling termination.

---

### User Story 3 - Event-based gateway targets receive task (Priority: P1)

A mid-process exclusive event-based gateway may target a receive task (and intermediate catch events). The first completing target cancels sibling waits.

**Acceptance Scenarios**:

1. **Given** exclusive event-based gateway → receive task + message catch, **When** PublishMessage completes the receive task, **Then** the sibling catch is terminated and the process follows the receive path.

---

### User Story 4 - Multiple exclusive instantiate event-based gateways (Priority: P2)

A process may have more than one exclusive instantiate event-based gateway as alternative entries. CreateInstance arms all of their catch targets. The first winning catch cancels waits belonging to the other instantiate entries as well.

**Acceptance Scenarios**:

1. **Given** two exclusive instantiate event-based gateways each with two catches, **When** CreateInstance runs, **Then** all four catches wait.
2. **Given** that race, **When** one message wins, **Then** the other three waits terminate and the process continues on the winning path.

---

### Edge Cases

- Parallel instantiate does not cancel siblings on the same gateway.
- Instantiate entries remain incompatible with process-level start events in this increment.
- PublishMessage still does not mint instances solely because an instantiate receive or gateway is deployed (CreateInstance remains the allocator).
- Call Activity callees may use instantiate receive or instantiate event-based gateway entries the same way CreateInstance does.

## Requirements *(mandatory)*

- **FR-001**: Deploy MUST accept instantiate receive task with no incoming flows when it is a valid process entry (no start event).
- **FR-002**: Deploy MUST accept parallel instantiate event-based gateway as process entry.
- **FR-003**: Deploy MUST allow event-based gateway outgoing targets to be intermediate catch events or receive tasks.
- **FR-004**: Deploy MUST accept multiple exclusive instantiate event-based gateways as alternative entries.
- **FR-005**: CreateInstance MUST arm all instantiate entries for the process.
- **FR-006**: Exclusive races (including cross-entry alternatives) MUST cancel remaining instantiate waits when one wait completes.
- **FR-007**: Receive task completion MUST cancel exclusive event-based gateway siblings when it is a gateway target.
- **FR-008**: Existing exclusive instantiate and mid-process suites MUST remain green.

## Success Criteria *(mandatory)*

- **SC-001**: Instantiate receive fixture CreateInstance → PublishMessage → completed.
- **SC-002**: Parallel instantiate fixture both catches can complete without mutual cancel.
- **SC-003**: Gateway-to-receive fixture cancels sibling on receive win.
- **SC-004**: Multiple instantiate gateways arm all catches; one win clears the rest.
- **SC-005**: Processing regression suite remains green.

## Assumptions

- Condition and message languages unchanged.
- Documentation and AGENTS updates spell out “event-based gateway” and “event sub-process” without abbreviation.
