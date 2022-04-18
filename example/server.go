package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"

	telemetry "github.com/bancodobrasil/gin-telemetry"
)

func main() {

	externalServiceURL := os.Getenv("EXTERNAL_SERVICE_URL")
	serviceName := os.Getenv("SERVICE_NAME")
	serverPort := os.Getenv("PORT")
	if serverPort == "" {
		serverPort = "7000"
	}

	router := gin.Default()

	router.GET("/user/:name", telemetry.Middleware(serviceName), func(c *gin.Context) {
		name := c.Param("name")
		if externalServiceURL != "" {
			ctx := c.Request.Context()
			req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/user/%s", externalServiceURL, name), nil)

			if err != nil {
				log.Fatal(err)
			}

			res, err := telemetry.HTTPClient.Do(req)

			if err != nil {
				log.Fatal(err)
			}

			defer res.Body.Close()

			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Fatal(err)
			}
			name = string(data)
			c.String(http.StatusOK, "Response from %s: '%s'", externalServiceURL, name)
		} else {
			c.String(http.StatusOK, "Hello %s", name)
		}
	})

	router.Run(fmt.Sprintf(":%s", serverPort))
}
