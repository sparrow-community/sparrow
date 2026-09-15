# Feature Specification: MI Complex Behavior and Behavior Event Refs

**Feature Branch**: `021-mi-complex-behavior`

**Created**: 2026-09-15

**Status**: Draft

**Input**: User description: "继续 Planned P2：Multi-instance complexBehaviorDefinition and none/one behavior event refs。Deploy 接受 behavior None/Complex 与 noneBehaviorEventRef / oneBehaviorEventRef；在实例完成里程碑抛出 signal/message；complex 条件用现有 loop 计数器求值。"

## User Scenarios & Testing *(mandatory)*

Multi-instance today supports All/One completion but rejects `complexBehaviorDefinition` and none/one behavior event refs. Authors need BPMN intermediate MI events (progress signals) and complex conditional throws without inventing non-OMG constructs.

### User Story 1 - oneBehaviorEventRef on first completion (Priority: P1)

An author models parallel MI User Tasks with `behavior="One"` and `oneBehaviorEventRef` pointing at a signal event definition. When the first inner instance Completes, the engine publishes that signal (and existing One completion cancels remaining inners).

**Why this priority**: Smallest additive behavior on an already-supported completion mode.

**Independent Test**: MI cardinality 3; Complete one inner; observe signal publication and host completion.

**Acceptance Scenarios**:

1. **Given** MI with `behavior="One"` and `oneBehaviorEventRef` to a signal, **When** the first inner Completes, **Then** the named signal is published and the MI host completes (remaining inners cancelled).

---

### User Story 2 - noneBehaviorEventRef after each inner (Priority: P1)

With `behavior="None"` and `noneBehaviorEventRef`, each inner Complete publishes the referenced signal/message; the host still completes when all inners finish (or `completionCondition`).

**Acceptance Scenarios**:

1. **Given** MI cardinality 2 with `behavior="None"` and a signal ref, **When** each inner Completes, **Then** the signal is published once per completion and the process completes after both.

---

### User Story 3 - complexBehaviorDefinition (Priority: P1)

With `behavior="Complex"` and one or more `complexBehaviorDefinition` entries (condition + implicit throw signal/message), after each inner Complete the engine evaluates conditions with loop counters; when a condition becomes true the associated signal/message is published once per definition.

**Acceptance Scenarios**:

1. **Given** Complex MI with condition `nrOfCompletedInstances >= 1` throwing signal `halfway`, **When** the first inner Completes, **Then** `halfway` is published once and further completes do not re-publish that definition’s event.
2. **Given** Complex MI without a matching throw type (e.g. only timer on the implicit throw), **When** Deploy runs, **Then** Deploy fails with `UNSUPPORTED_ELEMENT`.

---

### User Story 4 - Recover (Priority: P2)

After Recover mid-MI, already-satisfied complex / one milestones do not re-publish; new inner Completes still follow the rules.

**Acceptance Scenarios**:

1. **Given** one inner completed (and oneBehavior or complex event already due), **When** Recover then further Completes, **Then** the process finishes without duplicate milestone publishes for already-reached milestones.

---

### Edge Cases

- `noneBehaviorEventRef` / `oneBehaviorEventRef` MUST resolve to signal or message event definitions (or root signal/message usable as throw target); other kinds rejected at Deploy.
- Missing event ref on None/One is allowed (completion-only, no publish).
- `behavior="Complex"` without any `complexBehaviorDefinition` is rejected.
- Combining with `standardLoopCharacteristics` remains rejected.
- Standalone `implicitThrowEvent` FlowElements remain a separate Planned reject item; only nested MI complex events are in scope here.

## Requirements *(mandatory)*

- **FR-001**: Deploy MUST accept `behavior` None and Complex in addition to existing One/All.
- **FR-002**: Deploy MUST accept `noneBehaviorEventRef` and `oneBehaviorEventRef` when they resolve to supported signal/message throws.
- **FR-003**: Deploy MUST accept `complexBehaviorDefinition` entries with a condition and a signal/message implicit throw event.
- **FR-004**: After each inner Complete with `behavior=None` and a none ref, the engine MUST publish that signal/message.
- **FR-005**: When the first inner Completes with a one ref, the engine MUST publish that signal/message (One completion semantics unchanged).
- **FR-006**: After each inner Complete with Complex definitions, the engine MUST evaluate each condition once-fired; publish when newly satisfied.
- **FR-007**: Recover MUST NOT duplicate publishes for milestones already reached before restart.
- **FR-008**: Existing MI suites MUST remain green.

## Success Criteria *(mandatory)*

- **SC-001**: oneBehaviorEventRef fixture publishes once on first Complete and completes the host.
- **SC-002**: noneBehaviorEventRef fixture publishes once per inner Complete.
- **SC-003**: complexBehaviorDefinition fixture publishes when its condition first holds and not again.
- **SC-004**: Processing regression suite remains green.

## Assumptions

- Throw kinds limited to signal and message (existing Publication path).
- Loop counter names match existing MI completion evaluation (`nrOfCompletedInstances`, etc.).
- One-fire tracking for complex definitions is projection state rebuilt so already-true conditions after Recover are treated as already fired.
