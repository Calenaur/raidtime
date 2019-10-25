package handler

import (
	"github.com/labstack/echo"
	"github.com/calenaur/raidtime/store"
	"github.com/calenaur/raidtime/model"
	"github.com/calenaur/raidtime/config"
)

const CODE_OK = 200
const CODE_ERROR_INVALID_ARGUMENTS = 400
const CODE_ERROR_NO_SESSION = 401
const CODE_ERROR_NO_SIGNUP = 402
const CODE_ERROR_INTERNAL_SERVER_ERROR = 500

type Handler struct {
	us *store.UserStore
	es *store.EventStore
	ds *store.DiscordStore
	cfg *config.Config
}

func New(userStore *store.UserStore, eventStore *store.EventStore, discordStore *store.DiscordStore, config *config.Config) *Handler {
	return &Handler{
		us: userStore,
		es: eventStore,
		ds: discordStore,
		cfg: config,
	}
}

func (h *Handler) RegisterRoutes(e *echo.Echo) {
	//Pages
	e.GET("/", h.IndexHandler)
	e.GET("/auth", h.AuthenticationHandler)
	e.GET("/dashboard", h.DashboardHandler)
	e.GET("/calendar/:offset", h.CalendarHandler)
	e.GET("/signup/:event/:type", h.SignupHandler)
	e.GET("/signup/cancel/:event/:type", h.SignupHandler)
	e.GET("/logout", h.LogoutHandler)

	//Statics
	e.File("/static/css", "static/css/style.css")
	e.File("/static/css/index", "static/css/index.css")
	e.File("/static/css/dashboard", "static/css/dashboard.css")
	e.File("/static/image/wow_logo", "static/image/wow_logo.png")
	e.File("/static/image/wow_logo_inverted", "static/image/wow_logo_inverted.png")
	e.File("/static/javascript/sidenav", "static/javascript/sidenav.js")
	e.File("/static/javascript/default", "static/javascript/default.js")
	e.File("/static/javascript/dashboard", "static/javascript/dashboard.js")
}

func (h *Handler) GetUserByContext(c echo.Context) (*model.User, error) {
	cookie, err := c.Cookie("session")
	if err != nil {
		return nil, err
	}

	session := cookie.Value
	user, err := h.us.GetBySession(session)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (h *Handler) HasValidSessionToken(c echo.Context) bool {
	cookie, err := c.Cookie("session")
	if err != nil {
		return false
	}

	return h.us.ValidateSession(cookie.Value)
}