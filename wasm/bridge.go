// Copyright 2026 The Sparrow community and contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build js && wasm

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"syscall/js"
	"time"

	"github.com/sparrow-community/sparrow/processing"
	"google.golang.org/protobuf/encoding/protojson"
)

// api is a thin JS façade over processing.Engine.
// Timers: host reads nextDueUnixMs and schedules FireDue (e.g. setTimeout).
// Jobs / scripts: host Activate(wait=0) and runs workers in JS, then Complete/Fail.
type api struct {
	eng *processing.Engine
}

func newAPI(eng *processing.Engine) *api {
	return &api{eng: eng}
}

func (a *api) value() js.Value {
	o := js.Global().Get("Object").New()
	o.Set("deploy", a.fn(a.deploy))
	o.Set("createInstance", a.fn(a.createInstance))
	o.Set("complete", a.fn(a.complete))
	o.Set("throwError", a.fn(a.throwError))
	o.Set("resolveIncident", a.fn(a.resolveIncident))
	o.Set("fireDue", a.fn(a.fireDue))
	o.Set("nextDueUnixMs", a.fn(a.nextDueUnixMs))
	o.Set("setNowUnixMs", a.fn(a.setNowUnixMs))
	o.Set("publishMessage", a.fn(a.publishMessage))
	o.Set("publishSignal", a.fn(a.publishSignal))
	o.Set("evaluateConditions", a.fn(a.evaluateConditions))
	o.Set("evaluateConditionalStarts", a.fn(a.evaluateConditionalStarts))
	o.Set("activate", a.fn(a.activate))
	o.Set("fail", a.fn(a.fail))
	o.Set("heartbeat", a.fn(a.heartbeat))
	o.Set("getDeployment", a.fn(a.getDeployment))
	o.Set("getInstance", a.fn(a.getInstance))
	o.Set("listInstanceIds", a.fn(a.listInstanceIds))
	o.Set("listEvents", a.fn(a.listEvents))
	o.Set("enableIntervention", a.fn(a.enableIntervention))
	o.Set("disableIntervention", a.fn(a.disableIntervention))
	o.Set("setBreakpoints", a.fn(a.setBreakpoints))
	o.Set("continueIntervention", a.fn(a.continueIntervention))
	o.Set("stepInto", a.fn(a.stepInto))
	o.Set("stepOver", a.fn(a.stepOver))
	o.Set("getInterventionState", a.fn(a.getInterventionState))
	o.Set("setVariables", a.fn(a.setVariables))
	return o
}

func (a *api) fn(fn func(js.Value, []js.Value) (any, error)) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		// Never panic out of a JS callback: that exits the Go WASM runtime
		// ("Go program has already exited") and breaks later calls.
		var out any
		var err error
		func() {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("%v", r)
				}
			}()
			out, err = fn(this, args)
		}()
		if err != nil {
			return toJS(map[string]any{"$error": err.Error()})
		}
		return toJS(out)
	})
}

func (a *api) deploy(_ js.Value, args []js.Value) (any, error) {
	xml, err := argString(args, 0, "bpmnXml")
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(xml) == "" {
		return nil, fmt.Errorf("INVALID_ARGUMENT: bpmnXml is empty")
	}
	id, err := a.eng.Deploy(context.Background(), []byte(xml))
	if err != nil {
		return nil, err
	}
	processID := ""
	if dep, ok := a.eng.GetDeployment(id); ok {
		processID = dep.ProcessID()
	}
	return map[string]any{"deploymentId": id, "processId": processID}, nil
}

func (a *api) createInstance(_ js.Value, args []js.Value) (any, error) {
	req, err := argObject(args, 0)
	if err != nil {
		return nil, err
	}
	vars, err := asAnyMap(req["variables"])
	if err != nil {
		return nil, fmt.Errorf("INVALID_ARGUMENT: variables: %w", err)
	}
	id, err := a.eng.CreateInstanceRequest(context.Background(), processing.CreateInstanceRequest{
		DeploymentID:   asString(req["deploymentId"]),
		ProcessID:      asString(req["processId"]),
		ProcessVersion: asInt32(req["processVersion"]),
		Variables:      vars,
	})
	if err != nil {
		return nil, err
	}
	return map[string]any{"instanceId": id}, nil
}

func (a *api) complete(_ js.Value, args []js.Value) (any, error) {
	req, err := argObject(args, 0)
	if err != nil {
		return nil, err
	}
	vars, err := asAnyMap(req["variables"])
	if err != nil {
		return nil, fmt.Errorf("INVALID_ARGUMENT: variables: %w", err)
	}
	err = a.eng.Complete(context.Background(),
		asString(req["instanceId"]),
		asString(req["elementId"]),
		asString(req["tokenId"]),
		vars,
	)
	return map[string]any{"ok": err == nil}, err
}

