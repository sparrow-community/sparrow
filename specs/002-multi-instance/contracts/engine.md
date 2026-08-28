# Contracts: Multi-Instance loop fields

Additive Protobuf fields only. Regenerate with `cd protocol/proto && ./build.sh`. Do not renumber or reuse existing fields.

## event.v1

### `ActivityPayload` (extend)

| Field | Type | When |
|-------|------|------|
| `loop_instance_index` | int32 | Inner iteration ACTIVATED/COMPLETED/TERMINATED; 0-based. Omit/zero for non-multi-instance activities. |
| `loop_total_instances` | int32 | Optional snapshot on host ACTIVATED/COMPLETED for audit |
| `loop_completed_instances` | int32 | Optional snapshot on host COMPLETED |

No new `Element.Type`. No new Intent. Inner and host lifecycle reuse ACTIVATED/COMPLETED/TERMINATED.

## engine.v1

### `Token` (extend)

| Field | Type | Meaning |
|-------|------|---------|
| `loop_instance_index` | int32 | Set on inner-instance waiting tokens; -1 or unset for ordinary tokens |

### `CompleteRequest`

No new RPC. Clients complete inner instances by `token_id` (already required). `loop_instance_index` on the token in `GetInstance` is informational for operators.

Optional future: `loop_instance_index` on `CompleteRequest` as a guard that the token matches the claimed index — not required for v1 if `token_id` is authoritative.

## processing Go API (non-wire)

- `deploy` compiles `MultiInstanceSpec` per activity/sub-process.
- Handlers return Effects to spawn inner tokens, increment counters, cancel stragglers.
- `projection.ApplyEvent` maintains `MultiInstanceLoops` and `Token.LoopInstanceIndex`.

## Compatibility

Existing clients ignore new fields. Non-MI deployments unchanged. `GetInstance` gains optional `loop_instance_index` on tokens; clients that only use `token_id` for `Complete` keep working.

## Out of scope (this increment)

- New ListLoops RPC
- MI-specific Job API (jobs remain per inner Service Task token)
- Cross-instance MI (N/A — all inner instances share one `process_instance_id`)
