module github.com/sparrow-community/sparrow/processing

go 1.26.5

require (
	github.com/google/uuid v1.6.0
	github.com/sparrow-community/sparrow/bpmn v0.0.0
	github.com/sparrow-community/sparrow/protocol v0.0.0
)

require github.com/expr-lang/expr v1.17.8

require (
	golang.org/x/net v0.57.0 // indirect
	golang.org/x/text v0.40.0 // indirect
	google.golang.org/protobuf v1.36.11
)

replace (
	github.com/sparrow-community/sparrow/bpmn => ../bpmn
	github.com/sparrow-community/sparrow/protocol => ../protocol
)