func (a *api) throwError(_ js.Value, args []js.Value) (any, error) {
	req, err := argObject(args, 0)
	if err != nil {
		return nil, err
	}
	err = a.eng.ThrowError(context.Background(),
		asString(req["instanceId"]),
		asString(req["elementId"]),
		asString(req["tokenId"]),
		asString(req["errorCode"]),
	)
	return map[string]any{"ok": err == nil}, err
}

func (a *api) resolveIncident(_ js.Value, args []js.Value) (any, error) {
	req, err := argObject(args, 0)
	if err != nil {
		return nil, err
	}
	err = a.eng.ResolveIncident(context.Background(),
		asString(req["instanceId"]),
		asString(req["elementId"]),
		asString(req["tokenId"]),
	)
	return map[string]any{"ok": err == nil}, err
}

func (a *api) fireDue(_ js.Value, _ []js.Value) (any, error) {
	err := a.eng.FireDue(context.Background())
	return map[string]any{"ok": err == nil}, err
}

func (a *api) nextDueUnixMs(_ js.Value, _ []js.Value) (any, error) {
	return a.eng.NextDueUnixMs(), nil
}

func (a *api) setNowUnixMs(_ js.Value, args []js.Value) (any, error) {
	if len(args) == 0 || args[0].IsNull() || args[0].IsUndefined() {
		a.eng.SetNow(nil)
		return map[string]any{"ok": true}, nil
	}
	ms := int64(args[0].Float())
	if ms <= 0 {
		a.eng.SetNow(nil)
		return map[string]any{"ok": true}, nil
	}
	a.eng.SetNow(func() time.Time { return time.UnixMilli(ms) })
	return map[string]any{"ok": true, "nowUnixMs": ms}, nil
}

func (a *api) publishMessage(_ js.Value, args []js.Value) (any, error) {
	req, err := argObject(args, 0)
	if err != nil {
		return nil, err
	}
	keys, err := asAnyMap(req["correlationKeys"])
	if err != nil {
		return nil, fmt.Errorf("INVALID_ARGUMENT: correlationKeys: %w", err)
	}
	vars, err := asAnyMap(req["variables"])
	if err != nil {
		return nil, fmt.Errorf("INVALID_ARGUMENT: variables: %w", err)
	}
	n, err := a.eng.PublishMessage(context.Background(), processing.PublishMessageRequest{
		Name:              asString(req["name"]),
		ProcessInstanceID: asString(req["instanceId"]),
		CorrelationKeys:   keys,
		Variables:         vars,
	})
	return map[string]any{"delivered": n}, err
}

func (a *api) publishSignal(_ js.Value, args []js.Value) (any, error) {
	req, err := argObject(args, 0)
	if err != nil {
		return nil, err
	}
	vars, err := asAnyMap(req["variables"])
	if err != nil {
		return nil, fmt.Errorf("INVALID_ARGUMENT: variables: %w", err)
	}
	n, err := a.eng.PublishSignal(context.Background(), processing.PublishSignalRequest{
		Name:              asString(req["name"]),
		ProcessInstanceID: asString(req["instanceId"]),
		Variables:         vars,
	})
	return map[string]any{"delivered": n}, err
}

func (a *api) evaluateConditions(_ js.Value, args []js.Value) (any, error) {
	req, err := argObject(args, 0)
	if err != nil {
		return nil, err
	}
	vars, err := asAnyMap(req["variables"])
	if err != nil {
		return nil, fmt.Errorf("INVALID_ARGUMENT: variables: %w", err)
	}
	n, err := a.eng.EvaluateConditions(context.Background(), processing.EvaluateConditionsRequest{
		ProcessInstanceID: asString(req["instanceId"]),
		Variables:         vars,
	})
	return map[string]any{"fired": n}, err
}

func (a *api) evaluateConditionalStarts(_ js.Value, args []js.Value) (any, error) {
	req, err := argObject(args, 0)
	if err != nil {
		return nil, err
	}
	vars, err := asAnyMap(req["variables"])
	if err != nil {
		return nil, fmt.Errorf("INVALID_ARGUMENT: variables: %w", err)
	}
	n, err := a.eng.EvaluateConditionalStarts(context.Background(), processing.EvaluateConditionalStartsRequest{
		DeploymentID:   asString(req["deploymentId"]),
		ProcessID:      asString(req["processId"]),
		ProcessVersion: asInt32(req["processVersion"]),
		Variables:      vars,
	})
	return map[string]any{"started": n}, err
}

