package endpoint

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Start init endpoint and all routes
func Start() *echo.Echo {
	e := echo.New()

	// Middlewarre
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Cors
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST},
	}))

	e.GET("/", homePage)
	e.POST("/file", uploadFileHandler)

	// Server
	e.Start(":8080")

	return e
}

func homePage(c echo.Context) error {
	return c.JSON(http.StatusOK, "say hi")
}
