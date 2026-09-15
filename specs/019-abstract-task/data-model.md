# Data model: Abstract Task

## Definition

- `element.Task` already on `FlowElements.Tasks` (`xml:"task"`).
- Activity fields (default, loopCharacteristics) apply as for Manual Task.

## Runtime

- Indexed as `Element_TYPE_TASK`.
- Waiting token until `Complete` (no JobType).
- Projection `waitingActivation` includes `TYPE_TASK`.

## No new EventLog record types
