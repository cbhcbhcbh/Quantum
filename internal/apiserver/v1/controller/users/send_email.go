package users

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/log"
	v1 "github.com/cbhcbhcbh/Quantum/pkg/api/v1"
	"github.com/cbhcbhcbh/Quantum/pkg/response"
	"github.com/gin-gonic/gin"
)

func (uc *UserController) SendEmail(c *gin.Context) {
	log.C(c).Infow("SendEmail function called")

	var r v1.SendEmailRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		response.FailResponse(http.StatusInternalServerError, err.Error()).WriteTo(c)
		return
	}

	if _, err := govalidator.ValidateStruct(r); err != nil {
		response.FailResponse(http.StatusInternalServerError, err.Error()).WriteTo(c)
		return
	}

	if err := uc.b.Users().SendEmail(c, &r); err != nil {
		response.FailResponse(http.StatusInternalServerError, err.Error()).WriteTo(c)
		return
	}

	response.SuccessResponse(nil).WriteTo(c)
}
