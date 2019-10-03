package consul

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
)

const (
	defaultInterval       = "3s"
	defaultTimeOut        = "3s"
	defaultDeregisterTime = "360s"
)

//Service ...
type Service struct {
	consulClient *consulapi.Client
	id           string
	name         string
	hostname     string
	port         int
	tags         []string
}

//NewService ...
func NewService(id, name, hostname string, port int, tags []string) (*Service, error) {
	consulClient, err := consulapi.NewClient(consulapi.DefaultConfig())
	if err != nil {
		return nil, err
	}

	consulService := &Service{
		consulClient: consulClient,
		id:           id,
		name:         name,
		hostname:     hostname,
		port:         port,
		tags:         tags,
	}

	return consulService, nil
}

//Register ...
func (s *Service) Register() error {
	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = s.id
	registration.Name = s.name
	registration.Kind = consulapi.ServiceKindTypical
	registration.Tags = s.tags
	registration.Port = s.port
	registration.Address = s.hostname
	registration.Check = &consulapi.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%v:%d", s.hostname, s.port),
		Interval:                       defaultInterval,
		DeregisterCriticalServiceAfter: defaultDeregisterTime,
	}

	s.consulClient.Agent().ServiceDeregister(s.id)
	return s.consulClient.Agent().ServiceRegister(registration)
}

//Deregister ...
func (s *Service) Deregister() error {
	return s.consulClient.Agent().ServiceDeregister(s.id)
}

//GetKV ...
func (s *Service) GetKV(kvname string) (string, error) {
	kvPair, _, err := s.consulClient.KV().Get(kvname, nil)
	if err != nil {
		return "", nil
	}

	return string(kvPair.Value), nil
}
