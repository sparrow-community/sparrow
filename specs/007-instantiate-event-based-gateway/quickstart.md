# Quickstart: Instantiate Event-Based Gateway

Validate each user story independently. Run from the repository root (`go.work`).

## Prerequisites

- Go 1.26.5
- No proto changes expected

```shell
go test ./processing/ ./gateway/ ./protocol/proto/event/v1/
```

## Story 1 — Deploy and start via instantiate exclusive EBG

1. Deploy `m11_instantiate_ebg_timer_message.bpmn` (no startEvent; instantiate EBG → timer + message catches → ends).
2. `CreateInstance` → two waiting catches.
3. `FireDue` **or** `PublishMessage` → one path completes; sibling terminated; process completed.

## Story 2 — Reject invalid instantiate configurations

Deploy each invalid fixture and expect `UNSUPPORTED_ELEMENT` (or clear deploy error):

- Incoming flow on instantiate EBG
- Parallel instantiate EBG
- startEvent + instantiate EBG
- Single outgoing catch

## Story 3 — Recover mid-race

1. CreateInstance on happy-path fixture; both catches waiting.
2. `Recover` / `Open` from EventLog + deployments.
3. PublishMessage (or FireDue) → same outcome as Story 1 without recover.

## Regression

```shell
go test ./processing/ -run 'EventBased|Instantiate' -count=1
go test ./processing/ ./gateway/ ./protocol/proto/event/v1/
```

## Expected test layout

- `processing/testdata/m11_instantiate_ebg_*.bpmn`
- `processing/event_based_gateway_test.go` (or dedicated instantiate tests)
- Recover coverage in `recover_test.go` or same package
