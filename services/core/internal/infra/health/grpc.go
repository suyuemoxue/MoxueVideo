package health

import (
	"context"
	"errors"

	"moxuevideo/core/internal/infra/grpcchat"
)

type GRPCChecker struct {
	client *grpcchat.Client
}

func NewGRPCChecker(client *grpcchat.Client) *GRPCChecker {
	return &GRPCChecker{client: client}
}

func (c *GRPCChecker) Check(ctx context.Context) error {
	if c == nil || c.client == nil {
		return errors.New("grpc unavailable")
	}
	return c.client.Ping(ctx)
}
