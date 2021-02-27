package service

import (
	"github.com/1024casts/snake/pkg/conf"
)

var (
	// Svc global service var
	Svc *Service
)

// Service struct
type Service struct {
	c *conf.Config
}

// New init service
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c: c,
	}
	return s
}

// Ping service
func (s *Service) Ping() error {
	return nil
}

// Close service
func (s *Service) Close() {
}
