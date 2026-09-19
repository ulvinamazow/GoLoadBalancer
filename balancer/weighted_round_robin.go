package balancer

import (
	"sync"
)

type Server struct {
	URL           string
	Weight        int
	CurrentWeight int
	Health        bool
}

type WeightedRoundRobin struct {
	servers     []*Server
	totalWeight int
	mu          sync.Mutex
}

func NewWeightedRoundRobin(servers []*Server) *WeightedRoundRobin {
	total := 0

	for _, s := range servers {
		total += s.Weight
	}

	return &WeightedRoundRobin{
		servers:     servers,
		totalWeight: total,
	}
}

func (wwr *WeightedRoundRobin) NextServer() *Server {
	wwr.mu.Lock()
	defer wwr.mu.Unlock()

	if wwr.totalWeight == 0 || len(wwr.servers) == 0 {
		return nil
	}

	var bestServer *Server
	maxCurrentWeight := -1

	for _, s := range wwr.servers {
		//
		s.CurrentWeight += s.Weight
		if s.CurrentWeight > maxCurrentWeight {
			maxCurrentWeight = s.CurrentWeight
			bestServer = s
		}
	}

	if bestServer != nil {
		bestServer.CurrentWeight -= wwr.totalWeight
	}

	if bestServer.IsHealth() == true {
		return bestServer
	} else {
		return nil
	}
}

func (server *Server) SetHealthy(status bool) {
	server.Health = status
}

func (server *Server) IsHealth() bool {
	return server.Health
}
