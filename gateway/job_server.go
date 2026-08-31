package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/sparrow-community/sparrow/processing"
	jobv1 "github.com/sparrow-community/sparrow/protocol/gen/go/job/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// JobServer is a thin JobService adapter over processing.Engine.
type JobServer struct {
	jobv1.UnimplementedJobServiceServer
	engine *processing.Engine
}

func NewJobServer(eng *processing.Engine) *JobServer {
	return &JobServer{engine: eng}
}

// RegisterJobService registers JobService on s.
func RegisterJobService(s grpc.ServiceRegistrar, eng *processing.Engine) {
	jobv1.RegisterJobServiceServer(s, NewJobServer(eng))
}

func (s *JobServer) ActivateJobs(ctx context.Context, req *jobv1.ActivateJobsRequest) (*jobv1.ActivateJobsResponse, error) {
	if s.engine == nil {
		return nil, status.Error(codes.FailedPrecondition, "engine is required")
	}
	jobs, err := s.engine.Activate(ctx, processing.ActivateRequest{
		JobType:      req.GetJobType(),
		MaxJobs:      int(req.GetMaxJobs()),
		Wait:         msDuration(req.GetWaitMs()),
		WorkerID:     req.GetWorkerId(),
		LockDuration: msDuration(req.GetLockDurationMs()),
	})
	if err != nil {
		return nil, statusFromEngine(err)
	}
	out := &jobv1.ActivateJobsResponse{Jobs: make([]*jobv1.Job, 0, len(jobs))}
	for _, job := range jobs {
		out.Jobs = append(out.Jobs, jobToProto(job))
	}
	return out, nil
}

func (s *JobServer) CompleteJob(ctx context.Context, req *jobv1.CompleteJobRequest) (*jobv1.CompleteJobResponse, error) {
	if s.engine == nil {
		return nil, status.Error(codes.FailedPrecondition, "engine is required")
	}
	vars, err := variablesFromJSONMap(req.GetVariables())
	if err != nil {
		return nil, statusFromEngine(err)
	}
	if err := s.engine.Complete(ctx, req.GetProcessInstanceId(), req.GetElementId(), req.GetTokenId(), vars); err != nil {
		return nil, statusFromEngine(err)
	}
	return &jobv1.CompleteJobResponse{}, nil
}

func (s *JobServer) FailJob(ctx context.Context, req *jobv1.FailJobRequest) (*jobv1.FailJobResponse, error) {
	if s.engine == nil {
		return nil, status.Error(codes.FailedPrecondition, "engine is required")
	}
	if err := s.engine.Fail(ctx, req.GetProcessInstanceId(), req.GetElementId(), req.GetTokenId(), req.GetErrorMessage(), req.GetNoRetry()); err != nil {
		return nil, statusFromEngine(err)
	}
	return &jobv1.FailJobResponse{}, nil
}

func (s *JobServer) Heartbeat(ctx context.Context, req *jobv1.HeartbeatRequest) (*jobv1.HeartbeatResponse, error) {
	if s.engine == nil {
		return nil, status.Error(codes.FailedPrecondition, "engine is required")
	}
	if err := s.engine.Heartbeat(ctx, req.GetProcessInstanceId(), req.GetTokenId(), req.GetWorkerId(), msDuration(req.GetLockDurationMs())); err != nil {
		return nil, statusFromEngine(err)
	}
	return &jobv1.HeartbeatResponse{}, nil
}

func jobToProto(job processing.Job) *jobv1.Job {
	return &jobv1.Job{
		JobType:            job.JobType,
		ProcessInstanceId:  job.ProcessInstanceID,
		DeploymentId:       job.DeploymentID,
		ElementId:          job.ElementID,
		TokenId:            job.TokenID,
		Variables:          job.Variables,
		WorkerId:           job.WorkerID,
		LockDeadlineUnixMs: job.LockDeadline.UnixMilli(),
	}
}

func msDuration(ms int64) time.Duration {
	if ms <= 0 {
		return 0
	}
	return time.Duration(ms) * time.Millisecond
}

func variablesFromJSONMap(m map[string]string) (map[string]any, error) {
	if len(m) == 0 {
		return nil, nil
	}
	out := make(map[string]any, len(m))
	for k, raw := range m {
		var v any
		if err := json.Unmarshal([]byte(raw), &v); err != nil {
			return nil, fmt.Errorf("INVALID_ARGUMENT: variable %q: %w", k, err)
		}
		out[k] = v
	}
	return out, nil
}

func statusFromEngine(err error) error {
	if err == nil {
		return nil
	}
	msg := err.Error()
	switch {
	case strings.HasPrefix(msg, "NOT_FOUND"):
		return status.Error(codes.NotFound, msg)
	case strings.HasPrefix(msg, "INVALID_ARGUMENT"):
		return status.Error(codes.InvalidArgument, msg)
	case strings.HasPrefix(msg, "INVALID_STATE"):
		return status.Error(codes.FailedPrecondition, msg)
	case strings.HasPrefix(msg, "INCIDENT_OPEN"), strings.HasPrefix(msg, "NO_INCIDENT"):
		return status.Error(codes.FailedPrecondition, msg)
	default:
		return status.Error(codes.Internal, msg)
	}
}
