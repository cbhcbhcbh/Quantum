package message

import (
	"net/http"
	"strconv"

	"github.com/cbhcbhcbh/Quantum/internal/pkg/date"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/known"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/log"
	v1 "github.com/cbhcbhcbh/Quantum/pkg/api/v1"
	"github.com/cbhcbhcbh/Quantum/pkg/response"
	"github.com/gin-gonic/gin"
)

// TODO: 基于ID的分页（游标分页）
func (mc *MessageController) PrivateMessage(c *gin.Context) {
	formId := c.GetInt64(known.XIdKey)
	toIdStr, exists := c.GetQuery("to_id")
	if !exists {
		response.FailResponse(http.StatusInternalServerError, "to_id is required").ToJson(c)
		return
	}

	toId, err := strconv.ParseInt(toIdStr, 10, 64)
	if err != nil {
		response.FailResponse(http.StatusInternalServerError, "to_id must be a valid int64").ToJson(c)
		return
	}

	messages, total, err := mc.b.Message().PrivateMessage(c, formId, toId)
	if err != nil {
		response.FailResponse(http.StatusInternalServerError, "failed to retrieve messages").ToJson(c)
		return
	}

	response.SuccessResponse(gin.H{
		"message": messages,
		"nums":    total,
	}, http.StatusOK).ToJson(c)
}

func (mc *MessageController) GroupMessage(c *gin.Context) {

}

func (mc *MessageController) SendMessage(c *gin.Context) {
	log.C(c).Infow("MessageController SendMessage function called")

	var form v1.PrivateMessageRequest
	if err := c.ShouldBindJSON(&form); err != nil {
		response.FailResponse(http.StatusBadRequest, "Invalid input").ToJson(c)
		return
	}

	form.SendTime = date.NewDate()

	if code, err := mc.b.Message().SendMessage(c, form); err != nil {
		response.FailResponse(code, err.Error()).ToJson(c)
		return
	}
	response.SuccessResponse(form).ToJson(c)
}

func (mc *MessageController) SendVideoMessage(c *gin.Context) {

}

func (mc *MessageController) RecallMessage(c *gin.Context) {

}
