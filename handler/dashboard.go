package handler

import (
	"net/http"
	"github.com/labstack/echo"
	"github.com/calenaur/raidtime/model"
)

type CalendarDay struct {
	Day int
	Events []*model.Event
}

func (h *Handler) DashboardHandler(c echo.Context) error {
	user, err := h.GetUserByContext(c)
	if err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/")
		return c.String(http.StatusOK, err.Error())
	}

	data := echo.Map {
		"header": true,
		"footer": false,
		"title": "Raid Time! - Dashboard",
		"side_nav": true,
		"user": user,
	}
	return c.Render(http.StatusOK, "dashboard", data)
}