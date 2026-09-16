// Copyright 2025 The Sparrow community and contributors
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

package processing

import (
	"context"
	"fmt"

	"github.com/sparrow-community/sparrow/processing/deploy"
	"github.com/sparrow-community/sparrow/processing/handlers"
	"github.com/sparrow-community/sparrow/processing/projection"
)

// runAdHocEnable enables inner activities of an Ad-Hoc SubProcess. The host
// token stays parked on the scope; each enabled activity gets its own token.
func (x *Executor) runAdHocEnable(
	ctx context.Context,
	dep *deploy.Deployment,
	inst *projection.Instance,
	hostTokenID string,
	enable *handlers.AdHocEnable,
	emit Emitter,
) ([]handlers.Publication, error) {
	var pubs []handlers.Publication
	if enable == nil {
		return pubs, nil
	}
	for _, activityID := range enable.ActivityIDs {
		more, err := x.enableAdHocActivity(ctx, dep, inst, hostTokenID, activityID, emit)
		pubs = append(pubs, more...)
		if err != nil {
			return pubs, err
		}
	}
	return pubs, nil
}

func (x *Executor) enableAdHocActivity(
	ctx context.Context,
	dep *deploy.Deployment,
	inst *projection.Instance,
	hostTokenID, activityID string,
	emit Emitter,
) ([]handlers.Publication, error) {
	tokenID, err := spawnScopeChildToken(inst, hostTokenID)
	if err != nil {
		return nil, err
	}
	return x.Enter(ctx, dep, inst, tokenID, activityID, emit)
}

// advanceAdHoc reacts to one inner activity completing: re-evaluate the
// completion condition, then cancel remaining work, enable the next activity,
// or complete the scope when nothing is enabled anymore.
func (x *Executor) advanceAdHoc(
	ctx context.Context,
	dep *deploy.Deployment,
	inst *projection.Instance,
	adHocID, innerTokenID, innerElementID string,
	emit Emitter,
) ([]handlers.Publication, error) {
	spec, ok := dep.AdHocSubProcessSpecOf(adHocID)
	if !ok {
		return nil, fmt.Errorf("NOT_FOUND: adHocSubProcess %q", adHocID)
	}
	hostTokenID, ok := adHocHostToken(inst, adHocID, innerTokenID)
	if !ok {
		return nil, fmt.Errorf("NOT_FOUND: adHocSubProcess host token for %q", adHocID)
	}
	delete(inst.Tokens, innerTokenID)

	met, err := dep.EvalAdHocCompletion(adHocID, inst.Variables)
	if err != nil {
		return nil, err
	}
	remaining := adHocEnabledTokens(dep, inst, adHocID, hostTokenID)
	if met {
		if len(remaining) > 0 && !spec.CancelRemaining {
			// Keep running instances; the last one completes the scope.
			return nil, nil
		}
		return x.completeAdHocScope(ctx, dep, inst, adHocID, hostTokenID, len(remaining) > 0, emit)
	}
	if spec.Sequential {
		if nextID, ok := spec.NextSequentialActivity(innerElementID); ok {
			return x.enableAdHocActivity(ctx, dep, inst, hostTokenID, nextID, emit)
		}
	}
	if len(remaining) > 0 {
		return nil, nil
	}
	// No enabled inner activity remains: finish the scope instead of deadlocking.
	return x.completeAdHocScope(ctx, dep, inst, adHocID, hostTokenID, false, emit)
}

// completeAdHocScope optionally cancels inner activities still enabled, then
// completes the scope on its host token.
func (x *Executor) completeAdHocScope(
	ctx context.Context,
	dep *deploy.Deployment,
	inst *projection.Instance,
	adHocID, hostTokenID string,
	cancelRemaining bool,
	emit Emitter,
) ([]handlers.Publication, error) {
	if cancelRemaining {
		if err := terminateScopeTokens(dep, inst, adHocID, emit, scopeTerminateOpts{
			DropTokens:   true,
			KeepTokenIDs: map[string]bool{hostTokenID: true},
		}); err != nil {
			return nil, err
		}
	}
	return x.Complete(ctx, dep, inst, hostTokenID, adHocID, nil, emit)
}

// adHocHostToken returns the token parked on the Ad-Hoc SubProcess scope.
func adHocHostToken(inst *projection.Instance, adHocID, innerTokenID string) (string, bool) {
	if tok := inst.Tokens[innerTokenID]; tok != nil && tok.ScopeHostTokenID != "" {
		if host := inst.Tokens[tok.ScopeHostTokenID]; host != nil && host.ElementID == adHocID {
			return tok.ScopeHostTokenID, true
		}
	}
	for tid, tok := range inst.Tokens {
		if tok != nil && tok.ElementID == adHocID {
			return tid, true
		}
	}
	return "", false
}

// adHocEnabledTokens lists tokens of inner activities still enabled in the scope.
// Membership comes from the element scope, so Recover (which does not carry the
// host link) sees the same set.
func adHocEnabledTokens(dep *deploy.Deployment, inst *projection.Instance, adHocID, hostTokenID string) []string {
	var out []string
	for tid, tok := range inst.Tokens {
		if tok == nil || tid == hostTokenID {
			continue
		}
		if scope, ok := dep.ScopeOf(tok.ElementID); !ok || scope != adHocID {
			continue
		}
		if tok.Status == projection.TokenWaiting || tok.Status == projection.TokenBlocked || tok.Status == projection.TokenActive {
			out = append(out, tid)
		}
	}
	return out
}
