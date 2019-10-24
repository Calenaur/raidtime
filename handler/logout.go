package handler

import (
	"time"
	"net/http"
	"github.com/labstack/echo"
)

func (h *Handler) LogoutHandler(c echo.Context) error {

	if !h.HasValidSessionToken(c) {
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	user, err := h.GetUserByContext(c)
	if err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	cookie, err := c.Cookie("session")
	if err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	cookie.Expires = time.Now()
	c.SetCookie(cookie)
	h.us.Logout(user)
	return c.Redirect(http.StatusTemporaryRedirect, "/")
}