# Feature Specification: Multi-Instance Activities

**Feature Branch**: `002-multi-instance`

**Created**: 2026-08-28

**Status**: Draft

**Input**: User description: "下一增量：实现 BPMN Multi-Instance（并行与顺序循环、基数与集合输入、完成条件），补齐可执行语义缺口；暂缓复杂行为事件、Multi-Instance Call Activity、集群与产品套件。"

## User Scenarios & Testing *(mandatory)*

This increment does not restart the engine. Authors can already deploy, start instances, complete waiting work, route through gateways, call child processes, and handle errors. The stories below add **multi-instance loop characteristics** on activities and embedded sub-processes so one modeled element can run multiple times before the flow continues.

### User Story 1 - Parallel multi-instance with fixed cardinality (Priority: P1)

An author marks a User Task (or Service Task) with parallel multi-instance loop characteristics and a fixed loop cardinality (for example `3`). When the flow reaches that activity, three independent inner instances of the activity start at once. Each inner instance can be completed separately. When all inner instances have completed (default completion), the activity completes once and the single outgoing token continues.

**Why this priority**: Parallel multi-instance with a numeric cardinality is the smallest slice that proves spawn, track, join, and continue — the core loop model.

**Independent Test**: Deploy a process with a parallel multi-instance user task (`loopCardinality` = 3) followed by an end event. Start the instance. Three waiting user tasks appear for the same activity id. Complete all three; the process reaches the end.

**Acceptance Scenarios**:

1. **Given** a parallel multi-instance activity with `loopCardinality` 3, **When** the flow arrives, **Then** three inner activity instances are active and the process does not continue past the activity until they are all completed.
2. **Given** two of three inner instances are completed, **When** the third is completed, **Then** the activity completes exactly once and the outgoing sequence flow is taken once.
3. **Given** a parallel multi-instance Service Task, **When** the flow arrives, **Then** three jobs are created (one per inner instance) and the activity completes only after all jobs succeed.

---

### User Story 2 - Sequential multi-instance (Priority: P1)

An author marks an activity with sequential multi-instance loop characteristics and a fixed cardinality. Inner instances run one after another: only one is active at a time; completing it starts the next until all are done.

**Why this priority**: Sequential loops share the same completion machinery as parallel but constrain concurrency — a common BPMN pattern and a distinct execution mode.

**Independent Test**: Deploy a sequential multi-instance user task with cardinality 3. Only one waiting user task exists at a time. After three completions in order, the process ends.

**Acceptance Scenarios**:

1. **Given** a sequential multi-instance activity with cardinality 3, **When** the flow arrives, **Then** exactly one inner instance is active.
2. **Given** the first inner instance is completed, **When** the engine advances, **Then** the second inner instance becomes active (the first does not re-open).
3. **Given** all three inner instances have completed in order, **When** the last completes, **Then** the activity completes once and the flow continues.

---

### User Story 3 - Collection-driven input (Priority: P2)

An author binds multi-instance input to a collection variable (`loopDataInputRef`) instead of a fixed cardinality. The number of inner instances equals the collection size. Each inner instance receives the current element (per `inputDataItem` name) as a variable for use inside the activity body.

**Why this priority**: Real processes iterate over lists (order lines, approvals). Cardinality-only loops are insufficient for data-driven iteration.

**Independent Test**: Start an instance with variable `items` = `["a","b"]`. A collection-driven parallel multi-instance user task creates two inner instances; each exposes the element variable (`item` = `"a"` or `"b"`). Complete both; optional output collection is written back if configured.

**Acceptance Scenarios**:

1. **Given** a collection variable with N elements and a parallel multi-instance activity referencing it, **When** the flow arrives, **Then** N inner instances start.
2. **Given** an inner instance is active, **When** an operator inspects variables for that loop context, **Then** the element variable holds the corresponding collection entry.
3. **Given** an empty collection, **When** the flow arrives, **Then** zero inner instances run and the activity completes immediately (no waiting work).

---

### User Story 4 - Completion condition before all instances finish (Priority: P2)

An author sets a completion condition expression (for example "at least two instances complete" or the standard `One` behavior). Inner instances still start per parallel/sequential rules, but the activity may complete early when the condition becomes true; remaining active inner instances are cancelled.

**Why this priority**: Many workflows stop after the first success or a quorum without waiting for stragglers.

**Independent Test**: Parallel multi-instance with cardinality 5 and completion condition `nrOfCompletedInstances >= 2`. Complete any two inner instances; the activity completes and any still-waiting inner instances are terminated.

**Acceptance Scenarios**:

1. **Given** a completion condition requiring two completed instances out of five parallel inner instances, **When** any two complete, **Then** the activity completes and the flow continues.
2. **Given** the completion condition is satisfied while other inner instances are still waiting, **When** the activity completes, **Then** those waiting inner instances are cancelled and do not block the process.
3. **Given** no completion condition is declared, **When** the activity runs, **Then** the default is that all started inner instances must complete (equivalent to `All` behavior).

---

### User Story 5 - Multi-instance embedded sub-process (Priority: P3)

An author applies multi-instance loop characteristics to an embedded Sub-Process. Each inner instance runs the sub-process body independently (parallel or sequential per `isSequential`), then the sub-process join semantics apply per inner instance before the outer activity completes.

**Why this priority**: Sub-process MI is common for "for each item, run this mini-flow" and reuses the same loop machinery on a scope host.

**Independent Test**: Parallel multi-instance sub-process (cardinality 2), each containing start → user task → end. Two sub-process scopes run; completing both user tasks finishes the multi-instance sub-process once.

