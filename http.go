package peterstd

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type IHandler interface {
	Mount(e *echo.Group) error
}

type HTTPConfig struct {
	Address string
}

// ---------------------------------------------------------------------------------------------------------------------
// HTTP Server

type HTTPServer struct {
	e        *echo.Echo
	Address  string
	apiGroup *echo.Group
	rp       *redis.Pool
}

func NewHTTPServer(rp *redis.Pool, ei *echo.Echo, config HTTPConfig, handlers ...IHandler) (*HTTPServer, error) {
	if ei == nil {
		ei = echo.New()
	}

	ei.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAcceptEncoding, "X-Ivanka-Token"},
	}))

	ei.Use(ApiMetricsToRedisMiddle(rp))

	s := &HTTPServer{
		Address:  config.Address,
		apiGroup: ei.Group("api"),
		e:        ei,
	}

	if err := s.Register(handlers...); err != nil {
		return nil, err
	}

	return s, nil
}

// Register HTTP handler
func (s HTTPServer) Register(handlers ...IHandler) error {
	for _, h := range handlers {
		if err := h.Mount(s.apiGroup); err != nil {
			return err
		}
	}
	return nil
}

// Run Server
func (s HTTPServer) ServeForever(errorChan chan<- error) {
	if cap(errorChan) == 0 {
		panic("Capacity of error channel shoule > 0")
	}

	go func(errorChan chan<- error) {
		if err := s.e.Start(s.Address); err != nil {
			errorChan <- err
		}
	}(errorChan)
}

// ---------------------------------------------------------------------------------------------------------------------
// Handler

type (
	Response struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}

	HTTPHandler        func(echo.Context) (interface{}, error)
	HTTPHandlerWrapper func(HTTPHandler) HTTPHandler
)

func EchoWrapper(handler HTTPHandler, wrappers ...HTTPHandlerWrapper) echo.HandlerFunc {
	for _, w := range wrappers {
		handler = w(handler)
	}

	return func(ctx echo.Context) error {
		// add request_id to request context
		rid := ctx.Request().Header.Get(echo.HeaderXRequestID)
		if rid != "" {
			c := context.WithValue(ctx.Request().Context(), ContextRequestID, rid)
			ctx.SetRequest(ctx.Request().WithContext(c))
		}

		res, err := handler(ctx)

		if err != nil {
			resp := Response{
				Code:    50000,
				Message: err.Error(),
				Data:    nil,
			}
			err = ctx.JSON(http.StatusOK, resp)
		} else {
			resp := Response{
				Code:    20000,
				Message: "OK",
				Data:    res,
			}
			err = ctx.JSON(http.StatusOK, resp)
		}
		if err != nil {
			str, ok := res.([]byte)
			if ok {
				WithError(err).WithField("result", str).Error("Error in http request")
			}
			WithError(err).WithField("result", res).Error("Error in http request")
		}
		return err
	}
}

type CompressedResult struct {
	Fields []string      `json:"fields"`
	Values []interface{} `json:"values"`
}

func WrapperJSONCompress(handler HTTPHandler) HTTPHandler {
	return func(c echo.Context) (interface{}, error) {
		res, err := handler(c)

		resmap := make([]map[string]interface{}, 0)
		if err := JSONTranslate(res, &resmap); err != nil {
			return nil, err
		}
		if len(resmap) == 0 {
			return []int{}, nil
		}

		comres := new(CompressedResult)
		for k := range resmap[0] {
			comres.Fields = append(comres.Fields, k)
		}

		for _, item := range resmap {
			var values []interface{}
			for _, value := range item {
				values = append(values, value)
			}
			comres.Values = append(comres.Values, values)
		}

		return comres, err
	}
}

// ------------
func ApiMetricsToRedisMiddle(rp *redis.Pool) echo.MiddlewareFunc {
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
				Info(err)
				return err
			}

			return next(c)

		}
	}
}
