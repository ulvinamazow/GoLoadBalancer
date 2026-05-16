package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ulvinamazow/load_balancer/balancer"
	"github.com/ulvinamazow/load_balancer/proxy"
)

func main() {
	servers := []*balancer.Server{
		{URL: "http://localhost:9090", Weight: 5},
		{URL: "http://localhost:9091", Weight: 4},
	}

	lb := balancer.NewWeightedRoundRobin(servers)

	r := gin.Default()

	r.Any("/*path", proxy.NewHandler(lb))

	if err := r.Run(":8080"); err != nil {
		panic("Server can not started")
	}
}