func (a *api) activate(_ js.Value, args []js.Value) (any, error) {
	req, err := argObject(args, 0)
	if err != nil {
		return nil, err
	}
	// Browser host must not long-poll on the JS thread; wait is forced to 0.
	jobs, err := a.eng.Activate(context.Background(), processing.ActivateRequest{
		JobType:      asString(req["jobType"]),
		MaxJobs:      asInt(req["maxJobs"]),
		Wait:         0,
		WorkerID:     asString(req["workerId"]),
		LockDuration: time.Duration(asInt64(req["lockDurationMs"])) * time.Millisecond,
	})
	if err != nil {
		return nil, err
	}
	out := make([]map[string]any, 0, len(jobs))
	for _, j := range jobs {
		out = append(out, map[string]any{
			"jobType":      j.JobType,
			"instanceId":   j.ProcessInstanceID,
			"deploymentId": j.DeploymentID,
			"elementId":    j.ElementID,
			"tokenId":      j.TokenID,
			"variables":    j.Variables,
			"workerId":     j.WorkerID,
			"lockDeadline": j.LockDeadline.UnixMilli(),
			"scriptFormat": j.ScriptFormat,
			"script":       j.Script,
		})
	}
	return map[string]any{"jobs": out}, nil
}

func (a *api) fail(_ js.Value, args []js.Value) (any, error) {
	req, err := argObject(args, 0)
	if err != nil {
		return nil, err
	}
	err = a.eng.Fail(context.Background(),
		asString(req["instanceId"]),
		asString(req["elementId"]),
		asString(req["tokenId"]),
		asString(req["message"]),
		asBool(req["noRetry"]),
	)
	return map[string]any{"ok": err == nil}, err
}

func (a *api) heartbeat(_ js.Value, args []js.Value) (any, error) {
	req, err := argObject(args, 0)
	if err != nil {
		return nil, err
	}
	err = a.eng.Heartbeat(context.Background(),
		asString(req["instanceId"]),
		asString(req["tokenId"]),
		asString(req["workerId"]),
		time.Duration(asInt64(req["lockDurationMs"]))*time.Millisecond,
	)
	return map[string]any{"ok": err == nil}, err
}

func (a *api) getDeployment(_ js.Value, args []js.Value) (any, error) {
	id, err := argString(args, 0, "deploymentId")
	if err != nil {
		return nil, err
	}
	dep, ok := a.eng.GetDeployment(id)
	if !ok {
		return nil, fmt.Errorf("NOT_FOUND: deployment %q", id)
	}
	xml, err := a.eng.GetDeploymentXML(id)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"deploymentId": id,
		"processId":    dep.ProcessID(),
		"bpmnXml":      string(xml),
	}, nil
}

func (a *api) getInstance(_ js.Value, args []js.Value) (any, error) {
	id, err := argString(args, 0, "instanceId")
	if err != nil {
		return nil, err
	}
	inst, ok := a.eng.GetInstance(id)
	if !ok {
		return nil, fmt.Errorf("NOT_FOUND: instance %q", id)
	}
	return instanceJSON(inst), nil
}

func (a *api) listInstanceIds(_ js.Value, _ []js.Value) (any, error) {
	return a.eng.ListInstanceIDs(), nil
}

