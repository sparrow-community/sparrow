# Engine contract

## Deploy

- Accept event sub-process start with exactly one `compensateEventDefinition` when parent is embedded SubProcess.
- Reject: process-root compensation event sub-process; multiple per SubProcess; SubProcess also has compensation boundary; non-empty start `activityRef`.

## Runtime

- Enclosing SubProcess COMPLETED → subscribe compensation handler = event sub-process id.
- Compensate throw queues Enter(event sub-process); complete event sub-process → AdvanceCompensation.
- Compensate throw inside compensation event sub-process collects subscriptions under ParentScopeID.
- No ordinary event sub-process arming for compensate kind.
