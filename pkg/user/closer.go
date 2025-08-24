package user

import (
	"github.com/cbhcbhcbh/Quantum/pkg/infra"
)

type InfraCloser struct{}

func NewInfraCloser() *InfraCloser {
	return &InfraCloser{}
}

func (closer *InfraCloser) Close() error {
	return infra.RedisClient.Close()
}
