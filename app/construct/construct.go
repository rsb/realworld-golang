package construct

import (
	"net/http"
	"time"

	"github.com/gofiber/contrib/fiberzap"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"

	"github.com/rsb/failure"

	"github.com/rsb/realworld-golang/app"
	"github.com/rsb/realworld-golang/app/conf"
	"github.com/rsb/realworld-golang/foundation/logging"
)

const (
	DefaultHTTPClientTimeout             = 5 * time.Second
	DefaultHTTPClientMaxIde              = 100
	DefaultHTTPClientMaxConnsPerHost     = 100
	DefaultHTTPClientMaxIdleConnsPerHost = 100
)

func NewLogger(appVersion string) (*zap.SugaredLogger, error) {
	l, err := logging.NewLogger(app.ServiceName, appVersion)
	if err != nil {
		return nil, failure.Wrap(err, "logging.NewLogger failed")
	}

	return l, nil
}

func NewAPIMux(d app.Dependencies, c conf.API) *fiber.App {

	app := fiber.New(c.NewFiberConfig())
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(fiberzap.New(
		fiberzap.Config{
			Logger: d.Logger.Desugar(),
		},
	))

	return app
}

func NewHttpClient(config conf.HTTPClient) *http.Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = config.MaxIdleConn
	t.MaxConnsPerHost = config.MaxConnPerHost
	t.MaxIdleConnsPerHost = config.MaxIdleConnPerHost

	return &http.Client{
		Timeout:   config.Timeout,
		Transport: t,
	}
}

func NewDefaultHTTPClient() *http.Client {
	config := conf.HTTPClient{
		Timeout:            DefaultHTTPClientTimeout,
		MaxIdleConn:        DefaultHTTPClientMaxIde,
		MaxConnPerHost:     DefaultHTTPClientMaxConnsPerHost,
		MaxIdleConnPerHost: DefaultHTTPClientMaxIdleConnsPerHost,
	}

	return NewHttpClient(config)
}
