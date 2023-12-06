package main

import (
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.POST("/order", handlePlaceOrder)

	e.Start(":3000")
}

func handlePlaceOrder(c echo.Context) error {
	return c.JSON(200, "tst")
}
