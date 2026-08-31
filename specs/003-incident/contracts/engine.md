# Contracts: Incident lifecycle

Additive Protobuf fields and RPCs only. Regenerate with `cd protocol/proto && ./build.sh`. Do not renumber or reuse existing fields.

## event.v1

### `Element.Intent` (extend enum)

| Value | Name | Meaning |
|-------|------|---------|
| 14 | `INTENT_INCIDENT_OPENED` | Service Task token blocked after retry policy exhausted |
| 15 | `INTENT_INCIDENT_RESOLVED` | Operator resolve; token returns to waiting job semantics |

No new `Element.Type`.

### `ActivityPayload` (extend)

| Field | Type | When |
|-------|------|------|
| `job_fail_count` | int32 | FAILED (current consecutive fail count); INCIDENT_OPENED (count at open) |

Existing `error_message` used on FAILED and INCIDENT_OPENED.

## engine.v1

### `Token` (extend)

| Field | Type | Meaning |
|-------|------|---------|
| `status` | string | Existing values plus **`blocked`** |
| `incident_error_message` | string | Last worker error when blocked; empty otherwise |
| `job_fail_count` | int32 | Consecutive fails since last activation/resolve; informational on GetInstance |

### `ResolveIncidentRequest`

| Field | Type | Required |
|-------|------|----------|
| `process_instance_id` | string | yes |
| `element_id` | string | yes |
| `token_id` | string | yes |

### `ResolveIncidentResponse`

Empty on success.

### `EngineService` (extend)

```protobuf
rpc ResolveIncident(ResolveIncidentRequest) returns (ResolveIncidentResponse);
```

Maps to `processing.Engine.ResolveIncident`.

## job.v1

### `FailJobRequest` (extend)

| Field | Type | Meaning |
|-------|------|---------|
| `no_retry` | bool | When true, open incident after this fail regardless of fail count |

Default false — threshold policy applies.

No change to Activate/Complete/Heartbeat semantics except: **blocked tokens are not returned from ActivateJobs**.

## processing Go API (non-wire)

```go
func (e *Engine) ResolveIncident(ctx context.Context, instanceID, elementID, tokenID string) error
// Fail signature gains no_retry bool (or FailOption); internal threshold from deploy spec
```

- `deploy.ServiceTaskSpec` (or equivalent) carries compiled `IncidentThreshold`.
- `projection.ApplyEvent` handles INCIDENT_OPENED / INCIDENT_RESOLVED.
- `jobs.go` Fail: after FAILED emit, if policy triggers, emit INCIDENT_OPENED.
- `engine.go` Complete: reject when `TokenBlocked`.
- Job activation queue skips blocked tokens.

## Compatibility

- Existing clients ignore new enum values and fields.
- Processes without failures behave unchanged.
- `GetInstance` clients that only check `status == "waiting"` must treat **`blocked`** separately (document in quickstart).

## Out of scope (this increment)

- ListIncidents RPC (use GetInstance tokens + ListEvents filter)
- Incidents on User Task, ThrowError, or non-job failures
- Incident on Call Activity host (unless later specified for child job failures — child instance is separate)
- Operator UI / product suite
- Cross-deployment CallActivity
