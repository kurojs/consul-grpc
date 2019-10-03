package consul

import (
	"fmt"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"google.golang.org/grpc/naming"
)

//Resolver ...
type Resolver struct {
	serviceName  string
	tag          string
	consulClient *consulapi.Client
	done         chan interface{}
	updateCh     chan []*naming.Update
}

//NewResolver ...
func NewResolver(servicename, tag string) (*Resolver, error) {
	client, err := consulapi.NewClient(consulapi.DefaultConfig())
	if err != nil {
		return nil, err
	}

	resolver := &Resolver{
		serviceName:  servicename,
		tag:          tag,
		consulClient: client,
		done:         make(chan interface{}),
		updateCh:     make(chan []*naming.Update, 1),
	}

	instances, lastIndex, err := resolver.getInstances(0, true)
	if err != nil {
		// log
	}

	updates := resolver.updateInstances(nil, instances)
	if len(updates) > 0 {
		resolver.updateCh <- updates
	}

	go resolver.worker(instances, lastIndex)

	return resolver, nil
}

// Resolve ...
func (r *Resolver) Resolve(target string) (naming.Watcher, error) {
	return r, nil
}

// Next ...
func (r *Resolver) Next() ([]*naming.Update, error) {
	return <-r.updateCh, nil
}

// Close ...
func (r *Resolver) Close() {
	select {
	case <-r.done:
	default:
		close(r.done)
		close(r.updateCh)
	}
}

func (r *Resolver) getInstances(lastIndex uint64, passOnly bool) ([]string, uint64, error) {
	services, metadata, err := r.consulClient.Health().Service(
		r.serviceName,
		r.tag, passOnly,
		&consulapi.QueryOptions{
			WaitIndex: lastIndex,
		},
	)
	if err != nil {
		return nil, 0, err
	}

	instances := []string{}

	for _, service := range services {
		addr := service.Service.Address
		if len(addr) == 0 {
			addr = service.Node.Address
		}

		address := fmt.Sprintf("%s:%d", addr, service.Service.Port)
		instances = append(instances, address)
	}

	return instances, metadata.LastIndex, nil
}

// updateInstance mix 2 array and truncat duplicate elements
func (r *Resolver) updateInstances(oldInstances, newInstances []string) []*naming.Update {
	oldAddr := make(map[string]bool, len(oldInstances))
	for _, instance := range oldInstances {
		oldAddr[instance] = true
	}

	newAddr := make(map[string]bool, len(newInstances))
	for _, instance := range newInstances {
		newAddr[instance] = true
	}

	var updates []*naming.Update
	for addr := range newAddr {
		if _, ok := oldAddr[addr]; !ok {
			updates = append(updates, &naming.Update{
				Op:   naming.Add,
				Addr: addr,
			})
		}
	}

	for addr := range oldAddr {
		if _, ok := newAddr[addr]; !ok {
			updates = append(updates, &naming.Update{
				Op:   naming.Delete,
				Addr: addr,
			})
		}
	}

	return updates
}

// worker is background process, it query to consul and detect config change
func (r *Resolver) worker(instances []string, lastIndex uint64) {
	var err error
	var newInstances []string

	for {
		time.Sleep(5 * time.Second)
		select {
		case <-r.done:
			return
		default:
			newInstances, lastIndex, err = r.getInstances(lastIndex, true)
			if err != nil {
				// log
				continue
			}

			updatedInstances := r.updateInstances(instances, newInstances)
			if len(updatedInstances) > 0 {
				r.updateCh <- updatedInstances
			}
			instances = newInstances
		}
	}
}
