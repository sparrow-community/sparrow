# Data Model: Terminate End and Remaining Tasks

## Terminate end

- Indexed as `TYPE_END_EVENT` with `terminateEnds[id]=true`.
- Enter: full lifecycle COMPLETED on the end, then cancel other tokens in `ScopeOf(end)` (keep current token and SubProcess host), `PublicationTerminateChild` for Call Activity children, then `TryCompleteProcess`.

## Tasks

| Type | Wait | Payload |
|------|------|---------|
| `TYPE_MANUAL_TASK` | Complete | ActivityPayload + boundaries |
| `TYPE_RECEIVE_TASK` | PublishMessage | ActivityPayload.MessageName + boundaries |
| `TYPE_SEND_TASK` | none | publish message name |
| `TYPE_BUSINESS_RULE_TASK` | job | ActivityPayload.JobType + boundaries |

## Projection

- `waitingActivation` includes Manual, Receive, Business Rule (and Send only if MI host waits).
- Receive: `tok.MessageName` is the task message, distinct from boundary waits.
- Job claim: any waiting token with `JobType` (already true); notify on BR ACTIVATED.

## Validation

- Abstract `bpmn:task` rejected.
- Receive `instantiate` rejected.
- End with terminate plus other event definitions rejected.
