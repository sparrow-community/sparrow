# Data Model: Process-Level Typed Start

## Start classification (process-level only)

| Kind | Event definition | CreateInstance | Mint trigger |
|------|------------------|----------------|--------------|
| None | no event definitions | yes | CreateInstance |
| Message | exactly one messageEventDefinition | no (alone) | PublishMessage (unmatched) |
| Timer | exactly one timerEventDefinition | no (alone) | FireDue when schedule due |
| Signal | exactly one signalEventDefinition | no (alone) | PublishSignal (unmatched) |
| Conditional | exactly one conditionalEventDefinition | no (alone) | EvaluateConditionalStarts |
| Error | errorEventDefinition | — | Deploy reject |
| Other / mixed defs | — | — | Deploy reject |

Embedded SubProcess none starts unchanged. Event Sub-Process typed starts unchanged (existing ESP index).

## Deployment indexes

- `noneStartID` — process-level none start id if present
- `messageStarts[name] → []{deployment-local startEventID}` (or list of start ids on this deployment)
- `signalStarts[name] → []startEventID`
- `timerStarts[startEventID] → TimerCatchSpec` (reuse timer catch parsing)
- `conditionalStarts[startEventID] → condition expression text`

Instantiate EBG entry remains via existing `StartEventID` / `instantiateEntryID` when `StartEvents` empty — **not** registered in typed start maps.

## Runtime helper (non-ledger)

- **Timer start schedule**: `{deploymentID, startEventID, dueUnixMs}` armed at Deploy / Recover; updated on cycle re-arm; removed after one-shot fire.
- Optional `runtime.Store` persistence of schedules (same class as message buffer).

## Instance mint

Shared path produces the same COMMAND/EVENT chain as today’s CreateInstance, but `Enter` targets the chosen `startEventID` (none or typed). Variables from CreateInstance / publish / evaluate become process variables.

## Validation

- Process-level error / escalation / compensate / link / terminate start → reject.
- Exactly one event definition for each typed start.
- Call Activity called process must have none start or instantiate EBG.
- Cannot combine any startEvent with instantiate EBG (existing rule).
