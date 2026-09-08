# Quickstart: Nested Event Sub-Process

```shell
go test ./processing/ -run 'EventSubProcess|EspInEsp|Nested' -count=1
```

## Story 1 — ESP in ESP interrupting

1. Deploy `m13_esp_in_esp_message.bpmn`.
2. CreateInstance → PublishMessage outer → waiting in outer; nested armed.
3. PublishMessage inner → nested completes; outer work terminated.

## Story 2 — Disarm on outer complete

1. Same fixture; complete outer path without inner message.
2. Nested arm cleared; PublishMessage inner delivers 0.

## Story 3 — Recover

1. After outer trigger with nested armed, Open/Recover.
2. PublishMessage inner → completed as Story 1.

## Regression

```shell
go test ./processing/ -run 'EventSubProcess|ErrorEventSubProcess' -count=1
go test ./processing/ ./gateway/ ./protocol/proto/event/v1/
```
