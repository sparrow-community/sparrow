package deploy

import "sort"

// Revision is one numbered snapshot of a BPMN process id.
type Revision struct {
	Version      int32
	DeploymentID string
}

// AssignProcessVersions sets Deployment.Version for each deployment as the
// 1-based rank among deployments that share the same process id, ordered by
// deployment id (UUIDv7 is time-ordered, so this matches Deploy order).
func AssignProcessVersions(deps map[string]*Deployment) map[string][]Revision {
	byProcess := make(map[string][]string)
	for id, dep := range deps {
		if dep == nil {
			continue
		}
		pid := dep.ProcessID()
		byProcess[pid] = append(byProcess[pid], id)
	}
	out := make(map[string][]Revision, len(byProcess))
	for pid, ids := range byProcess {
		sort.Strings(ids)
		revs := make([]Revision, 0, len(ids))
		for i, id := range ids {
			v := int32(i + 1)
			deps[id].Version = v
			revs = append(revs, Revision{Version: v, DeploymentID: id})
		}
		out[pid] = revs
	}
	return out
}

// ResolveRevision returns the deployment id for processID at version.
// version <= 0 selects the latest revision.
func ResolveRevision(revs map[string][]Revision, processID string, version int32) (string, bool) {
	list := revs[processID]
	if len(list) == 0 {
		return "", false
	}
	if version <= 0 {
		return list[len(list)-1].DeploymentID, true
	}
	for _, r := range list {
		if r.Version == version {
			return r.DeploymentID, true
		}
	}
	return "", false
}
