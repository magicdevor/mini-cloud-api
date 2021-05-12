package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	db "xiebeitech.com/mini-cloud-api/db/sqlc"
)

type LoginParams struct {
	XWXAppid       string `json:"X-WX-APPID"`
	XWXOpenid      string `json:"X-WX-OPENID"`
	XWXUnionid     string `json:"X-WX-UNIONID"`
	XWXAppidFrom   string `json:"X-WX-FROM-APPID"`
	XWXOpenidFrom  string `json:"X-WX-FROM-OPENID"`
	XWXUnionidFrom string `json:"X-WX-FROM-UNIONID"`
	XWXSource      string `json:"X-WX-SOURCE"`

	SelfWXOpenid     string `json:"self_wx_openid" binding:"required"`
	SelfWXSessionKey string `json:"self_wx_session_key" binding:"required"`
}

func (s *Server) Login(ctx *gin.Context) {
	var params LoginParams
	err := ctx.ShouldBindHeader(&params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": err.Error()})
		return
	}

	arg := db.GetUserIDByAppidAndOpenidParams{
		Appid:  params.XWXAppid,
		Openid: params.SelfWXOpenid,
	}
	id, err := s.db.GetUserIDByAppidAndOpenid(ctx, arg)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": err.Error()})
		return
	}

	carg := db.CreateUserParams{ID: id}
	s.db.CreateUser(ctx, carg)
}
