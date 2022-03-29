package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	telemetry "github.com/bancodobrasil/gin-telemetry"
)

func main() {

	externalServiceUrl := os.Getenv("EXTERNAL_SERVICE_URL")
	serviceName := os.Getenv("SERVICE_NAME")

	router := gin.Default()

	router.GET("/user/:name", telemetry.New(serviceName), func(c *gin.Context) {
		name := c.Param("name")
		if externalServiceUrl != "" {
			res, err := http.Get(fmt.Sprintf("%s/user/%s", externalServiceUrl, name))
			if err != nil {
				log.Fatal(err)
			}
			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Fatal(err)
			}
			name = string(data)
			res.Body.Close()
		}

		c.String(http.StatusOK, "Hello %s", name)
	})

	router.Run(":7000")
}
