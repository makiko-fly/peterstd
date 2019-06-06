package http

import (
	"net/http"

	"github.com/labstack/echo"
)

type Handler interface {
	Mount(e *echo.Group) error
}

type (
	Response struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}

	HTTPHandler        func(echo.Context) (interface{}, error)
	HTTPHandlerWrapper func(HTTPHandler) HTTPHandler
)

func EchoHandler(handler HTTPHandler, wrappers ...HTTPHandlerWrapper) echo.HandlerFunc {
	for _, w := range wrappers {
		handler = w(handler)
	}

	return func(ctx echo.Context) error {
		res, err := handler(ctx)
		if err != nil {
			resp := Response{
				Code:    50000,
				Message: err.Error(),
				Data:    nil,
			}
			return ctx.JSON(http.StatusInternalServerError, resp)
		}
		resp := Response{
			Code:    20000,
			Message: "OK",
			Data:    res,
		}
		return ctx.JSON(http.StatusOK, resp)
	}
}
