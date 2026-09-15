# Quickstart: Message and signal end

```shell
go test ./processing/ -run 'MessageEnd|SignalEnd'
go test ./processing/
```

Fixtures `processing/testdata/m20_*.bpmn`:

1. Message end alone → process completed; publishes name.
2. Message catch waiter + message end → waiter delivered.
3. Signal end alone + signal catch waiter → delivered.
4. SubProcess message/signal end → parent completes.
5. Mixed end defs → Deploy rejected.
6. Recover after delivery matches continuous.
