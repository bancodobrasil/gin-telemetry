package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	telemetry "github.com/bancodobrasil/gin-telemetry"
)

func main() {
	router := gin.Default()
	// This handler will match /user/john but will not match /user/ or /user
	router.GET("/user/:name", telemetry.New("example-test"), func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})

	router.Run(":7000")
}
