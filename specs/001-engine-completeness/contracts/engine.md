# Contracts: EngineService and event payloads

Additive Protobuf fields only. Regenerate with `cd protocol/proto && ./build.sh`. Do not renumber or reuse existing fields.

## event.v1

### `ProcessPayload`

| Field | Type | When |
|-------|------|------|
| `parent_process_instance_id` | string | child PROCESS ACTIVATED / COMPLETED / TERMINATED |
| `parent_element_id` | string | Call Activity id on the parent |
| `parent_token_id` | string | parked host token |

### `ActivityPayload`

| Field | Type | When |
|-------|------|------|
| `called_process_instance_id` | string | CALL_ACTIVITY ACTIVATED while the child exists; cleared conceptually after COMPLETED/TERMINATED (field may still appear on those records for audit) |

### `EventPayload`

| Field | Type | When |
|-------|------|------|
| `error_code` | string | already used on ERROR_THROWN; also on error ESP START_EVENT ACTIVATED (arm) |

No new `Element.Type`. No new Intent required for the happy path (reuse ACTIVATED/COMPLETED/TERMINATED). Child start is a PROCESS on the child instance id.

## engine.v1

### `Instance`

| Field | Type | Meaning |
|-------|------|---------|
| `parent_process_instance_id` | string | empty on root instances |
| `parent_element_id` | string | Call Activity id when this instance was called |
| `process_id` | string | BPMN process id bound at start |

### `Token`

| Field | Type | Meaning |
|-------|------|---------|
| `called_process_instance_id` | string | set on a waiting Call Activity host token |

### `CreateInstanceRequest`

| Field | Type | Meaning |
|-------|------|---------|
| `deployment_id` | string | existing; still valid |
| `process_id` | string | optional; start latest revision of this process id (requires empty `deployment_id` or MUST match) |
| `process_version` | int32 | optional; with `process_id`, select that revision; `0` means latest |

Rejection: `NOT_FOUND` if process id / version is missing from the store. `INVALID_ARGUMENT` if `deployment_id` and `process_id` disagree.

### `DeployResponse`

| Field | Type | Meaning |
|-------|------|---------|
| `deployment_id` | string | existing |
| `process_id` | string | primary process id in the definitions (the one `CreateInstance` would start from this deployment) |
| `process_version` | int32 | assigned revision for that process id |

`GetInstance` / `ListEvents` stay keyed by `process_instance_id`. Listing a called instance uses its own id; parent lookup is via `Instance.parent_process_instance_id`. No new ListChildren RPC in this increment.

## processing Go API (non-wire)

`CreateInstance` gains optional process-id/version selection equivalent to the proto. `Complete` / `ThrowError` / `FireDue` / `Publish*` already take `process_instance_id` and apply to caller or called without a new method.

Child start and parent resume are engine-internal; they MUST still append COMMAND+EVENT on the instance they mutate.

## Compatibility

Existing clients that only send `deployment_id` keep working. Existing projections ignore unknown proto fields. Gateway maps new fields 1:1.
