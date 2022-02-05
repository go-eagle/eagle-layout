package service

// Svc global var
var Svc Service

// Service define all service
type Service interface {
	Greeter() IGreeterService
}

// service struct
type service struct {
}

// New init service
func New() Service {
	return &service{}
}

func (s *service) Greeter() IGreeterService {
	return newGreeterService(s)
}
