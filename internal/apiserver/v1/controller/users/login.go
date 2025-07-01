package users

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/log"
	v1 "github.com/cbhcbhcbh/Quantum/pkg/api/v1"
	"github.com/cbhcbhcbh/Quantum/pkg/response"
	"github.com/gin-gonic/gin"
)

func (uc *UserController) Login(c *gin.Context) {
	log.C(c).Infow("Login function called")

	var r v1.LoginRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		response.FailResponse(http.StatusInternalServerError, err.Error()).WriteTo(c)

		return
	}

	if _, err := govalidator.ValidateStruct(r); err != nil {
		response.FailResponse(http.StatusInternalServerError, err.Error()).WriteTo(c)

		return
	}

	resp, err := uc.b.Users().Login(c, &r)
	if err != nil {
		response.FailResponse(http.StatusInternalServerError, err.Error()).WriteTo(c)

		return
	}

	response.SuccessResponse(resp).WriteTo(c)
}