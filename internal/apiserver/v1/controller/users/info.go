package users

import (
	"net/http"

	"github.com/cbhcbhcbh/Quantum/internal/pkg/enum"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/log"
	v1 "github.com/cbhcbhcbh/Quantum/pkg/api/v1"
	"github.com/cbhcbhcbh/Quantum/pkg/response"
	"github.com/gin-gonic/gin"
)

func (uc *UserController) Info(c *gin.Context) {
	log.C(c).Infow("Info function called")

	var r v1.Person
	var err error
	if err = c.ShouldBindUri(&r); err != nil {
		response.FailResponse(enum.ParamError, err.Error()).ToJson(c)
		return
	}

	var users *v1.UserDetails
	if users, err = uc.b.Users().GetUserInfo(c, &r); err != nil {
		response.FailResponse(http.StatusInternalServerError, err.Error()).ToJson(c)
		return
	}

	response.SuccessResponse(users).ToJson(c)
}
