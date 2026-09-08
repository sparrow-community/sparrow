# Quickstart: Escalation Events

```shell
go test ./processing/ -count=1 -run Escalation
go test ./processing/ ./gateway/ ./protocol/proto/event/v1/ -count=1
```

Fixtures under `processing/testdata/m15_escalation_*.bpmn`.
