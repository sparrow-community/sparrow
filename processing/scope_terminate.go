package processing

import (
	"github.com/sparrow-community/sparrow/processing/deploy"
	"github.com/sparrow-community/sparrow/processing/projection"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

// scopeTerminateOpts controls how tokens inside a scope are cancelled.
type scopeTerminateOpts struct {
	// IncludeHost also terminates a token whose ElementID equals scopeID.
	IncludeHost bool
	// DropTokens removes tokens from the projection after TERMINATED so
	// ApplyEvent cannot revive them as active (non-catch TERMINATED keeps tokens).
	DropTokens bool
}

// terminateScopeTokens emits TERMINATING/TERMINATED for tokens in scopeID.
// When scopeID is the root process, every token is terminated.
func terminateScopeTokens(
	dep *deploy.Deployment,
	inst *projection.Instance,
	scopeID string,
	emit Emitter,
	opts scopeTerminateOpts,
) error {
	processID := dep.ProcessID()
	ids := make([]string, 0, len(inst.Tokens))
	for tid, tok := range inst.Tokens {
		if tok == nil {
			continue
		}
		if scopeID != processID {
			if !tokenInOrIsScope(dep, tok, scopeID) {
				continue
			}
			if tok.ElementID == scopeID && !opts.IncludeHost {
				continue
			}
		}
		ids = append(ids, tid)
	}
	for _, tid := range ids {
		tok := inst.Tokens[tid]
		if tok == nil {
			continue
		}
		tokType, err := dep.TypeOf(tok.ElementID)
		if err != nil {
			tokType = eventv1.Element_TYPE_UNSPECIFIED
		}
		for _, intent := range []eventv1.Element_Intent{
			eventv1.Element_INTENT_TERMINATING,
			eventv1.Element_INTENT_TERMINATED,
		} {
			if err := emit(&eventv1.Element{
				Intent:  intent,
				Type:    tokType,
				Id:      tok.ElementID,
				TokenId: tid,
			}); err != nil {
				return err
			}
		}
		if opts.DropTokens {
			delete(inst.Tokens, tid)
		}
	}
	return nil
}

// terminateEmbeddedScope records SubProcess TERMINATING/TERMINATED for an embedded scope.
// Empty tokenID means audit-only: the EventLog still shows the cancel, but ApplyEvent
// does not revive a projection token. Non-empty tokenID places that token on the
// SubProcess so the caller can take outgoing (e.g. error / timer / signal boundary).
func terminateEmbeddedScope(scopeID, tokenID string, emit Emitter) error {
	for _, intent := range []eventv1.Element_Intent{
		eventv1.Element_INTENT_TERMINATING,
		eventv1.Element_INTENT_TERMINATED,
	} {
		if err := emit(&eventv1.Element{
			Intent:  intent,
			Type:    eventv1.Element_TYPE_SUB_PROCESS,
			Id:      scopeID,
			TokenId: tokenID,
		}); err != nil {
			return err
		}
	}
	return nil
}

// tokenInOrIsScope reports whether tok sits on scopeID itself or inside it.
func tokenInOrIsScope(dep *deploy.Deployment, tok *projection.Token, scopeID string) bool {
	if tok.ElementID == scopeID {
		return true
	}
	tokScope, _ := dep.ScopeOf(tok.ElementID)
	return tokScope == scopeID || isInScope(dep, tokScope, scopeID)
}
