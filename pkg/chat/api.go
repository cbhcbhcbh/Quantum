package chat

import (
	"net/http"
	"strconv"

	"github.com/cbhcbhcbh/Quantum/pkg/common/known"
	"github.com/cbhcbhcbh/Quantum/pkg/common/response"
	"github.com/gin-gonic/gin"
)

func (h *HttpServer) StartChat(c *gin.Context) {
	uid := c.Query("uid")
	userID, err := strconv.ParseUint(uid, 10, 64)
	if err != nil {
		response.ErrorResponse(http.StatusBadRequest, known.ErrInvalidParam.Error()).SetHttpCode(http.StatusUnauthorized).WriteTo(c)
		return
	}

	r.
}

func (h *HttpServer) GetChannelUsers(c *gin.Context) {

}

func (h *HttpServer) GetOnlineUsers(c *gin.Context) {

}

func (h *HttpServer) ListMessages(c *gin.Context) {

}

func (h *HttpServer) DeleteChannel(c *gin.Context) {

}
