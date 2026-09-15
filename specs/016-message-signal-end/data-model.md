# Data Model: Message End and Signal End

| Kind | Index | OnEnter |
|------|-------|---------|
| Message end | `messageEnds[id]=name` | Publish message + complete end + TryCompleteProcess |
| Signal end | `signalEnds[id]=name` | Publish signal + complete end + TryCompleteProcess |

Payload on ACTIVATED/COMPLETED may carry `MessageName` / `SignalName` like throws.

Validation: exactly one matching event definition; no timers or other catch defs.
