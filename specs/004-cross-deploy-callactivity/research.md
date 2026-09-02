# Research: Cross-deployment Call Activity

## 1. Deploy-time validation: embedded vs external callee

**Decision**: **Dual-mode** compile. If `calledElement` matches a process in the same definitions catalog, keep today's behavior: embed callee into `calledProcesses`, compile `StartEventID`, validate callee scope. If `calledElement` is non-empty but **not** in the catalog, compile a `CallActivity` ref with `ExternalCallee: true`, IO mappings only, and **no** deploy failure.

**Rationale**: FR-001 and spec Story 1 scenario 3 — caller-only deploy must succeed. Bundled same-file tests (SC-005) keep the embedded path unchanged.

**Alternatives considered**:

- Require callee BPMN in deploy store at caller deploy time — rejected; blocks independent deploy order.
- Always external (never embed) — rejected; regresses same-file optimization and M1 validation of bundled callee.
- Separate BPMN extension `sparrow:externalCall` — rejected for MVP; absence from catalog is sufficient signal.

## 2. Runtime callee resolution

**Decision**: Resolve `calledElement` process id via existing `Engine.resolveDeploymentLocked("", processID, 0)` → `deploy.ResolveRevision` for **latest** revision. Reuse the same rule as `CreateInstance` when only `process_id` is supplied.

**Rationale**: Spec assumption and edge-case table — latest revision; single-node engine already maintains `e.revisions` index.

**Alternatives considered**:

- Pin caller deployment's compile-time revision — impossible when callee not embedded.
- Global "deployment alias" table — extra store; out of scope.
- Resolve at deploy and store callee deployment id on CallActivity — callee may not exist yet at caller deploy.

## 3. FR-006: rejection without waiting projection

**Decision**: For `ExternalCallee`, **pre-resolve** callee deployment inside `CallActivityHandler.OnEnter` **before** emitting ACTIVATING/ACTIVATED. On `NOT_FOUND`, return error from `OnEnter` with no `Effect.Records` — executor emits nothing; upstream command path surfaces REJECTION. Do **not** defer resolution solely to `startCalledInstance` after ACTIVATED (current flush path would leave caller waiting with `called_process_instance_id`).

**Rationale**: FR-006 — caller must not appear to have started a successful call. `OnEnter` error before emit matches constitution VI.

**Alternatives considered**:

- Emit ACTIVATED then compensate with TERMINATED on child start failure — extra ledger noise and race with boundaries.
- REJECTION only on flushPublications error — parent projection already mutated; rejected.

**Embedded path**: No runtime resolution; compile-time `StartEventID` and caller deployment id unchanged.

## 4. Child instance deployment id and Enter context

**Decision**: `startCalledInstance` creates child `projection.Instance` and PROCESS COMMAND/EVENT with **`calleeDep.ID`** and **`calleeDep.Version`**. `executor.Enter` and `emitEventSubProcessStartArms` use **`calleeDep`**, not parent deployment. IO mappings still read from **caller** `CallActivitySpec` on parent deployment.

**Rationale**: FR-002, FR-003 — child ledger and GetInstance show callee deployment; parent linkage fields unchanged on child PROCESS payload.

**Alternatives considered**:

- Child shares parent deployment id with different process id — breaks GetInstance, Recover deployment load, and callee element index.
- Store callee BPMN snapshot on caller deployment — parallel graph; rejected.

## 5. Publication and redrive consistency

**Decision**: Extend `handlers.Publication` with `CalledDeploymentID` set at OnEnter (from resolver). `startCalledInstance` uses this id directly rather than re-resolving on redrive, so Recover redrive does not bind to a **newer** callee revision if one was deployed between crash and recover.

**Rationale**: FR-007 idempotency — same command chain → same child deployment id. Matches child CREATE_INSTANCE COMMAND already appended with fixed `deployment_id`.

**Alternatives considered**:

- Re-resolve latest on every redrive — could start different callee revision mid-recover; rejected.

## 6. Recover with two deployments

**Decision**: No Recover changes required beyond correct child `deployment_id` in ledger. `loadDeployments()` loads all files from `deploy.Store`; replay projects parent (caller dep) and child (callee dep) independently. Redrive `startCalledInstance` from unfinished parent Enter publication uses stored `CalledDeploymentID` and preallocated `ChildInstanceID`.

**Rationale**: US3 — both deployments must be present in store (same engine node assumption). Existing recover tests pattern applies.

**Alternatives considered**:

- Lazy-load callee deployment on replay — unnecessary if store has both artifacts (spec assumption).

## 7. Wire / API surface

**Decision**: **No new RPC.** `GetInstance` already returns `deployment_id`, `parent_process_instance_id`, and token `called_process_instance_id` (FR-008). Optional additive `ActivityPayload.called_process_id` on CALL_ACTIVITY ACTIVATED for log audit — implement only if tasks want explicit trace without joining instances.

**Rationale**: Minimize proto churn; instance query suffices for operators.

**Alternatives considered**:

- `ResolveCallActivity` RPC — product surface; rejected.
- Required proto change for callee deployment on token — redundant with child instance record.

## 8. Negative and edge behavior

**Decision**:

- Callee not deployed at activation → `NOT_FOUND` rejection on parent token enter; no child COMMAND.
- Callee redeployed while child active → existing child keeps original deployment id from its CREATE_INSTANCE COMMAND.
- New calls after callee redeploy → latest revision at next OnEnter resolve.
- Boundary interrupt on Call Activity → unchanged `terminateCalledInstance` using `inst.DeploymentID` (callee).
- Nested cross-deploy calls → each level resolves independently per `calledElement`.

**Rationale**: Spec edge cases; parity with same-deployment boundary behavior.

**Alternatives considered**:

- Fail active children when callee deployment file removed — spec says existing children continue.
