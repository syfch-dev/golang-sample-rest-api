package auth

import "github.com/labstack/echo/v5"

type Handlers interface {
	Register() echo.HandlerFunc
}
