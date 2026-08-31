# Quickstart: Incident (blocked execution state)

Validate each user story independently. Run from the repository root (`go.work`).

## Prerequisites

- Go 1.26.5
- Buf CLI when proto files change (`cd protocol/proto && ./build.sh`)

```shell
go test ./processing/ ./gateway/ ./protocol/proto/event/v1/
```

## Story 1 — Job failure opens an incident

1. Deploy fixture `m7_incident_service.bpmn`: start → Service Task (`job_type` set) → end.
2. `CreateInstance`.
3. `ActivateJobs` → claim job.
4. `FailJob` with retriable errors until default threshold (3): after each fail, `GetInstance` shows token `status=waiting`, `job_fail_count` incrementing.
5. Third `FailJob`: `GetInstance` shows token `status=blocked`, `incident_error_message` set.
6. `ListEvents`: three `INTENT_FAILED` then one `INTENT_INCIDENT_OPENED` for that token.
7. `Complete` or second `ActivateJobs` on blocked token → REJECTION with code `INCIDENT_OPEN`; projection unchanged.

**Non-retriable variant**: single `FailJob` with `no_retry=true` → blocked after one FAILED.

## Story 2 — Resolve incident and retry

1. Continue from Story 1 blocked state (or repeat steps 1–5).
2. `ResolveIncident` on the blocked token.
3. `GetInstance`: `status=waiting`, incident fields cleared, `job_fail_count=0`.
4. `ListEvents`: `INTENT_INCIDENT_RESOLVED` appended.
5. `ActivateJobs` → `CompleteJob` → process `status=completed`.

**Negative**: `ResolveIncident` on waiting token (no incident) → REJECTION `NO_INCIDENT`.

## Story 3 — Recover with open incident

1. Open incident (Story 1 step 5).
2. Build fresh engine: new `Engine` + `Recover` from same EventLog and deployments (mirror `processing/recover_test.go` pattern).
3. `GetInstance`: same blocked token and `incident_error_message`.
4. `ResolveIncident` → activate → complete → process completes.
5. **Redrive**: simulate crash after Fail COMMAND appended but before INCIDENT_OPENED EVENT; `Recover` finishes chain once; still exactly one INCIDENT_OPENED per fail command.

## Boundary edge (plan test)

1. Fixture with Service Task + interrupting timer boundary; fail to incident.
2. Fire timer boundary before resolve.
3. Host token terminated; no stray blocked token; process follows boundary path.

## Expected test layout

- Fixtures: `processing/testdata/m7_incident_*.bpmn`
- Tests: `processing/incident_test.go`, optional `gateway/engine_server_test.go` for ResolveIncident RPC
- Recover: extend or parallel `processing/recover_test.go` incident case

## Regression

```shell
go test ./processing/ ./gateway/ ./protocol/proto/event/v1/
```

Pre-existing job fail/retry tests (`TestJobServiceFailThenActivate`) still pass for sub-threshold fails.
