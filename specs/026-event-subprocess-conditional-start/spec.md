# Feature Specification: Event Sub-Process Conditional Start

**Feature Branch**: `main` (developed on main)

**Created**: 2026-09-15

**Status**: Draft

**Input**: User description: "继续推进 Planned：Event Sub-Process conditional start。保持在主分支。说明与文档中不要使用简写。"

## User Scenarios & Testing *(mandatory)*

Event sub-process starts already support message, signal, timer, error, escalation, and compensate. Conditional starts are rejected at Deploy. Intermediate conditional catches and process-level conditional starts already use `EvaluateConditions` / `EvaluateConditionalStarts`. This increment accepts a conditional start on an event sub-process and triggers it via `EvaluateConditions` when the condition is true.

### User Story 1 - Interrupting conditional event sub-process (Priority: P1)

An author nests a `triggeredByEvent` SubProcess with a conditional start. While the host process waits on a user task, `EvaluateConditions` with matching variables starts the event sub-process, cancels the host wait (interrupting), and the handler path completes the process.

**Acceptance Scenarios**:

1. **Given** an armed conditional event sub-process and a waiting user task, **When** EvaluateConditions makes the condition true, **Then** the event sub-process runs and the user task is terminated.
2. **Given** the same model, **When** EvaluateConditions runs with a false condition, **Then** the event sub-process stays armed and the user task remains waiting.

### User Story 2 - Non-interrupting conditional event sub-process (Priority: P2)

With `isInterrupting="false"`, EvaluateConditions starts the event sub-process without cancelling the host user task; both can complete.

### User Story 3 - Recover while armed (Priority: P2)

Recover while the conditional event sub-process is armed and the host is waiting; after Recover, EvaluateConditions still triggers the same outcome.

### Edge Cases

- Condition expression language is the same as intermediate conditional catch (`expr`).
- Empty condition expression is rejected at Deploy.
- Compensation event sub-process remains separate (not armed as a live conditional wait).
- Documentation spells out “event sub-process” without abbreviation.

## Requirements *(mandatory)*

- **FR-001**: Deploy MUST accept an event sub-process whose start has exactly one `conditionalEventDefinition` with a non-empty condition.
- **FR-002**: Scope entry MUST arm the conditional event sub-process start (START_EVENT ACTIVATED) like other live event sub-process kinds.
- **FR-003**: EvaluateConditions MUST trigger armed conditional event sub-processes whose condition evaluates true (using instance variables plus request overlay).
- **FR-004**: Interrupting vs non-interrupting semantics MUST match existing event sub-process trigger behavior.
- **FR-005**: Existing event sub-process and conditional suites MUST remain green.
- **FR-006**: Recover MUST restore the armed conditional start so EvaluateConditions still works.

## Success Criteria *(mandatory)*

- **SC-001**: Interrupting fixture: EvaluateConditions fires event sub-process; host wait terminated; process completes.
- **SC-002**: False condition leaves host waiting and event sub-process armed.
- **SC-003**: Non-interrupting fixture keeps host wait after trigger.
- **SC-004**: Recover then EvaluateConditions matches continuous run.
- **SC-005**: Processing regression suite remains green.
