module github.com/sparrow-community/sparrow/processing

go 1.26.5

require (
	github.com/sony/sonyflake/v2 v2.2.0
	github.com/sparrow-community/sparrow/bpmn v0.0.0
	github.com/sparrow-community/sparrow/protocol v0.0.0
)

require google.golang.org/protobuf v1.36.11 // indirect

replace (
	github.com/sparrow-community/sparrow/bpmn => ../bpmn
	github.com/sparrow-community/sparrow/protocol => ../protocol
)
