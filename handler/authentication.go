package handler

import (
	"time"
	"net/http"
	"github.com/labstack/echo"
)

func (h *Handler) AuthenticationHandler(c echo.Context) error {
	code := c.FormValue("code")
	if len(code) < 1 {
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	credentials, err := h.ds.GetCredentialsByCode(code)
	if err != nil {
		return c.String(http.StatusOK, "1: " + err.Error())
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	_, session, err := h.us.Login(credentials)
	if err != nil {
		return c.String(http.StatusOK, "2: " + err.Error())
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	cookie := new(http.Cookie)
	cookie.Name = "session"
	cookie.Value = session
	cookie.Expires = time.Now().Add(time.Duration(h.cfg.Session.SessionDuration) * time.Second)
	c.SetCookie(cookie)

	return c.Redirect(http.StatusTemporaryRedirect, "/dashboard")
}