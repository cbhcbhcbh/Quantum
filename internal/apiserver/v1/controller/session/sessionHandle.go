package session

import (
	"net/http"

	"github.com/cbhcbhcbh/Quantum/internal/pkg/enum"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/known"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/log"
	v1 "github.com/cbhcbhcbh/Quantum/pkg/api/v1"
	"github.com/cbhcbhcbh/Quantum/pkg/response"
	"github.com/gin-gonic/gin"
)

func (sc *SessionController) Index(c *gin.Context) {
	log.C(c).Infow("Session Index function called")

	var sessions *[]v1.SessionDetail
	var err error

	id := c.GetInt64(known.XIdKey)
	if sessions, err = sc.b.Sessions().GetSessions(c, id); err != nil {
		response.FailResponse(http.StatusInternalServerError, err.Error()).ToJson(c)
		return
	}
	response.SuccessResponse(sessions).ToJson(c)
}

func (sc *SessionController) Store(c *gin.Context) {
	log.C(c).Infow("Session Store function called")

	var form v1.SessionForm
	if err := c.ShouldBindJSON(&form); err != nil {
		response.FailResponse(http.StatusBadRequest, "Invalid input").ToJson(c)
		return
	}

	id := c.GetInt64(known.XIdKey)
	session, err := sc.b.Sessions().CreateSession(c, id, form.ToID, form.ChannelType)
	if err != nil {
		response.FailResponse(http.StatusInternalServerError, err.Error()).ToJson(c)
		return
	}
	response.SuccessResponse(session).ToJson(c)
}

func (sc *SessionController) Update(c *gin.Context) {
	log.C(c).Infow("Session Update function called")

	var r v1.SessionInfo
	var err error
	if err = c.ShouldBindUri(&r); err != nil {
		response.FailResponse(enum.ParamError, err.Error()).ToJson(c)
		return
	}

	var form v1.SessionUpdateForm
	if err = c.ShouldBindJSON(&form); err != nil {
		response.FailResponse(http.StatusBadRequest, "Invalid input").ToJson(c)
		return
	}

	err = sc.b.Sessions().UpdateSession(c, r.SessionId, form.TopStatus, form.Note)
	if err != nil {
		response.FailResponse(http.StatusInternalServerError, err.Error()).ToJson(c)
		return
	}

	response.SuccessResponse().ToJson(c)
}

func (sc *SessionController) Delete(c *gin.Context) {
	log.C(c).Infow("Session Delete function called")

	var r v1.SessionInfo
	if err := c.ShouldBindUri(&r); err != nil {
		response.FailResponse(enum.ParamError, err.Error()).ToJson(c)
		return
	}

	err := sc.b.Sessions().DeleteSession(c, r.SessionId)
	if err != nil {
		response.FailResponse(http.StatusInternalServerError, err.Error()).ToJson(c)
		return
	}

	response.SuccessResponse().ToJson(c)
}
