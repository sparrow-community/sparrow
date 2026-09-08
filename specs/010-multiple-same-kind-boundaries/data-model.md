# Data Model / Contracts / Quickstart

## Data model

`WaitingBoundary` on activity ACTIVATED: boundary_id, kind, due/duration, message_name, signal_name. Token.BoundaryWaits mirrors that list.

## Contracts

No new RPCs. Deploy accepts multiple same-kind boundaries with unique message/signal names. Runtime arms all; one interrupting trigger cancels siblings.

## Quickstart

```shell
cd protocol/proto && ./build.sh
go test ./processing/ -run 'MultiBoundary|Boundary' -count=1
```

Story 1: `m14_multi_message_boundary.bpmn` — three messages; publish one.  
Story 2: `m14_multi_timer_boundary.bpmn` — PT0S vs PT1H; FireDue.  
Story 3: Recover mid dual-message wait then publish.
