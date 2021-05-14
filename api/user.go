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

	if id, err := s.db.GetUserIDByAppidAndOpenid(ctx, arg); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			token, err := LoginWithOpenData(s, ctx, params)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "token": token, "message": "success"})
			return
		}
		arg := db.UpdateUserParams{
			ID:         id,
			SessionKey: params.SelfWXSessionKey,
		}
		if err = s.db.UpdateUser(ctx, arg); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		token := ""
		ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "token": token, "message": "success"})
		return
	}
}

func LoginWithOpenData(s *Server, ctx *gin.Context, args LoginParams) (token string, err error) {
	carg := db.CreateUserParams{
		Appid:       args.XWXAppid,
		Unionid:     args.XWXUnionid,
		Openid:      args.SelfWXOpenid,
		SessionKey:  args.SelfWXOpenid,
		AppidFrom:   args.XWXAppidFrom,
		OpenidFrom:  args.XWXOpenidFrom,
		UnionidFrom: args.XWXUnionidFrom,
	}
	u, err := s.db.CreateUser(ctx, carg)

	if err != nil {
		return
	}
	id := u.ID
	print(id)
	// generate token
	return "", nil
}
