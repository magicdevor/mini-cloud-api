package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"

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

	SelfWXOpenid     string `json:"self_wx_openid"`
	SelfWXSessionKey string `json:"self_wx_session_key"`
}

func (s *Server) Login(ctx *gin.Context) {
	var params LoginParams
	err := ctx.ShouldBindHeader(&params)
	params.SelfWXOpenid = ctx.MustGet("self_wx_openid").(string)
	params.SelfWXSessionKey = ctx.MustGet("self_wx_session_key").(string)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": err.Error()})
		return
	}

	arg := db.GetUserIDByAppidAndOpenidParams{
		Appid:  GetString(params.XWXAppid, s.config.WXAppid),
		Openid: params.SelfWXOpenid,
	}

	id, err := s.db.GetUserIDByAppidAndOpenid(ctx, arg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			token, err := LoginWithOpenData(s, ctx, params)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "token": token, "message": "success"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	uarg := db.UpdateUserParams{
		ID:         id,
		SessionKey: params.SelfWXSessionKey,
	}
	if err = s.db.UpdateUser(ctx, uarg); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	token, err := s.maker.CreateToken(id, s.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "token": token, "message": "success"})
}

func LoginWithOpenData(s *Server, ctx *gin.Context, args LoginParams) (token string, err error) {
	id, err := s.snowflake.NextID()
	if err != nil {
		return
	}
	carg := db.CreateUserParams{
		ID:          strconv.FormatInt(id, 16),
		Unionid:     GetString(args.XWXUnionid, "unionid"),
		Appid:       GetString(args.XWXAppid, s.config.WXAppid),
		Openid:      args.SelfWXOpenid,
		SessionKey:  args.SelfWXSessionKey,
		AppidFrom:   GetString(args.XWXAppidFrom, s.config.WXAppid),
		UnionidFrom: GetString(args.XWXUnionidFrom, "unionid"),
		OpenidFrom:  GetString(args.XWXOpenidFrom, args.SelfWXOpenid),
	}
	log.Printf("appid: %s", carg.Appid)
	u, err := s.db.CreateUser(ctx, carg)

	if err != nil {
		return
	}
	token, err = s.maker.CreateToken(u.ID, s.config.AccessTokenDuration)
	return token, err
}

func GetString(origin, target string) string {
	if origin != "" {
		return origin
	}
	return target
}
