# Contracts: Instantiate Event-Based Gateway

No new `EngineService` or `JobService` RPCs.

## Deploy contract

- Exclusive `eventBasedGateway` MAY set `instantiate="true"` when it has **no** incoming sequence flows and **at least two** outgoing flows to intermediate catch events.
- Process MAY omit `startEvent` when the sole process entry is that instantiate exclusive EBG.
- Deploy MUST reject:
  - `instantiate` + Parallel `eventGatewayType`
  - `instantiate` with any incoming sequence flow
  - `instantiate` with fewer than two catch targets
  - any `instantiate` EBG combined with a process-level `startEvent`
  - `instantiate` nested inside a subProcess (unsupported in this increment)
- Non-instantiate event-based gateways (exclusive/parallel mid-process) unchanged.

## Runtime contract

1. `CreateInstance` on an instantiate-only deployment enters the instantiate gateway (not a startEvent).
2. Both (all) catch targets become waiting; exclusive first-wins cancels siblings.
3. `PublishMessage` / `PublishSignal` / `FireDue` do **not** create instances for instantiate deployments without a prior `CreateInstance`.
4. Mid-process exclusive/parallel EBG behavior unchanged.

## Recover

Unfinished CreateInstance COMMAND / waiting catches after Recover: winning timer or message yields the same completion path as a continuous run.

## event.v1 / engine.v1

No schema changes. First element after PROCESS start may be `TYPE_EVENT_BASED_GATEWAY` instead of `TYPE_START_EVENT` for instantiate-only processes.
