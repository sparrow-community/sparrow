# Data Model

| Concept | Representation |
|---------|----------------|
| Unfinished Call Activity | Caller token on Call Activity with `CalledProcessInstanceID` set and child StatusActive |
| PublicationCompensateUnfinishedChild | Engine publication after parent startCompensation |
| PendingCompensation.WaitingChildren | Count of unfinished-call children still compensating |
| Child PendingCompensation | Queue of child handler ids; empty ThrowTokenID + NotifyParentUnfinishedCall |

No proto changes.
