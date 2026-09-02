# Implementation Plan: Cross-deployment Call Activity

**Branch**: `004-cross-deploy-callactivity` | **Date**: 2026-09-02 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/004-cross-deploy-callactivity/spec.md`

## Summary

Allow Call Activity `calledElement` to reference a process deployed **outside** the caller's BPMN file. Deploy compiles caller-only definitions when the callee is absent from the catalog; at activation the engine resolves `calledElement` to the **latest** deployed revision, starts the child with the **callee deployment id**, and keeps IO mapping, parent/child linkage, boundaries, and Recover behavior aligned with same-deployment calls. Same-file bundled callee behavior is unchanged (SC-005).

## Technical Context

**Language/Version**: Go 1.26.5 (workspace `go.work`)

**Primary Dependencies**: in-tree `bpmn`, `protocol` (Protobuf + generated Go), `processing` (`deploy`, `handlers/call_activity.go`, `call_child.go`, `engine.go`, `recover.go`, `projection`), `gateway` (adapters only — no new RPC expected)

**Storage**: append-only `log.EventLog` + `deploy.Store` (multiple deployment files); Recover loads all deployments from store

**Testing**: `go test ./processing/ ./gateway/ ./protocol/proto/event/v1/`; new fixtures `m8_cross_call_*.bpmn`; existing `m4_call_activity.bpmn` / `m5_call_*.bpmn` must pass unchanged

**Target Platform**: single-node process (`gateway/cmd/sparrow`)

**Project Type**: Go multi-module BPMN execution engine

**Performance Goals**: correctness and recoverability; callee resolution is O(1) map lookup on revision index

**Constraints**: constitution (no parallel graph, Element subject unchanged, serial commands per instance); child COMMAND/EVENT must record callee `deployment_id`; FR-006 requires no waiting projection when callee missing

**Scale/Scope**: three user stories (cross-deploy spawn, IO mapping, Recover). MVP: latest revision by process id. Deferred: explicit version pin on Call Activity, multi-instance Call Activity, cross-node registry.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status |
|-----------|--------|
| I. Engine, not product suite | Pass — deploy/call semantics only; no console |
| II. Ledger is source of truth | Pass — child PROCESS COMMAND/EVENT carry callee `deployment_id`; Recover replays from log |
| III. Element is the subject | Pass — CALL_ACTIVITY and PROCESS intents unchanged; no new element type |
| IV. No parallel graph; thin modules | Pass — `deploy` compile split + `call_child` resolution; handlers get injected resolver hook |
| V. Standard BPMN | Pass — `calledElement` is standard; external reference is deploy/runtime policy |
| VI. Serial commands, explicit rejection | Pass — callee NOT_FOUND before ACTIVATED → REJECTION; no phantom waiting token |

Post-design re-check: still pass. Optional additive `ActivityPayload.called_process_id` is audit-only; projection uses existing `Instance.deployment_id` and token `called_process_instance_id`.

## Project Structure

### Documentation (this feature)

```text
specs/004-cross-deploy-callactivity/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── engine.md
└── tasks.md              # Phase 2 (/speckit-tasks)
```

### Source Code (repository root)

```text
processing/
├── deploy/
│   ├── call_activity.go   # dual-mode validate: embedded vs external ref
│   └── deploy.go          # indexCallActivities: skip embed when external
├── handlers/
│   ├── call_activity.go   # pre-resolve external callee; Publication.CalledDeploymentID
│   └── handler.go         # Publication struct extension
├── call_child.go          # child instance + Enter on callee deployment
├── engine.go              # wire callee resolver into handlers (existing resolveDeploymentLocked)
└── call_activity_test.go  # cross-deploy + recover cases; regression on m4/m5

processing/testdata/
└── m8_cross_call_*.bpmn   # caller-only + callee-only fixtures
```

**Structure Decision**: Extend existing `processing/deploy` and `call_child` paths. No new module or runtime store.

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| None | — | — |

## Phase 0 & 1 Outputs

- [research.md](./research.md) — dual-mode deploy, runtime resolution, FR-006 pre-resolve, Recover deployment ids
- [data-model.md](./data-model.md) — `CallActivity` compile fields, publication fields, rejection codes
- [contracts/engine.md](./contracts/engine.md) — optional proto audit field; deploy validation contract
- [quickstart.md](./quickstart.md) — per-story validation steps

## Implementation Notes (for tasks phase)

1. **Deploy**: Split `validateCallActivity` into ref extraction (always) + embedded catalog validation (when callee in file). Add `CallActivity.ExternalCallee bool` and empty `StartEventID` for external refs. `indexCallActivities` skips `calledProcesses` embed when external.
2. **Handler hook**: `handlers.SetCalleeResolver(func(processID string) (deploymentID string, err error))` called from `processing` init; external OnEnter resolves before emitting ACTIVATED; failure returns `NOT_FOUND` without records.
3. **Publication**: Add `CalledDeploymentID`; handler sets it for external (and embedded may set caller id). `startCalledInstance` uses callee deployment for child COMMAND, projection, and `Enter`.
4. **startCalledInstance**: `calleeDep := e.deployments[p.CalledDeploymentID]`; `StartEventID` from compile or `calleeDep.StartEventID()`; `emitEventSubProcessStartArms(calleeDep, ...)`.
5. **Recover**: Child CREATE_INSTANCE COMMAND already keyed by `deployment_id` in log; ensure redrive uses same id from publication/command (no re-resolve to newer revision mid-redrive). Existing `loadDeployments` + `recover.go` redrive sufficient when both deployments in store.
6. **Tests**: SC-001–SC-005 per quickstart; `recover_test.go` cross-deploy case (US3).
7. **Proto (optional)**: `ActivityPayload.called_process_id` on CALL_ACTIVITY ACTIVATED for audit; skip if GetInstance linkage is enough for MVP.
8. **AGENTS.md**: move cross-deployment Call Activity to Implemented after implement phase.
