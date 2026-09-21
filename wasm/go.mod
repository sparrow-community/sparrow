module github.com/sparrow-community/sparrow/wasm

go 1.26.5

replace (
	github.com/sparrow-community/sparrow/bpmn => ../bpmn
	github.com/sparrow-community/sparrow/processing => ../processing
	github.com/sparrow-community/sparrow/protocol => ../protocol
)

require (
	github.com/sparrow-community/sparrow/processing v0.0.0-00010101000000-000000000000
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/expr-lang/expr v1.17.8 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/sparrow-community/sparrow/bpmn v0.0.0 // indirect
	github.com/sparrow-community/sparrow/protocol v0.0.0 // indirect
	golang.org/x/net v0.57.0 // indirect
	golang.org/x/text v0.40.0 // indirect
)
