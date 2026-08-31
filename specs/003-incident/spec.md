# Feature Specification: Incident (blocked execution state)

**Feature Branch**: `003-incident`

**Created**: 2026-08-31

**Status**: Draft

**Input**: User description: "Engine increment 003 — Incident: durable failed/blocked state with audit trail; retry and resume via explicit COMMAND; Recover testable. First P1 item after multi-instance per AGENTS.md Remaining list."

## User Scenarios & Testing *(mandatory)*

Callers can already deploy, start instances, complete waiting work, activate/fail Service Task jobs, and recover from mid-command crashes. Today a failed job keeps the activity waiting and the worker may retry indefinitely; there is no durable **incident** state that marks an instance as blocked pending operator action, and no first-class **resolve/retry** command that is auditable separately from another anonymous Fail/Complete attempt.

### User Story 1 - Job failure opens an incident (Priority: P1)

An operator runs a process with a Service Task job. The worker reports failure enough times (or with an explicit “no retry” outcome) that the engine records an **incident** on that activity. The process instance remains identifiable as blocked: the waiting work is not silently retriable without an operator resolve step. The event log contains an auditable incident lifecycle (opened at minimum).

**Why this priority**: Without a durable incident, production and agent-driven flows cannot distinguish “try again later” from “human must intervene,” and blocked instances are invisible except by reading many FAILED events.

**Independent Test**: Deploy a one Service Task process. Start an instance, activate the job, fail it until incident threshold. `GetInstance` shows blocked/incident state on the activity token; completing the process without resolve is rejected or impossible.

**Acceptance Scenarios**:

1. **Given** a waiting Service Task job, **When** the worker fails with retry still allowed, **Then** the token stays waiting, no incident is opened, and another worker may activate the job.
2. **Given** a waiting Service Task job, **When** failures reach the configured incident threshold (or a fail marked as non-retriable), **Then** an incident is opened on that element/token and the instance is blocked at that activity.
3. **Given** an open incident, **When** a caller attempts an unrelated Complete on another token, **Then** behavior follows existing rules; the incident does not corrupt other tokens.
4. **Given** an open incident, **When** a caller attempts Complete on the incident activity without resolve, **Then** the command is rejected with a stable code and a REJECTION is appended; projection is unchanged.

---

### User Story 2 - Resolve incident and retry work (Priority: P1)

An operator (or automation via explicit COMMAND) **resolves** an open incident so the blocked activity returns to a normal waiting job state. The worker may activate and complete the job afterward; the process continues as before the incident.

**Why this priority**: Incident without resolve is a dead end; resolve/retry is the minimum viable operator loop.

**Independent Test**: Open an incident on a Service Task (Story 1). Issue resolve/retry. Activate job again, complete successfully, process reaches end.

**Acceptance Scenarios**:

1. **Given** an open incident on a Service Task, **When** resolve/retry is commanded, **Then** the incident closes, the token is waiting with job type armed, and FAILED/resolve events appear in the audit trail.
2. **Given** a resolved incident, **When** a worker activates and completes the job, **Then** the process continues along the normal outgoing flow.
3. **Given** no open incident on the token, **When** resolve/retry is commanded, **Then** the command is rejected; no silent success.

---

### User Story 3 - Recover with an open incident (Priority: P2)

After a crash, Recover rebuilds instance state including any open incident. Resolve/retry and subsequent job completion still succeed without duplicating incident open/close facts.

**Why this priority**: Completeness definition requires Recover tests; incidents must be ledger-backed, not memory-only.

**Independent Test**: Open incident, rebuild engine from EventLog only, confirm incident still visible, resolve and finish process.

**Acceptance Scenarios**:

1. **Given** an open incident persisted in the log, **When** Recover runs, **Then** projection shows the same blocked/incident state and waiting token identity.
2. **Given** Recover after resolve was commanded but EVENTs were partially written, **Then** redrive finishes the chain without duplicate incident opens (same idempotency rules as other commands).

---

### Edge Cases

- What happens when Fail is called on a User Task or non-job waiting token? No incident; existing REJECTION/behavior unchanged.
- What happens when an interrupting boundary fires on an activity with an open incident? Boundary semantics win; incident is closed or superseded consistently with termination of the host activity (exact path tested in plan).
- What happens when the process instance completes or terminates while resolving? Resolve on a closed instance is rejected.
- What happens when multiple Fail commands arrive concurrently for the same job? Serial instance lock preserves a single incident open and consistent retry count.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The engine MUST record incident open (and close on resolve) as append-only EVENT facts tied to process instance, element id, and token id.
- **FR-002**: While an incident is open on an activity token, the instance MUST be queryable as blocked at that activity (distinct from normal active/waiting-only job retry).
- **FR-003**: Fail on a Service Task job MUST support either a retry budget before opening an incident, or an explicit non-retriable fail that opens an incident immediately (default threshold documented in plan).
- **FR-004**: Resolve/retry MUST be a explicit COMMAND that closes the incident and restores normal waiting job semantics on that token.
- **FR-005**: Invalid resolve, Complete on incident token without resolve, or Complete on wrong token MUST append REJECTION and MUST NOT mutate projection as if successful.
- **FR-006**: Recover MUST rebuild open incidents and retry counts from the log; redrive MUST not duplicate incident EVENTs for the same command.
- **FR-007**: `GetInstance` (and list-by-instance audit) MUST expose enough incident information for an operator to identify element, token, and error context without reading raw logs only.

### Key Entities

- **Incident**: A blocked execution point on an instance — links to element id, token id, optional error message, opened/closed lifecycle, optional retry count at open.
- **Incident threshold**: Policy for how many retriable fails occur before open (per deployment default or job extension — resolved in plan).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: In a fixture process, failing a job until threshold produces exactly one open incident visible in instance query and in the event log.
- **SC-002**: Resolve/retry followed by successful job completion completes the process in one pass without manual projection repair.
- **SC-003**: Recover after incident open restores blocked state; resolve and complete after recover matches non-recover outcome.
- **SC-004**: Invalid operator commands on incident tokens produce REJECTION records in 100% of tested negative cases (wrong token, double resolve, complete while blocked).

## Assumptions

- MVP scope is **Service Task job failures**; uncaught BPMN errors, error boundaries, and process termination keep existing semantics unless they naturally close an incident host.
- Incident does not replace Job Fail/Activate/Complete APIs; it adds lifecycle on top.
- Retry threshold default is small (e.g. 3 retriable fails) unless BPMN extension specifies otherwise — exact default chosen in plan.
- Operator-facing resolve is exposed through the same engine/gateway surface as other commands (details in plan/contracts).
- Cross-deployment CallActivity, live migration, and AI-specific tooling are separate increments.
