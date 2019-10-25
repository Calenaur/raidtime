package handler

import (
	"time"
	"strconv"
	"net/http"
	"github.com/jinzhu/now"
	"github.com/labstack/echo"
	"github.com/calenaur/raidtime/model"
)

type CalendarResponse struct {
	Code int 						`json:"code"`
	Events []*model.Event 			`json:"events,omitempty"`
	Me int64 	 					`json:"me,omitempty"`
	Status string 					`json:"status,omitempty"`
	Error string 					`json:"error,omitempty"`
}

func (h *Handler) CalendarHandler(c echo.Context) error {
	response := &CalendarResponse{}
	user, err := h.GetUserByContext(c)
	if err != nil {
		response.Code = CODE_ERROR_NO_SESSION
		response.Error = "Invalid credentials"
		return c.JSON(http.StatusOK, response)
	}

	o := c.Param("offset")
	if len(o) < 1 {
		o = "0"
	}

	offset, err := strconv.Atoi(o)
	if err != nil {
		response.Code = CODE_ERROR_INVALID_ARGUMENTS
		response.Error = "Invalid arguments"
		return c.JSON(http.StatusOK, response)
	}

	date := now.New(time.Now())
	date = now.New(date.AddDate(0, offset, 0))
	response.Events, err = h.es.GetEventsByDateRange(date.BeginningOfMonth(), date.EndOfMonth())
	if err != nil {
		response.Code = CODE_ERROR_INTERNAL_SERVER_ERROR
		response.Error = "Internal server error"
		return c.JSON(http.StatusOK, response)
	}

	response.Me = user.ID
	response.Code = CODE_OK
	response.Status = "Success"
	return c.JSON(http.StatusOK, response)
}