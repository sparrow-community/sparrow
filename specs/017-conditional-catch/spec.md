# Feature Specification: Conditional Catch (Intermediate and Boundary)

**Feature Branch**: `017-conditional-catch`

**Created**: 2026-09-15

**Status**: Draft

**Input**: User description: "继续推进 Planned P1：Conditional catch — intermediate and boundary。Deploy 接受；等待后由 EvaluateConditions 在条件为真时完成；已为真时可立即通过；Recover 覆盖。"

## User Scenarios & Testing *(mandatory)*

Conditional sequence flows and process-level conditional starts exist. Intermediate and boundary conditional catches are still rejected at Deploy. This increment accepts them, waits when the condition is false, and completes when `EvaluateConditions` finds the condition true (same expression language as flows/starts).

### User Story 1 - Intermediate conditional catch (Priority: P1)

An author models start → conditional catch (`approved == true`) → end. If variables already make the condition true at enter, the catch completes immediately. Otherwise the instance waits until EvaluateConditions supplies or sees a true condition.

**Why this priority**: Core mid-process conditional wait.

**Independent Test**: False at enter → wait; EvaluateConditions true → completed. True at enter → completed without Evaluate.

**Acceptance Scenarios**:

1. **Given** start → conditional catch → end and CreateInstance with `{approved: false}`, **When** enter reaches the catch, **Then** the instance waits on the catch.
2. **Given** that wait, **When** EvaluateConditions runs with overlay/vars `{approved: true}`, **Then** the catch completes and the process completes.
3. **Given** CreateInstance with `{approved: true}`, **When** the catch is entered, **Then** the process completes without a separate EvaluateConditions.

---

### User Story 2 - Conditional boundary (Priority: P1)

An author attaches a conditional boundary to a User Task. While the task waits, EvaluateConditions with a true condition fires the boundary (interrupting cancels the task; non-interrupting leaves the task waiting and takes the boundary path).

**Why this priority**: Completes boundary parity with timer/message/signal.

**Independent Test**: User Task + interrupting conditional boundary; EvaluateConditions true → boundary path; task terminated.

**Acceptance Scenarios**:

1. **Given** User Task with interrupting conditional boundary (`approved == true`), **When** instance waits on the task, **Then** the boundary is armed.
2. **Given** that wait, **When** EvaluateConditions with `{approved: true}` runs, **Then** the boundary completes, the User Task is terminated, and the process follows the boundary path.
3. **Given** a non-interrupting conditional boundary, **When** EvaluateConditions makes it true, **Then** the boundary path runs and the User Task remains waiting until Complete.
4. **Given** two conditional boundaries on one activity with the same condition text, **When** Deploy runs, **Then** deploy is rejected (unique condition discriminator per activity).

---

### User Story 3 - Recover (Priority: P2)

Recover while waiting on a conditional catch or armed conditional boundary; EvaluateConditions then finishes matches continuous.

**Acceptance Scenarios**:

1. **Given** a wait on intermediate conditional catch, **When** Recover then EvaluateConditions true, **Then** process completes.
2. **Given** a wait on User Task with conditional boundary armed, **When** Recover then EvaluateConditions true, **Then** interrupting boundary path completes as continuous.

---

### Edge Cases

- Event Sub-Process conditional start remains Planned (out of scope).
- Expression language matches conditional flows / conditional starts (`processing/expr`).
- EvaluateConditions may target one instance or all; optional variables overlay the evaluation environment and are applied on successful complete.
- Mixed event definitions on a catch/boundary rejected.
- Existing timer/message/signal/link catches unchanged.

## Requirements *(mandatory)*

- **FR-001**: Deploy MUST accept intermediate catch with exactly one conditionalEventDefinition and a non-empty condition.
- **FR-002**: Deploy MUST accept boundary events with exactly one conditionalEventDefinition (interrupting and non-interrupting).
- **FR-003**: Intermediate conditional catch MUST wait when false at enter and complete when EvaluateConditions finds true (or complete immediately when true at enter).
- **FR-004**: Conditional boundaries MUST arm while the host waits and fire via EvaluateConditions when true, with interrupting/non-interrupting semantics matching other boundaries.
- **FR-005**: Duplicate identical condition text on two conditional boundaries of the same activity MUST be rejected at Deploy.
- **FR-006**: Engine MUST expose EvaluateConditions (deployment-agnostic instance scan; optional instance id and variables).
- **FR-007**: Recover MUST restore waits so EvaluateConditions matches continuous.
- **FR-008**: Existing suites MUST remain passing.

## Success Criteria *(mandatory)*

- **SC-001**: Intermediate catch false→wait→Evaluate true→completed; true at enter→completed.
- **SC-002**: Interrupting conditional boundary fires via EvaluateConditions.
- **SC-003**: Non-interrupting conditional boundary fires without canceling host.
- **SC-004**: Duplicate condition boundaries rejected; Recover scenarios green.
- **SC-005**: Processing regression suite remains green.

## Assumptions

- No proto change required; condition text lives in deploy index; boundary kind string `"conditional"`.
- EvaluateConditions does not create process-level conditional starts (EvaluateConditionalStarts remains separate).
- Variable overlay on EvaluateConditions is passed into Complete when a catch/boundary fires.
