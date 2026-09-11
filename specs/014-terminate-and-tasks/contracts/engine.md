# Engine contract: Terminate and remaining tasks

No proto or gRPC changes.

## Deploy

Accepts `terminateEventDefinition` on end events (exclusive). Accepts `manualTask`, `receiveTask`, `sendTask`, `businessRuleTask`. Rejects abstract `task` and `receiveTask instantiate="true"`.

## Commands

| Situation | Command |
|-----------|---------|
| Manual Task | `Complete` |
| Receive Task | `PublishMessage` (or `Complete` on the waiting token) |
| Send Task | none (runs on Enter) |
| Business Rule Task | `Activate` → `Complete` / `Fail` |
| Terminate end | none (runs on Enter) |

## Events

Distinct `Element.Type` per task. Terminate end uses `TYPE_END_EVENT` COMPLETED; cancelled peers `TERMINATING`/`TERMINATED`. Process after process-level terminate: `TYPE_PROCESS` COMPLETED.

## Recover

Redrive unfinished COMMANDs; waiting Manual/Receive/BR tokens rebuild from EVENTs.
