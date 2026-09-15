# Engine contract: Message and signal end

No proto or gRPC changes.

## Deploy

Accepts `messageEventDefinition` or `signalEventDefinition` alone on end events. Rejects mixed definitions.

## Commands

No new commands. Publication runs inside CreateInstance / Complete enter chains via existing flush.

## Events

`TYPE_END_EVENT` lifecycle COMPLETED; optional EventPayload message/signal name. Deliveries to waiters use existing Complete / Publish paths.

## Recover

Redrive unfinished commands; buffered messages / waiters rebuild as today.
