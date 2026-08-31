# Implementation Plan: Incident (blocked execution state)

**Branch**: `003-incident` | **Date**: 2026-08-31 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/003-incident/spec.md`

## Summary

Add a durable **incident** lifecycle for Service Task job failures: retriable fails accumulate until a default threshold (3) or an explicit `no_retry` fail opens an incident; blocked tokens are queryable via `GetInstance`; operators **resolve** through a new auditable COMMAND/RPC; Recover rebuilds blocked state from the event log. Implementation extends existing SERVICE_TASK intents and token projection — no new ledger subject or parallel graph.

## Technical Context

**Language/Version**: Go 1.26.5 (workspace `go.work`)

**Primary Dependencies**: in-tree `bpmn` (optional `sparrow:failedJobIncidentThreshold` extension), `protocol` (Protobuf + generated Go), `processing` (`jobs.go`, `engine.go`, `projection`, `deploy`, `recover`), `gateway` (Engine + Job adapters)

**Storage**: append-only `log.EventLog` + `deploy.Store` + optional `runtime.Store` (job leases unchanged; blocked tokens excluded from activate queue)

**Testing**: `go test ./processing/ ./gateway/ ./protocol/proto/event/v1/`; fixtures `m7_incident_*.bpmn`; Recover + redrive cases per Story 3

**Target Platform**: single-node process (`gateway/cmd/sparrow`)

**Project Type**: Go multi-module engine + gRPC adapter

**Performance Goals**: correctness and recoverability; one COMMAND at a time per `process_instance_id`

**Constraints**: constitution (Element subject, no parallel graph, no gRPC in `processing`, no product-suite scope); `buf` FILE compatibility for proto changes

**Scale/Scope**: three user stories (open incident, resolve/retry, Recover). MVP: Service Task job failures only. Deferred: User Task incidents, ListIncidents RPC, MI-specific incident aggregation.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status |
|-----------|--------|
| I. Engine, not product suite | Pass — operator commands via existing gRPC; no console |
| II. Ledger is source of truth | Pass — INCIDENT_OPENED/RESOLVED are EVENT facts; projection rebuilds blocked state |
| III. Element is the subject | Pass — intents on SERVICE_TASK; no INCIDENT element type |
| IV. No parallel graph; thin modules | Pass — handler changes in `jobs`/`engine`/`projection`; proto in `protocol`; RPC in `gateway` |
| V. Standard BPMN | Pass — optional Sparrow extension for threshold; core behavior is job retry policy |
| VI. Serial commands, explicit rejection | Pass — same instance lock; REJECTION codes for blocked-token commands |

Post-design re-check: still pass. Fail count and blocked status derive from EVENT replay; job leases remain outside ledger.

## Project Structure

### Documentation (this feature)

```text
specs/003-incident/
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
bpmn/                     # optional sparrow:failedJobIncidentThreshold on serviceTask
protocol/proto/event/v1/  # INTENT_INCIDENT_* ; ActivityPayload.job_fail_count
protocol/proto/engine/v1/ # ResolveIncident RPC; Token blocked fields
protocol/proto/job/v1/    # FailJobRequest.no_retry
protocol/gen/go/          # generated via protocol/proto/build.sh
processing/
├── deploy/               # compile incident threshold per Service Task
├── projection/           # TokenBlocked; ApplyEvent for incident intents
├── jobs.go               # Fail threshold → INCIDENT_OPENED; skip blocked on activate
├── engine.go             # Complete rejects blocked; ResolveIncident
├── recover.go            # redrive INCIDENT_OPENED after FAILED
└── incident_test.go / testdata/m7_incident_*.bpmn
gateway/
├── engine_server.go      # ResolveIncident RPC
└── job_server.go         # FailJob no_retry mapping
```

**Structure Decision**: Extend the existing four modules. No `incident` package or runtime store for incident state.

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| None | — | — |

## Phase 0 & 1 Outputs

- [research.md](./research.md) — ledger intents, threshold default (3), resolve command, boundary/Recover behavior
- [data-model.md](./data-model.md) — token/projection fields, state transitions, rejection codes
- [contracts/engine.md](./contracts/engine.md) — proto and RPC additions
- [quickstart.md](./quickstart.md) — per-story validation steps

## Implementation Notes (for tasks phase)

1. **Proto first**: add intents 14–15, payload field, engine ResolveIncident, job `no_retry`, token fields; run `build.sh`.
2. **Projection**: `TokenBlocked`, `JobFailCount`, `IncidentErrorMessage`; ApplyEvent for FAILED (increment), INCIDENT_OPENED, INCIDENT_RESOLVED, TERMINATED (clear blocked).
3. **Fail pipeline**: after FAILED EVENT, if `fail_count >= threshold` or `no_retry`, emit INCIDENT_OPENED in same command chain.
4. **Guards**: Complete, Fail, Activate skip/reject blocked; Resolve only when blocked.
5. **Deploy**: parse optional threshold extension; default 3.
6. **Tests**: SC-001–SC-004 via quickstart fixtures; recover_test redrive case.
7. **AGENTS.md**: move incident from Remaining to implemented after implement phase.
