# Engine contract: Process-level typed start

No proto change required for MVP (FireDue / typed mint are processing APIs). Gateway may later expose EvaluateConditionalStarts.

## Deploy

| Start | Result |
|-------|--------|
| None | Accepted (CreateInstance entry) |
| Message / timer / signal / conditional | Accepted; indexed as start subscriptions |
| Error (process-level) | `UNSUPPORTED_ELEMENT` |
| Mixed / unknown defs on start | `UNSUPPORTED_ELEMENT` |
| Instantiate EBG + any startEvent | Reject (existing) |
| Call Activity → typed-only called process | Reject |

## Commands

| Situation | Command |
|-----------|---------|
| None start (or none among alternatives) | `CreateInstance` |
| Instantiate EBG only | `CreateInstance` (unchanged; no auto-create from publish/fire) |
| Message start | `PublishMessage` when no waiter/ESP/scope delivery |
| Signal start | `PublishSignal` when no waiter/ESP delivery |
| Timer start | `FireDue` when deployment schedule due |
| Conditional start | `EvaluateConditionalStarts` (deployment/process + variables) |
| Typed-only + CreateInstance | Reject |

## Events

Minted instances use existing `TYPE_PROCESS` + `TYPE_START_EVENT` lifecycle. No new Element.Type. Timer-start arms are not Events.

## Recover

- Redrive unfinished CreateInstance / mint COMMANDs as today.
- Rebuild timer-start schedules from deployments + runtime store.
- Message/signal start creation is command-driven (Publish*); after Recover, the next publish/fire/evaluate matches continuous outcomes for waiting work.
