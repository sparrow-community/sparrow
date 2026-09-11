# Research: Terminate End and Remaining Tasks

## Terminate scope vs instance status

**Decision**: Cancel remaining tokens in the enclosing scope, then complete that scope. Process-level terminate → instance **completed** (not `terminated`). Unhandled error remains `terminated`.

**Rationale**: BPMN terminate consumes remaining tokens and completes the process/sub-process. Sparrow already uses `StatusTerminated` for unhandled errors; mixing them would make Call Activity parents treat a successful terminate as a failed child.

**Alternatives considered**: Mark process `TERMINATED` — rejected (Call Activity resume `Completed=false`).

## Send vs job

**Decision**: Send Task publishes a message (name from `messageRef` / name / id) and continues; no job wait.

**Rationale**: Matches message throw and BPMN “send then continue”. Workers already cover Service/Business Rule.

**Alternatives considered**: Always job (Zeebe-style) — extra API for a completeness item that already has throw.

## Receive vs catch

**Decision**: Distinct `TYPE_RECEIVE_TASK`; wait with `MessageName` on the activity payload; PublishMessage completes it. Instantiate receive rejected.

**Rationale**: Constitution forbids collapsing types. Instantiate receive is a process-start variant, deferred like parallel instantiate EBG.

## Business Rule vs DMN

**Decision**: Job wait; implementation/name/id → job type; incident threshold extensions reuse Service Task indexing.

**Rationale**: Engine must not embed a rule engine; constitution maps custom behavior onto Service Task-like jobs.

## Manual vs User

**Decision**: Same wait/Complete path; `TYPE_MANUAL_TASK`.

**Rationale**: BPMN distinguishes them; ledger Type must match.
