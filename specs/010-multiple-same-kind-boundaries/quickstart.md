# Quickstart: Multiple Same-Kind Boundaries

```shell
cd protocol/proto && ./build.sh
go test ./processing/ -run 'MultiBoundary|Boundary|SignalBoundary' -count=1
go test ./processing/ ./gateway/ ./protocol/proto/event/v1/
```
