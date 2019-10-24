package handler

import (
	"net/http"
	"github.com/labstack/echo"
	"strconv"
)

const CODE_OK = 200
const CODE_ERROR_INVALID_ARGUMENTS = 400
const CODE_ERROR_NO_SESSION = 401
const CODE_ERROR_NO_SIGNUP = 402

type SignupResponse struct {
	Code int 			`json:"code"`
	Status string 		`json:"status"`
	Error string 		`json:"error"`
}

func (h *Handler) SignupHandler(c echo.Context) error {
	response := &SignupResponse{}
	user, err := h.GetUserByContext(c)
	if err != nil {
		response.Code = CODE_ERROR_NO_SESSION
		response.Error = "Invalid credentials"
		return c.JSON(http.StatusOK, response)
	}

	t := c.Param("type")
	signupType, err := strconv.Atoi(t)
	if len(t) < 1 || err != nil {
		response.Code = CODE_ERROR_INVALID_ARGUMENTS
		response.Error = "Invalid arguments"
		return c.JSON(http.StatusOK, response)
	}

	e := c.Param("event")
	event, err := strconv.Atoi(e)
	if len(e) < 1 || err != nil {
		response.Code = CODE_ERROR_INVALID_ARGUMENTS
		response.Error = "Invalid arguments"
		return c.JSON(http.StatusOK, response)
	}

	if signupType == -1 {
		err = h.us.CancelSignup(user, event)
		if err != nil {
			response.Code = CODE_ERROR_NO_SIGNUP
			response.Error = "Could not cancel signup: No signup found."
			return c.JSON(http.StatusOK, response)
		}
	} else {
		err = h.us.SignupToEvent(user, event, signupType)
		if err != nil {
			response.Code = CODE_ERROR_INVALID_ARGUMENTS
			response.Error = "Invalid arguments"
			return c.JSON(http.StatusOK, response)
		}
	}

	response.Code = CODE_OK
	response.Status = "Success"
	return c.JSON(http.StatusOK, response)
}