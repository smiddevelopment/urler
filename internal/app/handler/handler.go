package handler

import (
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/smiddevelopment/urler.git/internal/app/shortener"
)

func EncodeUrl(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Please enter valid body")
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			c.Error(err)

			return
		}
	}(c.Request().Body)
	bodyString := string(body)
	if bodyString != "" {
		c.Response().Header().Set("Content-Type", "text/plain")
		c.Response().Header().Set("Content-Length", "30")
		c.Response().WriteHeader(http.StatusCreated)
		_, err := c.Response().Write([]byte(shortener.EncodeString(bodyString)))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Please enter valid id")
		}
	}
	return echo.NewHTTPError(http.StatusBadRequest, "Please enter not empty body!")
}

func DecodeUrl(c echo.Context) error {
	if c.Request().URL.Path == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Id is empty!")
	}

	c.Response().Header().Set("Content-Type", "text/plain")
	c.Response().Header().Set("Location", shortener.DecodeString(strings.TrimPrefix(c.Request().URL.Path, "/")))
	c.Response().WriteHeader(http.StatusTemporaryRedirect)

	return nil
}
