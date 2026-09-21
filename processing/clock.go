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

package processing

import "time"

// SetNow replaces the engine clock used by FireDue and timer arming.
// Pass nil to restore time.Now. Hosts (WASM, tests) inject a shared clock
// so due checks match the scheduler that calls FireDue.
func (e *Engine) SetNow(fn func() time.Time) {
	if e == nil {
		return
	}
	e.nowFn = fn
}
