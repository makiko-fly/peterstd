package http

import (
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/labstack/echo"
)

func ApiCountMiddleware(rp *redis.Pool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			accessTime := time.Now()

			path := c.Path()
			if !strings.Contains(path, "/api") {
				return err
			}

			dayStr := accessTime.Format("flash:api:access:2006-01-02")
			hourStr := accessTime.Format("flash:api:access:2006-01-02:15")

			conn := rp.Get()
			defer conn.Close()

			conn.Send("MULTI")
			conn.Send("ZINCRBY", dayStr, 1, path)
			conn.Send("ZINCRBY", hourStr, 1, path)
			_, err = conn.Do("EXEC")
			if err != nil {
				c.Logger().Error(err)
				return err
			}

			return next(c)

		}
	}
}
