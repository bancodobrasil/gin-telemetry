# gin-telemetry

Gin Telemetry is a middleware to exporting traces based on the open-telemetry project

## Install
To install execute the command below:
```
go get -u github.com/bancodobrasil/gin-telemetry
```
## Environemnt variables
Below are the environment variables that affect the behavior of the middleware: 
| Variable | Description | Default |
| -------- | ------------| ------- |
| TELEMETRY_EXPORTER_URL | Endpoint URL for the collector that spans are sent to  | http://localhost:14268 |
| TELEMETRY_HTTPCLIENT_TLS | Custom http client for passthrough TLS enabled | true |

## Register Tracing Middleware

You must register the tracing middleware to enable export traces

Example for all routes:

```go
import (
  "github.com/gin-gonic/gin"
	telemetry "github.com/bancodobrasil/gin-telemetry"
)

func main() {
  // ... extra code
  router := gin.Default()
  router.Use(telemtry.Middleware("service_name"))
  // ... extra code
}
```

Example for specific route:

```go
import (
  "github.com/gin-gonic/gin"
	telemetry "github.com/bancodobrasil/gin-telemetry"
)

func main() {
  // ... extra code
  
  router := gin.Default()
  router.GET("/user/:name", telemetry.Middleware(serviceName), handlerImpl)
  
  // ... extra code
}

func handlerImple(c *gin.Context) {
  // handler implementation
}
```
## Passthrough - Propagate context

Context propagation example through custom http client

```go
import (
  "net/http"

  "github.com/gin-gonic/gin"
	telemetry "github.com/bancodobrasil/gin-telemetry"
)

func main() {
  // ... extra code
  
  router := gin.Default()
  router.GET("/user/:name", telemetry.Middleware(serviceName), handlerImpl)
  
  // ... extra code
}

func handlerImple(c *gin.Context) {
  // recover context
  ctx := c.Request.Context()
  // create a request with context
  req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/user/%s", externalServiceURL, name), nil)

  // ... extra code

  // do request with httpclient telemetry ** MUST **
  res, err := telemetry.HTTPClient.Do(req)
  
  // ... extra code
}
```
## Example

in the `example` folder we have a docker-compose for testing traceability between 3 services and a active Jaeger as an exporter. 

To start docker-compose executes:

```bash
~example$ docker-compose up --build -d
```

The services are published with the address below:

```url
http://localhost:7001/user/:name
http://localhost:7002/user/:name
http://localhost:7003/user/:name
```
The jaeger is published with the address below:
```url
http://localhost:16686
```

Open the browser and execute request according to addresses above and check jaeger panel