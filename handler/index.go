package handler

import (
	"net/http"
	"github.com/labstack/echo"
)

func (h *Handler) IndexHandler(c echo.Context) error {

	if h.HasValidSessionToken(c) {
		return c.Redirect(http.StatusTemporaryRedirect, "/dashboard")
	}

	data := echo.Map {
		"header": false,
		"footer": false,
		"title": "Raid Time!",
		"RedirectUri": h.cfg.Discord.RedirectUri,
	}
	return c.Render(http.StatusOK, "index", data)
}