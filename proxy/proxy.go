package proxy

import (
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulvinamazow/load_balancer/balancer"
)

func NewHandler(lb *balancer.WeightedRoundRobin) gin.HandlerFunc {

	return func(c *gin.Context) {
		server := lb.NextServer()

		if server == nil {
			c.String(http.StatusServiceUnavailable, "No backend available")
			return
		}

		targetPath := c.Request.URL.Path

		backendURL := server.URL + targetPath

		request, err := http.NewRequestWithContext(c, c.Request.Method, backendURL, c.Request.Body)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to create backend request: "+err.Error())
			return
		}

		request.Header = c.Request.Header.Clone()

		if clientIP := c.ClientIP(); clientIP != "" {
			request.Header.Set("X-Forwarded-for", clientIP)
		}

		request.Header.Set("X-Forwarded-Proto", "http")
		request.Header.Set("X-Forwarded-Host", c.Request.Host)

		cliet := &http.Client{
			Timeout: 10 * time.Second,
		}

		response, err := cliet.Do(request)

		if err != nil {
			c.String(http.StatusBadGateway, "Bckend unreachable: "+err.Error())
			return
		}
		defer response.Body.Close()

		for key, vals := range request.Header {
			for _, v := range vals {
				c.Header(key, v)
			}
		}

		c.Status(response.StatusCode)

		if _, err := io.Copy(c.Writer, response.Body); err != nil {
		}
	}
}
