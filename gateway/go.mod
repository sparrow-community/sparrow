module github.com/sparrow-community/sparrow/gateway

go 1.26.5

replace (
	github.com/sparrow-community/sparrow/bpmn => ../bpmn
	github.com/sparrow-community/sparrow/processing => ../processing
	github.com/sparrow-community/sparrow/protocol => ../protocol
)

require (
	github.com/sparrow-community/sparrow/processing v0.0.0-00010101000000-000000000000
	github.com/sparrow-community/sparrow/protocol v0.0.0
	google.golang.org/grpc v1.83.2
)

require (
	github.com/expr-lang/expr v1.17.8 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/sparrow-community/sparrow/bpmn v0.0.0 // indirect
	golang.org/x/net v0.58.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.41.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260526163538-3dc84a4a5aaa // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)
