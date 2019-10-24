package handler

import (
	"fmt"
	"time"
	"strconv"
	"net/http"
	"github.com/jinzhu/now"
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

	date := now.New(time.Now())
	monthOffset := 0 
	m := c.Param("month")
	monthSet := false
	if len(m) > 0 {
		n, err := strconv.Atoi(m)
		if err == nil {
			if n > 10 {
				n = 10
			}
			if n < -10 {
				n = -10
			}
			monthSet = true
			monthOffset = n
			date = now.New(date.AddDate(0, monthOffset, 0))
		}
	}

	events, err := h.es.GetEventsByDateRange(date.BeginningOfMonth(), date.EndOfMonth())
	if err != nil {
		fmt.Println(err.Error())
	}
	//TODO::In hindsight this should be exported to clientside, only send []*Event to javascript
	//This is trashcan code
	currentDay := 1
	days := [][7]*CalendarDay{}
	firstDay := int(date.BeginningOfMonth().Weekday())
	lastDay := date.EndOfMonth().Day()
	for week := 0; week<(firstDay+lastDay-1)/7+1; week++ {
		days = append(days, [7]*CalendarDay{})
		for day, _ := range days[week] {
			if currentDay > lastDay {
				break;
			}
			if week != 0 || day >= firstDay {
				var e []*model.Event
				if events != nil {
					for _, event := range events {
						if event.Date.Day() == currentDay {
							e = append(e, event)
						}
					}
				}
				days[week][day] = &CalendarDay{currentDay, e}
				currentDay += 1
			}
		}
	}

	signedUpTo := []int{}
	for _, event := range events {
		for _, signup := range event.Signups {
			if signup.User.ID == user.ID {
				signedUpTo = append(signedUpTo, event.ID)
			}
		}
	}

	data := echo.Map {
		"header": true,
		"footer": false,
		"title": "Raid Time! - Dashboard",
		"side_nav": true,
		"user": user,
		"days": days,
		"month": date.Format("January 2006"),
		"month_set": monthSet,
		"next_month": monthOffset+1,
		"prev_month": monthOffset-1,
		"parsetime": func(t time.Time) string {
			return t.Format("02/01/2006")
		},
		"signedup": func(id int) bool {
			for _, eventId := range signedUpTo {
				if eventId == id {
					return true
				}
			}
			return false
		},
	}

	return c.Render(http.StatusOK, "dashboard", data)
}