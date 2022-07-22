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
| TELEMETRY_EXPORTER_JAEGER_AGENT_HOST | Agent Host for the collector that spans are sent to  | localhost |
| TELEMETRY_EXPORTER_JAEGER_AGENT_PORT | Agent Port URL for the collector that spans are sent to  | 6831 |
| TELEMETRY_HTTPCLIENT_TLS | Custom http client for passthrough TLS enabled | true |
| TELEMETRY_DISABLED | Disable middleware to export span for collectors | false |
| TELEMETRY_ENVIRONMENT | Environment to export span for collectors | test |

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

func handlerImpl(c *gin.Context) {
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
## Logrus Logging

This instrumentation records logrus log messages as events on the existing span that is passed via a context.Context.

**It does not record anything if a context does not contain a span.**

Example:
```go
import (
    "github.com/uptrace/opentelemetry-go-extra/otellogrus"
    "github.com/sirupsen/logrus"
)
// ...extra code
// Use ctx to pass the active span.
logrus.WithContext(ctx).
	WithError(errors.New("hello world")).
	WithField("foo", "bar").
	Error("something failed")
  // ...extra code
```

## Docker Compose Example

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

Open the browser and execute request according to addresses above and check jaeger panel. Thx!