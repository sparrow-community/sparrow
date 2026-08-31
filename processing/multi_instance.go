package processing

import (
	"encoding/json"

	"github.com/sparrow-community/sparrow/processing/deploy"
	"github.com/sparrow-community/sparrow/processing/projection"
)

// AppendOutputCollection writes assembled output items to the instance variable map.
func AppendOutputCollection(inst *projection.Instance, spec deploy.MultiInstanceSpec, loop *projection.MultiInstanceLoop) {
	if inst == nil || loop == nil || spec.OutputCollectionVar == "" {
		return
	}
	b, err := json.Marshal(loop.OutputItems)
	if err != nil {
		return
	}
	if inst.Variables == nil {
		inst.Variables = make(map[string]string)
	}
	inst.Variables[spec.OutputCollectionVar] = string(b)
}