**Acceptance Scenarios**:

1. **Given** a parallel multi-instance embedded sub-process with cardinality 2, **When** the flow arrives, **Then** two sub-process scopes are active.
2. **Given** one sub-process scope has completed and the other is still waiting, **When** default `All` completion applies, **Then** the outer multi-instance sub-process does not complete yet.
3. **Given** both scopes complete, **When** the join is satisfied, **Then** the multi-instance sub-process completes once.

---

### Edge Cases

- What happens when `loopCardinality` evaluates to zero or negative? Treat as zero instances: the activity completes immediately with no waiting work.
- What happens when the collection variable is missing or not a JSON array? Treat as empty collection (zero instances, immediate completion) — do not crash the instance.
- What happens when Recover runs mid-loop? Inner instance counters, active/waiting inner tokens, and partial output collection MUST rebuild so completing or cancelling remaining work still yields the same outcome as before the crash.
- What happens when an interrupting boundary fires on a multi-instance activity? All active inner instances and the loop host MUST be terminated; the boundary path is taken once.
- What happens when an error is thrown from one inner instance? Existing error-boundary and bubble rules apply within that inner scope; if uncaught, behavior matches a non-multi-instance activity (may terminate the instance or trigger Event Sub-Process).
- How does the system handle multi-instance on an unsupported element at deploy time? Deploy MUST reject with a clear unsupported error before any instance starts.
- What happens with output collection (`loopDataOutputRef`)? Each inner instance's output item variable is appended in loop order; if no output mapping is configured, no collection variable is written.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Activities and embedded sub-processes with `multiInstanceLoopCharacteristics` MUST be supported at deploy time when the element type is otherwise executable (User Task, Service Task, embedded Sub-Process).
- **FR-002**: Parallel multi-instance (`isSequential=false`) MUST allow multiple inner instances of the same activity to be active concurrently within one process instance.
- **FR-003**: Sequential multi-instance (`isSequential=true`) MUST allow at most one inner instance active at a time.
- **FR-004**: Loop cardinality from `loopCardinality` MUST determine the number of inner instances when no collection input is configured.
- **FR-005**: Collection input from `loopDataInputRef` MUST determine the number of inner instances and MUST expose each element to inner instances per `inputDataItem`.
- **FR-006**: Default completion MUST require all started inner instances to complete successfully (`All` behavior) when no `completionCondition` is present.
- **FR-007**: A declared `completionCondition` MUST be evaluated after each inner instance completes (or is cancelled) and MAY complete the activity early when true.
- **FR-008**: When the activity completes early, all still-active inner instances MUST be cancelled.
- **FR-009**: After the multi-instance activity completes, the outgoing sequence flow MUST be taken exactly once.
- **FR-010**: Optional `loopDataOutputRef` / `outputDataItem` MUST assemble inner outputs into a collection variable on the process instance.
- **FR-011**: Standard loop metadata (`loopCounter`, `numberOfInstances`, `numberOfActiveInstances`, `numberOfCompletedInstances`) MUST be available to completion conditions and MAY be exposed as variables during the loop.
- **FR-012**: Recover MUST restore multi-instance state (counts, waiting inner work, partial collections) without losing parent/child token relationships.
- **FR-013**: Interrupting boundaries on a multi-instance activity MUST cancel all inner instances and the loop host consistently.
- **FR-014**: Deploy MUST reject multi-instance on unsupported element types and reject `complexBehaviorDefinition` / none-one behavior event refs in this increment.
- **FR-015**: Existing non-multi-instance behavior for the same element types MUST remain unchanged.

### Key Entities

- **Multi-instance activity**: A flow element annotated with loop characteristics; acts as the loop host.
- **Inner instance**: One iteration of the loop (index, element input, own waiting work or sub-process scope).
- **Loop counters**: Active, completed, and total instance counts used for join and completion evaluation.
- **Collection binding**: Input collection variable and per-iteration element variable; optional output collection assembly.
- **Completion condition**: Expression evaluated against loop counters and instance variables to allow early join.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A parallel multi-instance user task with cardinality 3 exposes three completable waiting work items before the process can continue.
- **SC-002**: A sequential multi-instance task never has more than one active inner instance at a time.
- **SC-003**: A collection with N elements produces exactly N inner instances; an empty collection completes the activity with zero waiting work.
- **SC-004**: With a completion condition requiring 2 of 5 completions, the process continues after exactly two inner completes and does not leave stray waiting work.
- **SC-005**: Recover after a crash with 1 of 3 parallel inner instances completed restores one completed and two waiting; finishing the two waiting instances completes the activity.
- **SC-006**: Non-multi-instance regression suite (`go test ./processing/`) continues to pass.

## Assumptions

- M1–M4c and `001-engine-completeness` capabilities remain the baseline.
- `bpmn` already parses `multiInstanceLoopCharacteristics`; this increment adds **execution** only.
- Cardinality expressions and completion conditions use the same lightweight expression evaluator already used for sequence-flow conditions (JSON variables; missing names treated as absent/false).
- Only User Task, Service Task, and embedded Sub-Process support multi-instance in this increment.
- Multi-instance on Call Activity, complex behavior definitions, none/one behavior events, and collection expressions beyond variable references stay deferred.
- Cluster, product-suite UI, Incident records, live migration, cross-deployment Call Activity, and other deferred items from `AGENTS.md` stay out of scope.
- Inner instance indices are zero-based for `loopCounter` unless BPMN fixture tests in the repo dictate otherwise.
- Output collection ordering follows ascending loop index.
