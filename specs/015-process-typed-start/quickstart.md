# Quickstart: Process-level typed start

```shell
go test ./processing/ -run 'TypedStart|MessageStart|TimerStart|SignalStart|ConditionalStart'
go test ./processing/
```

Fixtures under `processing/testdata/m19_*.bpmn`:

1. Message start only → CreateInstance rejected; PublishMessage `order.created` → wait User Task → Complete → completed.
2. Timer start (e.g. PT0S) → CreateInstance rejected; FireDue → completed (or next wait).
3. Signal start → PublishSignal → progresses.
4. Conditional start → EvaluateConditionalStarts false → no instance; true → progresses.
5. None + message alternatives → CreateInstance via none; unmatched PublishMessage via message start.
6. Process-level error start → Deploy rejected; ESP error suite still green.
7. Recover: message/timer/signal start paths then finish waits match continuous run.
8. Instantiate EBG fixture still requires CreateInstance (no auto-create).
