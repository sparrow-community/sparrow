# Quickstart: Terminate and remaining tasks

```shell
go test ./processing/ -run 'Terminate|ManualTask|ReceiveTask|SendTask|BusinessRule'
go test ./processing/
```

Fixtures under `processing/testdata/m18_*.bpmn`:

1. Parallel User Task vs terminate end → process completed, User Task TERMINATED.
2. Manual Task → Complete → completed.
3. Receive Task → PublishMessage `order.confirmed` → completed.
4. Send Task publishes to a waiting Receive/catch.
5. Business Rule `decide.v1` → Activate → Complete.
6. Recover variants for Manual / Receive / Business Rule.
