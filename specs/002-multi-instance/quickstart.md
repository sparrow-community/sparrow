# Quickstart: Multi-Instance Activities

Validate each story independently. Run from the repository root (`go.work`).

## Prerequisites

- Go 1.26.5
- Buf CLI only if proto files changed (`cd protocol/proto && ./build.sh`)

```shell
go test ./processing/ ./gateway/ ./protocol/proto/event/v1/
```

## Story 1 — Parallel multi-instance (fixed cardinality)

1. Deploy `m6_mi_parallel_cardinality.bpmn`: start → parallel MI user task (`loopCardinality` 3) → end.
2. `CreateInstance`.
3. `GetInstance`: three waiting tokens on the same activity id with distinct `loop_instance_index` 0, 1, 2.
4. `Complete` each inner token (by `token_id`).
5. Process completes; audit shows three inner COMPLETED then one host completion and a single outgoing flow taken.

Recover: after step 3, rebuild from EventLog; complete remaining inner tokens; same outcome.

## Story 2 — Sequential multi-instance

1. Deploy `m6_mi_sequential_cardinality.bpmn` (cardinality 3, `isSequential=true`).
2. Start instance: exactly one waiting token (index 0).
3. Complete index 0 → index 1 appears; complete 1 → index 2; complete 2 → process ends.

## Story 3 — Collection input

1. Deploy `m6_mi_collection.bpmn` with `loopDataInputRef=items`, `inputDataItem=item`.
2. Start with variables `{items: ["a","b"]}`.
3. Two parallel inner instances; each exposes `item` as `"a"` or `"b"`.
4. Complete both; if output mapping configured, `results` collection matches inner outputs in order.
5. Empty collection fixture: zero waiting work; process completes immediately.

## Story 4 — Completion condition (early exit)

1. Deploy `m6_mi_completion_early.bpmn`: parallel cardinality 5, `completionCondition` requiring 2 completions.
2. Start; five inner tokens waiting.
3. Complete any two; activity completes; remaining three tokens terminated (no stray waiting work).
4. Process reaches end.

## Story 5 — Multi-instance sub-process

1. Deploy `m6_mi_subprocess.bpmn`: parallel MI embedded sub-process (cardinality 2), each scope has user task.
2. Start; two sub-process scopes active (two host tokens with `ScopeHost` and distinct loop index).
3. Complete both inner user tasks; MI sub-process completes once; process ends.

## Expected test layout

Fixtures under `processing/testdata/m6_mi_*.bpmn`. Tests in `processing/multi_instance_test.go`. Gateway mapping in `gateway/engine_server_test.go` if `Token.loop_instance_index` is exposed.

## Regression

```shell
go test ./processing/ ./gateway/ ./protocol/proto/event/v1/
```

All pre-existing tests must pass; non-MI processes unchanged.
