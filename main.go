package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulvinamazow/load_balancer/balancer"
	"github.com/ulvinamazow/load_balancer/health"
	"github.com/ulvinamazow/load_balancer/proxy"
)

func main() {
	servers := []*balancer.Server{
		{URL: "http://localhost:9090", Weight: 5},
		{URL: "http://localhost:9091", Weight: 4},
	}

	go func() {
		for {
			time.Sleep(10 * time.Second)

			for _, backend := range servers {
				code, err := health.Checker(http.MethodGet, backend.URL+"/health")

				if err != nil || code != http.StatusOK {
					backend.SetHealthy(true)
				} else {
					backend.SetHealthy(false)
				}
			}
		}
	}()

	lb := balancer.NewWeightedRoundRobin(servers)

	r := gin.Default()

	r.Any("/*path", proxy.NewHandler(lb))

	if err := r.Run(":8080"); err != nil {
		panic("Server can not started")
	}
}
