# Feature Specification: Link Events

**Feature Branch**: `012-link-events`

**Created**: 2026-09-09

**Status**: Draft

**Input**: User description: "引擎完整性优先；下一步 Remaining P3：Link — 支持 intermediate link throw/catch 按名称跳转（同进程内 goto）；Deploy 校验配对；Recover 覆盖。"

## User Scenarios & Testing *(mandatory)*

Authors use link events as off-page connectors: an intermediate link throw jumps execution to one or more intermediate link catches with the same link name in the same process, without a sequence flow between them.

### User Story 1 - Throw to single catch (Priority: P1)

An author models a link throw and a matching link catch (same name). When the throw is reached, the catch is entered and its outgoing path continues; the process can complete.

**Why this priority**: Core goto semantics.

**Independent Test**: Fixture throw → catch → user task → end; Complete task → completed.

**Acceptance Scenarios**:

1. **Given** a link throw and one catch with the same name, **When** Deploy runs, **Then** deploy succeeds.
2. **Given** such a process, **When** CreateInstance runs through to the catch’s outgoing User Task, **Then** the instance waits on that task.
3. **Given** a link throw with no matching catch name, **When** Deploy runs, **Then** deploy fails with UNSUPPORTED_ELEMENT.

---

### User Story 2 - One throw, multiple catches (Priority: P2)

An author attaches two link catches with the same name. The throw activates both catch paths (parallel tokens).

**Why this priority**: Matches Camunda/Zeebe multi-catch behavior.

**Independent Test**: Two catches → two ends (or join); both paths complete.

**Acceptance Scenarios**:

1. **Given** one throw and two catches with the same name, **When** the throw fires, **Then** both catch outgoing paths run.

---

### User Story 3 - Recover after link jump (Priority: P2)

Recover while waiting on a User Task reached via a link catch; Complete after Recover matches continuous run.

**Acceptance Scenarios**:

1. **Given** waiting after a link jump, **When** Recover/Open then Complete, **Then** outcome matches continuous execution.

---

### Edge Cases

- Link throw MUST NOT require outgoing sequence flows; link catch MUST NOT require incoming sequence flows.
- Link pairs are matched by `linkEventDefinition/@name` (fallback: event name) within the executable process (including nested SubProcesses).
- Link is not used as boundary, start, or end in this increment.
- Escalation/error Remaining items stay out of scope.

## Requirements *(mandatory)*

- **FR-001**: Deploy MUST accept intermediate link throw and catch events paired by name.
- **FR-002**: Deploy MUST reject a link throw with no catch of the same name in the process.
- **FR-003**: On link throw, the engine MUST enter every matching link catch and continue each catch’s outgoing flow.
- **FR-004**: Link catch enter MUST NOT wait for an external Complete (instantaneous).
- **FR-005**: Recover MUST restore waiting tokens after a link jump so subsequent Complete matches continuous outcome.
- **FR-006**: Existing suites MUST remain unchanged.

## Success Criteria *(mandatory)*

- **SC-001**: Single throw→catch fixture reaches User Task and completes after Complete.
- **SC-002**: Dual-catch fixture completes both paths.
- **SC-003**: Orphan throw is rejected at Deploy.
- **SC-004**: Recover mid-wait after link then Complete matches continuous outcome.
- **SC-005**: Existing `./processing` suites still pass.

## Assumptions

- Pairing is by link name string equality within one process definition (root + nested SubProcesses).
- No new public API; jump happens during Enter of the throw.
- Optional ledger payload may record the link name on ACTIVATED.
