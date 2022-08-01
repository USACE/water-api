package middleware

import (
	"log"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func PgFeatureservProxy(urlStr string) echo.MiddlewareFunc {

	url, err := url.Parse(urlStr)
	if err != nil {
		log.Fatal(err)
	}

	return middleware.Proxy(middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{
		{
			URL: url,
		},
	}))
}
