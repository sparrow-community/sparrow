# Data Model: Incident (blocked execution state)

## Entities

### Incident policy (compile-time, in `deploy`)

| Field | Notes |
|-------|--------|
| `element_id` | Service Task id |
| `incident_threshold` | Positive int; default 3 when BPMN extension absent |

Validation: only on `SERVICE_TASK`; threshold ≥ 1 when extension present.

BPMN extension (optional): `sparrow:failedJobIncidentThreshold="N"` on `serviceTask`.

Engine constant: `defaultIncidentThreshold = 3`.

### Token (runtime projection, extended)

| Field | Notes |
|-------|--------|
| `Status` | `active`, `waiting`, or **`blocked`** |
| `JobType` | Preserved while blocked (for post-resolve activate) |
| `JobFailCount` | Consecutive FAILED count since last ACTIVATED or INCIDENT_RESOLVED |
| `IncidentErrorMessage` | Error text snapshot at INCIDENT_OPENED; cleared on resolve |
| existing fields | Unchanged (boundaries, MI index, etc.) |

Invariant: at most one open incident per token (`Status == blocked`).

### ActivityPayload (ledger, extended)

| Field | Notes |
|-------|--------|
| `job_fail_count` | int32; set on FAILED (current count) and INCIDENT_OPENED (count at open) |
| `error_message` | Existing; worker reason on FAILED and copied on INCIDENT_OPENED |

No new payload type; reuse `ActivityPayload` on `SERVICE_TASK`.

### Element intents (ledger, new)

| Intent | When |
|--------|------|
| `INTENT_INCIDENT_OPENED` | After threshold fail or `no_retry` fail; token → blocked |
| `INTENT_INCIDENT_RESOLVED` | Operator resolve; token → waiting, fail count reset |

Existing `INTENT_FAILED` unchanged — emitted before INCIDENT_OPENED when threshold crossed.

## State transitions

### Service Task job — retriable fail below threshold

```text
ACTIVATED (waiting, job_type set, fail_count=0)
  → Fail → FAILED (fail_count=1, still waiting)
  → Activate → … → Fail → FAILED (fail_count=2, still waiting)
  → Activate → Complete → COMPLETED → flow continues
```

### Service Task job — incident on threshold

```text
ACTIVATED (waiting, fail_count=0)
  → Fail × (threshold-1) → FAILED each time, waiting, fail_count increments
  → Fail (Nth) → FAILED → INCIDENT_OPENED (blocked, incident_error set)
  → Activate → REJECTION INCIDENT_OPEN
  → Complete → REJECTION INCIDENT_OPEN
  → ResolveIncident → INCIDENT_RESOLVED (waiting, fail_count=0)
  → Activate → Complete → COMPLETED
```

### Non-retriable fail

```text
ACTIVATED (waiting)
  → Fail(no_retry=true) → FAILED → INCIDENT_OPENED (blocked)
```

### Boundary on blocked host

```text
blocked on SERVICE_TASK
  → boundary interrupt → host TERMINATING/TERMINATED
  → token removed / no longer blocked (no INCIDENT_RESOLVED EVENT)
```

## Validation rules

| Rule | Enforcement |
|------|-------------|
| Incident only on job-backed waiting/blocked SERVICE_TASK | Fail/Resolve handlers |
| Fail on User Task | Existing REJECTION; no incident |
| Resolve without open incident | REJECTION `NO_INCIDENT` |
| Double resolve | REJECTION `NO_INCIDENT` or `INVALID_STATE` |
| Complete/Activate on blocked token | REJECTION `INCIDENT_OPEN` |
| Fail on blocked token | REJECTION `INCIDENT_OPEN` |
| Resolve on completed/terminated instance | REJECTION `INVALID_STATE` |

## Rejection codes (stable)

| Code | Meaning |
|------|---------|
| `INCIDENT_OPEN` | Command not allowed while token blocked |
| `NO_INCIDENT` | Resolve when token is not blocked |
| `INVALID_STATE` | Instance ended or token not waiting/blocked as required |
| `NOT_FOUND` | Unknown instance/element/token |
| Existing codes | Unchanged for non-incident paths |

## Recover rebuild

Replay order per token:

1. `ACTIVATED` → waiting, fail_count=0, clear incident fields
2. `FAILED` → if not blocked, waiting, increment fail_count
3. `INCIDENT_OPENED` → blocked, set incident_error from payload
4. `INCIDENT_RESOLVED` → waiting, fail_count=0, clear incident fields
5. `TERMINATED` / token delete → drop blocked state

Redrive: Fail command emits FAILED then optionally INCIDENT_OPENED once per `source_record_id`.
