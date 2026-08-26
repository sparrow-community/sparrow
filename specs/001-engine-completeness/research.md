# Research: Engine Completeness (next increment)

## 1. Call Activity: independent instance vs keep inline tokens

**Decision**: Replace in-definition inline Call Activity with a real child `process_instance_id`. The caller host token stays parked on the Call Activity (`ScopeHost`). Child tokens live only on the child instance.

**Rationale**: Spec FR-001/FR-004 require independent query, variables, and audit. Earlier inline Call Activity was documented as design debt in `processing/README.md`. Embedded SubProcess stays inline (same instance) — that construct is a different BPMN element.

**Alternatives considered**:

- Keep inline tokens and mint a synthetic child id in the projection only — fails audit (child events would still carry the parent's `process_instance_id`).
- Nested EventLog per call — extra ledger type; violates constitution II/III.

## 2. Resuming the caller without holding two locks

**Decision**: Child start is a COMMAND on the **parent** that appends parent CALL_ACTIVITY ACTIVATED (payload: `called_process_instance_id`) then, after the parent lock is released, enqueues CreateInstance-equivalent work on the child id (or runs it under the child's lock). Child PROCESS COMPLETED / TERMINATED publishes a parent-side Complete/Terminate of the Call Activity, again after the child lock is released — same `Effect.Publication` idea as message/signal throw.

**Rationale**: Constitution VI forbids processing two instance ids under one lock. Publication already exists to avoid re-entrant PublishMessage during a COMMAND.

**Alternatives considered**:

- Hierarchical lock (parent covers children) — deadlocks with sibling calls and breaks the partition key story.
- Complete the Call Activity inside the child's COMMAND by writing parent EVENTs — mixes partition keys in one append batch; Recover would be ambiguous.

## 3. IO mapping source

**Decision**: Compile `element.CallActivity` `DataInputAssociations` / `DataOutputAssociations` at deploy (already parsed on `element.Activity`). v1 mapping is name-to-name (source item / target item). Missing source → skip copy. Unknown structure → `UNSUPPORTED_ELEMENT` at deploy.

**Rationale**: Spec FR-007–FR-009; bpmn model already has the XML types. No new core element.

**Alternatives considered**:

- Copy all variables when mapping is absent — spec says empty + no write-back.
- Camunda `zeebe:ioMapping` extensions first — non-standard; constitution V prefers OMG associations.

## 4. Error Event Sub-Process

**Decision**: Extend existing Event Sub-Process arming (`EventSubProcessArm`) with `error_code` (empty = catch-all). On ERROR_THROWN / uncaught bubble, resolve like error boundaries (innermost first: activity boundary, then scope Event Sub-Process, then parent). Reuse interrupting vs `isInterrupting=false` behavior already used for message/timer/signal Event Sub-Process.

**Rationale**: Spec FR-011–FR-013; M4 error boundary/end/throw already exist. This is the missing start kind.

**Alternatives considered**:

- Treat error Event Sub-Process as a rewritten error boundary on the process — loses non-interrupting Event Sub-Process semantics.

## 5. Version coexistence

**Decision**: Keep `deployment_id` as the immutable snapshot key (already how CreateInstance works). On Deploy, assign `process_version = max(versions for that process id) + 1` (or 1 if new). Index process id → ordered revisions. CreateInstance stays deployment_id-based and gains optional `process_id` (latest) plus optional `process_version`. Running instances keep `deployment_id` + `process_version` unchanged.

**Rationale**: Spec FR-014–FR-016. Live migration is explicitly out of scope. Avoids rewriting the store around a new primary key.

**Alternatives considered**:

- Overwrite the previous deployment for a process id — breaks SC-005.
- Instance migration API in this increment — earlier roadmap bundled "versioning and migration"; spec deferred migration.

## 6. Protocol compatibility

**Decision**: Add optional fields only (`ProcessPayload` parent/call ids, `ActivityPayload.called_process_instance_id`, `engine.v1.Instance` parent/called links, `CreateInstanceRequest.process_id` / `process_version`). Run `protocol/proto/build.sh`. No field reuse or renumber.

**Rationale**: constitution `buf` FILE compatibility.

**Alternatives considered**: a new `call.v1` RPC service — unnecessary; EngineService already starts and queries instances.
