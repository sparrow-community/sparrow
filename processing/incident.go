package processing

import (
	"github.com/sparrow-community/sparrow/processing/deploy"
	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

func shouldOpenIncident(threshold int, failCount int32, noRetry bool) bool {
	return noRetry || failCount >= int32(threshold)
}

func jobFailActivityPayload(tok *projection.Token, message string, failCount int32, noRetry bool) *eventv1.ActivityPayload {
	if tok == nil {
		return &eventv1.ActivityPayload{
			ErrorMessage: message,
			JobFailCount: failCount,
			NoRetry:      noRetry,
		}
	}
	return &eventv1.ActivityPayload{
		JobType:      tok.JobType,
		ErrorMessage: message,
		JobFailCount: failCount,
		NoRetry:      noRetry,
	}
}

func emitJobFailChain(dep *deploy.Deployment, inst *projection.Instance, typ eventv1.Element_Type, elementID, tokenID, message string, failCount int32, noRetry bool, emit Emitter) error {
	tok := inst.Tokens[tokenID]
	payload := jobFailActivityPayload(tok, message, failCount, noRetry)
	failed := &eventv1.Element{
		Intent:  eventv1.Element_INTENT_FAILED,
		Type:    typ,
		Id:      elementID,
		TokenId: tokenID,
		Payload: &eventv1.Element_ActivityPayload{ActivityPayload: payload},
	}
	if err := emit(failed); err != nil {
		return err
	}
	threshold := deploy.DefaultIncidentThreshold
	if dep != nil {
		threshold = dep.IncidentThreshold(elementID)
	}
	if !shouldOpenIncident(threshold, failCount, noRetry) {
		return nil
	}
	openPayload := &eventv1.ActivityPayload{
		JobType:      payload.GetJobType(),
		ErrorMessage: message,
		JobFailCount: failCount,
	}
	return emit(&eventv1.Element{
		Intent:  eventv1.Element_INTENT_INCIDENT_OPENED,
		Type:    typ,
		Id:      elementID,
		TokenId: tokenID,
		Payload: &eventv1.Element_ActivityPayload{ActivityPayload: openPayload},
	})
}
