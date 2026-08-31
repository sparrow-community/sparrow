package gateway

import (
	"context"
	"slices"

	"github.com/sparrow-community/sparrow/processing"
	"github.com/sparrow-community/sparrow/processing/projection"
	enginev1 "github.com/sparrow-community/sparrow/protocol/gen/go/engine/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// EngineServer is a thin EngineService adapter over processing.Engine.
type EngineServer struct {
	enginev1.UnimplementedEngineServiceServer
	engine *processing.Engine
}

func NewEngineServer(eng *processing.Engine) *EngineServer {
	return &EngineServer{engine: eng}
}

// RegisterEngineService registers EngineService on s.
func RegisterEngineService(s grpc.ServiceRegistrar, eng *processing.Engine) {
	enginev1.RegisterEngineServiceServer(s, NewEngineServer(eng))
}

func (s *EngineServer) Deploy(ctx context.Context, req *enginev1.DeployRequest) (*enginev1.DeployResponse, error) {
	if s.engine == nil {
		return nil, status.Error(codes.FailedPrecondition, "engine is required")
	}
	id, err := s.engine.Deploy(ctx, req.GetBpmnXml())
	if err != nil {
		return nil, statusFromEngine(err)
	}
	out := &enginev1.DeployResponse{DeploymentId: id}
	if dep, ok := s.engine.GetDeployment(id); ok {
		out.ProcessId = dep.ProcessID()
		out.ProcessVersion = dep.Version
	}
	return out, nil
}

func (s *EngineServer) CreateInstance(ctx context.Context, req *enginev1.CreateInstanceRequest) (*enginev1.CreateInstanceResponse, error) {
	if s.engine == nil {
		return nil, status.Error(codes.FailedPrecondition, "engine is required")
	}
	vars, err := variablesFromJSONMap(req.GetVariables())
	if err != nil {
		return nil, statusFromEngine(err)
	}
	id, err := s.engine.CreateInstanceRequest(ctx, processing.CreateInstanceRequest{
		DeploymentID:   req.GetDeploymentId(),
		ProcessID:      req.GetProcessId(),
		ProcessVersion: req.GetProcessVersion(),
		Variables:      vars,
	})
	if err != nil {
		return nil, statusFromEngine(err)
	}
	return &enginev1.CreateInstanceResponse{ProcessInstanceId: id}, nil
}

func (s *EngineServer) Complete(ctx context.Context, req *enginev1.CompleteRequest) (*enginev1.CompleteResponse, error) {
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
	return &enginev1.CompleteResponse{}, nil
}

func (s *EngineServer) PublishMessage(ctx context.Context, req *enginev1.PublishMessageRequest) (*enginev1.PublishMessageResponse, error) {
	if s.engine == nil {
		return nil, status.Error(codes.FailedPrecondition, "engine is required")
	}
	vars, err := variablesFromJSONMap(req.GetVariables())
	if err != nil {
		return nil, statusFromEngine(err)
	}
	keys, err := variablesFromJSONMap(req.GetCorrelationKeys())
	if err != nil {
		return nil, statusFromEngine(err)
	}
	n, err := s.engine.PublishMessage(ctx, processing.PublishMessageRequest{
		Name:              req.GetName(),
		ProcessInstanceID: req.GetProcessInstanceId(),
		CorrelationKeys:   keys,
		Variables:         vars,
	})
	if err != nil {
		return nil, statusFromEngine(err)
	}
	out := &enginev1.PublishMessageResponse{Delivered: int32(n)}
	if n == 0 {
		out.Buffered = 1
	}
	return out, nil
}

func (s *EngineServer) PublishSignal(ctx context.Context, req *enginev1.PublishSignalRequest) (*enginev1.PublishSignalResponse, error) {
	if s.engine == nil {
		return nil, status.Error(codes.FailedPrecondition, "engine is required")
	}
	vars, err := variablesFromJSONMap(req.GetVariables())
	if err != nil {
		return nil, statusFromEngine(err)
	}
	n, err := s.engine.PublishSignal(ctx, processing.PublishSignalRequest{
		Name:              req.GetName(),
		ProcessInstanceID: req.GetProcessInstanceId(),
		Variables:         vars,
	})
	if err != nil {
		return nil, statusFromEngine(err)
	}
	return &enginev1.PublishSignalResponse{Delivered: int32(n)}, nil
}

func (s *EngineServer) ThrowError(ctx context.Context, req *enginev1.ThrowErrorRequest) (*enginev1.ThrowErrorResponse, error) {
	if s.engine == nil {
		return nil, status.Error(codes.FailedPrecondition, "engine is required")
	}
	if err := s.engine.ThrowError(ctx, req.GetProcessInstanceId(), req.GetElementId(), req.GetTokenId(), req.GetErrorCode()); err != nil {
		return nil, statusFromEngine(err)
	}
	return &enginev1.ThrowErrorResponse{}, nil
}

func (s *EngineServer) ResolveIncident(ctx context.Context, req *enginev1.ResolveIncidentRequest) (*enginev1.ResolveIncidentResponse, error) {
	if s.engine == nil {
		return nil, status.Error(codes.FailedPrecondition, "engine is required")
	}
	if err := s.engine.ResolveIncident(ctx, req.GetProcessInstanceId(), req.GetElementId(), req.GetTokenId()); err != nil {
		return nil, statusFromEngine(err)
	}
	return &enginev1.ResolveIncidentResponse{}, nil
}

func (s *EngineServer) GetInstance(_ context.Context, req *enginev1.GetInstanceRequest) (*enginev1.GetInstanceResponse, error) {
	if s.engine == nil {
		return nil, status.Error(codes.FailedPrecondition, "engine is required")
	}
	inst, ok := s.engine.GetInstance(req.GetProcessInstanceId())
	if !ok {
		return nil, status.Errorf(codes.NotFound, "NOT_FOUND: instance %q", req.GetProcessInstanceId())
	}
	return &enginev1.GetInstanceResponse{Instance: instanceToProto(inst)}, nil
}

func (s *EngineServer) ListEvents(ctx context.Context, req *enginev1.ListEventsRequest) (*enginev1.ListEventsResponse, error) {
	if s.engine == nil {
		return nil, status.Error(codes.FailedPrecondition, "engine is required")
	}
	events, err := s.engine.ListEvents(ctx, req.GetProcessInstanceId())
	if err != nil {
		return nil, statusFromEngine(err)
	}
	return &enginev1.ListEventsResponse{Events: events}, nil
}

func instanceToProto(inst *projection.Instance) *enginev1.Instance {
	if inst == nil {
		return nil
	}
	out := &enginev1.Instance{
		Id:                      inst.ID,
		DeploymentId:            inst.DeploymentID,
		ProcessVersion:          inst.Version,
		Status:                  string(inst.Status),
		Variables:               inst.Variables,
		Tokens:                  make([]*enginev1.Token, 0, len(inst.Tokens)),
		ParentProcessInstanceId: inst.ParentProcessInstanceID,
		ParentElementId:         inst.ParentElementID,
		ProcessId:               inst.ProcessID,
	}
	ids := make([]string, 0, len(inst.Tokens))
	for id := range inst.Tokens {
		ids = append(ids, id)
	}
	slices.Sort(ids)
	for _, id := range ids {
		tok := inst.Tokens[id]
		if tok == nil {
			continue
		}
		out.Tokens = append(out.Tokens, &enginev1.Token{
			Id:                      tok.ID,
			ElementId:               tok.ElementID,
			Status:                  string(tok.Status),
			JobType:                 tok.JobType,
			MessageName:             tok.MessageName,
			DueUnixMs:               tok.DueUnixMs,
			BoundaryId:              tok.BoundaryID,
			SignalName:              tok.SignalName,
			CalledProcessInstanceId: tok.CalledProcessInstanceID,
			LoopInstanceIndex:       tok.LoopInstanceIndex,
			IncidentErrorMessage:    tok.IncidentErrorMessage,
			JobFailCount:            tok.JobFailCount,
		})
	}
	return out
}
