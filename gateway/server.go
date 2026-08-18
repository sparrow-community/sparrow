package gateway

import (
	"github.com/sparrow-community/sparrow/processing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// NewServer returns a gRPC server with EngineService, JobService, and Health.
func NewServer(eng *processing.Engine, opts ...grpc.ServerOption) *grpc.Server {
	s := grpc.NewServer(opts...)
	RegisterEngineService(s, eng)
	RegisterJobService(s, eng)
	hs := health.NewServer()
	hs.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(s, hs)
	return s
}
