# Feature Specification: Message End and Signal End

**Feature Branch**: `016-message-signal-end`

**Created**: 2026-09-15

**Status**: Draft

**Input**: User description: "继续推进 Planned P1：Message end；signal end。到达 end 时发布 message/signal（同 intermediate throw / Send Task），然后完成作用域；Deploy 接受；Recover 覆盖跨实例投递场景。"

## User Scenarios & Testing *(mandatory)*

Authors can throw messages and signals from intermediate throws and Send Tasks, and catch them mid-process. End events that throw a message or signal are still rejected at Deploy as unsupported. This increment accepts those ends and publishes on enter before completing the scope.

### User Story 1 - Message end publishes then completes (Priority: P1)

An author models start → User Task on process A waiting for a message, and process B that ends with a message end of the same name. When B reaches the message end, the message is published (same delivery as a message throw), A continues, and B completes successfully.

**Why this priority**: Message end is the natural counterpart of message start/throw for process completion paths.

**Independent Test**: Deploy catch + message-end fixtures; start catch waiter; start ender; waiter completes; ender completed.

**Acceptance Scenarios**:

1. **Given** a process with only a message end (`order.done`), **When** CreateInstance runs, **Then** the end publishes `order.done` and the process completes in the same command.
2. **Given** another instance waiting on message catch `order.done`, **When** the message-end process runs, **Then** the waiting instance is delivered and continues.
3. **Given** an end with message plus terminate (or other mixed defs), **When** Deploy runs, **Then** deploy is rejected.

---

### User Story 2 - Signal end publishes then completes (Priority: P1)

An author models a signal end. Reaching it broadcasts the signal (same as signal throw) and completes the process.

**Why this priority**: Symmetric to message end; PublishSignal path already exists.

**Independent Test**: Signal catch waiter + signal-end process; start both; catch completes; ender completed.

**Acceptance Scenarios**:

1. **Given** a process with only a signal end (`done`), **When** CreateInstance runs, **Then** the signal is published and the process completes.
2. **Given** a waiting signal catch for `done`, **When** the signal-end process runs, **Then** the waiter completes.
3. **Given** mixed signal+other event definitions on an end, **When** Deploy runs, **Then** deploy is rejected.

---

### User Story 3 - SubProcess message/signal end (Priority: P2)

A message or signal end inside an embedded SubProcess publishes, completes the SubProcess, and the parent continues.

**Why this priority**: Same end semantics inside scopes as process-level.

**Independent Test**: Parent → SubProcess → message/signal end → parent end; with optional external waiter.

**Acceptance Scenarios**:

1. **Given** an embedded SubProcess that ends with a message end, **When** the instance runs, **Then** the message is published, the SubProcess completes, and the parent completes.
2. **Given** the same for a signal end, **When** the instance runs, **Then** the signal is published and scopes complete.

---

### User Story 4 - Recover after message/signal end delivery (Priority: P2)

After a message or signal end has created delivery into a waiting catch, Recover then finishing remaining waits matches a continuous run.

**Why this priority**: Completeness requires Recover coverage for publish-on-end paths.

**Acceptance Scenarios**:

1. **Given** a waiter on message catch and a message-end process that has published, **When** Recover then Complete remaining work, **Then** outcomes match continuous.
2. **Given** signal end delivery into a waiter, **When** Recover then finish, **Then** outcomes match continuous.

---

### Edge Cases

- None end, error end, escalation end, terminate end, compensate end unchanged.
- Message/signal name resolution matches intermediate throw (`messageRef`/`signalRef` → definitions name, else end name, else id).
- Message end may also create process-level message-start instances when no waiter matches (existing PublishMessage rules).
- Collaboration / message flows remain Excluded.
- Implicit throw / empty throw alias stays Planned separately.

## Requirements *(mandatory)*

- **FR-001**: Deploy MUST accept end events with exactly one `messageEventDefinition` and no other event definitions.
- **FR-002**: Deploy MUST accept end events with exactly one `signalEventDefinition` and no other event definitions.
- **FR-003**: Reaching a message end MUST publish the resolved message name (same delivery as message throw / Send Task) and then complete the enclosing scope.
- **FR-004**: Reaching a signal end MUST publish the resolved signal name (same delivery as signal throw) and then complete the enclosing scope.
- **FR-005**: Mixed event definitions on an end (message/signal plus terminate/error/etc.) MUST be rejected at Deploy.
- **FR-006**: Existing none/error/escalation/terminate/compensate end suites MUST remain passing.
- **FR-007**: Recover MUST yield the same outcome as continuous for message/signal end delivery into waiters.

### Key Entities

- **Message end**: End event indexed with a message name; publishes on enter.
- **Signal end**: End event indexed with a signal name; publishes on enter.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Message-end fixture completes and delivers to a waiting message catch.
- **SC-002**: Signal-end fixture completes and delivers to a waiting signal catch.
- **SC-003**: SubProcess message/signal end publishes and parent completes.
- **SC-004**: Mixed-definition ends rejected at Deploy.
- **SC-005**: Recover scenarios match continuous outcomes.
- **SC-006**: Existing processing regression suite remains green.

## Assumptions

- No new public APIs; reuse PublicationMessage / PublicationSignal flush paths.
- Name resolution identical to intermediate throw.
- Message/signal end does not wait; publish is fire-and-continue then scope complete.
- Process-level message starts may be created by unmatched publish from a message end (existing 015 rules).