func (a *api) listEvents(_ js.Value, args []js.Value) (any, error) {
	id, err := argString(args, 0, "instanceId")
	if err != nil {
		return nil, err
	}
	events, err := a.eng.ListEvents(context.Background(), id)
	if err != nil {
		return nil, err
	}
	marshaler := protojson.MarshalOptions{EmitUnpopulated: false}
	out := make([]json.RawMessage, 0, len(events))
	for _, ev := range events {
		b, err := marshaler.Marshal(ev)
		if err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	// Return as parsed JSON array via round-trip so toJS sees []any.
	raw, err := json.Marshal(out)
	if err != nil {
		return nil, err
	}
	var arr []any
	if err := json.Unmarshal(raw, &arr); err != nil {
		return nil, err
	}
	return map[string]any{"events": arr}, nil
}

func (a *api) enableIntervention(_ js.Value, args []js.Value) (any, error) {
	req, err := argObject(args, 0)
	if err != nil {
		return nil, err
	}
	_, err = a.eng.EnableIntervention(context.Background(), processing.EnableInterventionRequest{
		InstanceID: asString(req["instanceId"]),
		Policy:     processing.RunPolicy(asString(req["policy"])),
	})
	return map[string]any{"ok": err == nil}, err
}

func (a *api) disableIntervention(_ js.Value, _ []js.Value) (any, error) {
	_, err := a.eng.DisableIntervention(context.Background(), processing.DisableInterventionRequest{})
	return map[string]any{"ok": err == nil}, err
}

func (a *api) setBreakpoints(_ js.Value, args []js.Value) (any, error) {
	req, err := argObject(args, 0)
	if err != nil {
		return nil, err
	}
	ids, err := asStringSlice(req["elementIds"])
	if err != nil {
		return nil, fmt.Errorf("INVALID_ARGUMENT: elementIds: %w", err)
	}
	resp, err := a.eng.SetBreakpoints(context.Background(), processing.SetBreakpointsRequest{
		InstanceID: asString(req["instanceId"]),
		ElementIDs: ids,
	})
	if err != nil {
		return nil, err
	}
	return map[string]any{"ok": true, "elementIds": resp.ElementIDs}, nil
}

func (a *api) continueIntervention(_ js.Value, args []js.Value) (any, error) {
	req, err := argObject(args, 0)
	if err != nil {
		return nil, err
	}
	resp, err := a.eng.Continue(context.Background(), processing.ContinueRequest{
		InstanceID: asString(req["instanceId"]),
	})
	if err != nil {
		return nil, err
	}
	return map[string]any{"ok": true, "paused": resp.Paused, "state": interventionStateJSON(&resp.State)}, nil
}

func (a *api) stepInto(_ js.Value, args []js.Value) (any, error) {
	req, err := argObject(args, 0)
	if err != nil {
		return nil, err
	}
	resp, err := a.eng.StepInto(context.Background(), processing.StepIntoRequest{
		InstanceID: asString(req["instanceId"]),
	})
	if err != nil {
		return nil, err
	}
	return map[string]any{"ok": true, "paused": resp.Paused, "state": interventionStateJSON(&resp.State)}, nil
}

func (a *api) stepOver(_ js.Value, args []js.Value) (any, error) {
	req, err := argObject(args, 0)
	if err != nil {
		return nil, err
	}
	resp, err := a.eng.StepOver(context.Background(), processing.StepOverRequest{
		InstanceID: asString(req["instanceId"]),
	})
	if err != nil {
		return nil, err
	}
	return map[string]any{"ok": true, "paused": resp.Paused, "state": interventionStateJSON(&resp.State)}, nil
}

func (a *api) getInterventionState(_ js.Value, args []js.Value) (any, error) {
	instanceID := ""
	if len(args) > 0 && !args[0].IsUndefined() && !args[0].IsNull() {
		if args[0].Type() == js.TypeString {
			instanceID = args[0].String()
		} else {
			req, err := argObject(args, 0)
			if err != nil {
				return nil, err
			}
			instanceID = asString(req["instanceId"])
		}
	}
	st, err := a.eng.GetInterventionState(context.Background(), instanceID)
	if err != nil {
		return nil, err
	}
	return interventionStateJSON(st), nil
}

func (a *api) setVariables(_ js.Value, args []js.Value) (any, error) {
	req, err := argObject(args, 0)
	if err != nil {
		return nil, err
	}
	vars, err := asAnyMap(req["variables"])
	if err != nil {
		return nil, fmt.Errorf("INVALID_ARGUMENT: variables: %w", err)
	}
	resp, err := a.eng.SetVariables(context.Background(), processing.SetVariablesRequest{
		InstanceID: asString(req["instanceId"]),
		Variables:  vars,
	})
	if err != nil {
		return nil, err
	}
	return map[string]any{"ok": true, "state": interventionStateJSON(&resp.State)}, nil
}

func interventionStateJSON(st *processing.InterventionState) map[string]any {
	if st == nil {
		return nil
	}
	out := map[string]any{
		"enabled":         st.Enabled,
		"focusInstanceId": st.FocusInstanceID,
		"policy":          string(st.Policy),
		"breakpoints":     st.Breakpoints,
		"paused":          st.Paused,
		"pauseReason":     string(st.PauseReason),
		"pauseElementId":  st.PauseElementID,
		"pauseTokenId":    st.PauseTokenID,
	}
	if st.Pending != nil {
		out["pending"] = map[string]any{
			"kind":            string(st.Pending.Kind),
			"fromElementId":   st.Pending.FromElementID,
			"tokenId":         st.Pending.TokenID,
			"takenFlowIds":    st.Pending.TakenFlowIDs,
			"nextElementIds":  st.Pending.NextElementIDs,
			"outgoingFlowId":  st.Pending.OutgoingFlowID,
			"enterChildId":    st.Pending.EnterChildID,
			"spawnChildToken": st.Pending.SpawnChildToken,
			"linkCatchIds":    st.Pending.LinkCatchIDs,
		}
	}
	return out
}
