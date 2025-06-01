package dependency

import (
	"canary/internal/types"
)

type Service interface {
	Resolve(req *types.DependencyCheckRequest) (string, error)
}

type service struct{}

func NewService() Service {
	return &service{}
}

func (s *service) Resolve(req *types.DependencyCheckRequest) (string, error) {
	return req.FileContent, nil
}
