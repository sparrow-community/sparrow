# Feature Specification: Script Task

**Feature Branch**: `018-script-task`

**Created**: 2026-09-15

**Status**: Draft

**Input**: User description: "继续推进 Planned P1：Script Task。引擎不内嵌脚本运行时；与 Business Rule / Service Task 一样走 Job Activate/Complete/Fail；scriptFormat/name/id 作为 job type；script 正文可供 worker 使用。"

## User Scenarios & Testing *(mandatory)*

Authors cannot deploy `scriptTask` today. This increment accepts Script Task as a job-backed wait (no in-engine script execution), matching Business Rule Task and the constitution’s external-behavior mapping.

### User Story 1 - Script Task waits as a job (Priority: P1)

An author models start → Script Task (`scriptFormat="javascript"`, optional script body) → end. The instance waits for a worker. Activate for the job type returns the job (including script format/body when present). Complete finishes the task and the process.

**Why this priority**: Core Script Task executable semantics without embedding a language runtime.

**Independent Test**: Deploy fixture; CreateInstance → Activate `javascript` → Complete → completed.

**Acceptance Scenarios**:

1. **Given** a Script Task with `scriptFormat="javascript"`, **When** the instance is created, **Then** Activate for `javascript` returns that job.
2. **Given** that job, **When** Complete runs, **Then** the process completes and events show `TYPE_SCRIPT_TASK`.
3. **Given** a Script Task with no scriptFormat but a name, **When** Activate uses that name as job type, **Then** the job is claimed.
4. **Given** neither format nor name, **When** Activate uses the element id, **Then** the job is claimed.

---

### User Story 2 - Fail / incident / Recover (Priority: P2)

Script Task uses the same Fail / incident threshold / Recover path as Service Task and Business Rule Task.

**Acceptance Scenarios**:

1. **Given** a waiting Script Task job, **When** Fail then Activate again (or incident per threshold), **Then** behavior matches Service Task job fail rules.
2. **Given** a wait on Script Task, **When** Recover then Activate and Complete, **Then** outcome matches continuous.

---

### User Story 3 - Boundaries and multi-instance (Priority: P2)

Script Task may host interrupting boundaries and multi-instance like Service Task / Business Rule.

**Acceptance Scenarios**:

1. **Given** Script Task with a message boundary, **When** PublishMessage fires while waiting, **Then** interrupting boundary path works.
2. **Given** multiInstanceLoopCharacteristics on Script Task, **When** collection is provided, **Then** loop host/inner job waits follow existing MI job rules.

---

### Edge Cases

- Engine does not evaluate or execute the script body.
- Abstract `bpmn:task` remains rejected (separate Planned item).
- Script body may be empty; Deploy still accepts.
- Collaboration and in-engine language runtimes remain out of scope.

## Requirements *(mandatory)*

- **FR-001**: Deploy MUST accept `scriptTask` and index `TYPE_SCRIPT_TASK`.
- **FR-002**: Script Task MUST wait as a job; job type from scriptFormat, else name, else id.
- **FR-003**: Activate MUST return Script / ScriptFormat on the Job when present on the definition.
- **FR-004**: Complete / Fail / incident / Recover MUST match Service Task / Business Rule Task job semantics.
- **FR-005**: Boundaries and multi-instance MUST be supported like Service Task.
- **FR-006**: Existing suites MUST remain passing.

## Success Criteria *(mandatory)*

- **SC-001**: Script Task fixture Activate → Complete → process completed with SCRIPT_TASK events.
- **SC-002**: Recover mid-wait then Activate/Complete matches continuous.
- **SC-003**: Processing regression suite remains green.

## Assumptions

- No in-engine JavaScript/Groovy/etc.
- Job type priority: scriptFormat → name → id.
- Proto Element.Type already has TYPE_SCRIPT_TASK.
