package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"xiebeitech.com/mini-cloud-api/util"
)

type AuthCode struct {
	Code string `json:"code"`
}

func AuthenticateMiddleware(config util.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var code AuthCode
		err := ctx.ShouldBindQuery(&code)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "code must be provided"})
			return
		}

		r, err := Code2Session(code.Code, config)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": err.Error()})
			return
		}

		if r.ErrCode != 0 {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": r.ErrCode, "message": err.Error()})
			return
		}

		ctx.Set("self_wx_openid", r.Openid)
		ctx.Set("self_wx_session_key", r.SessionKey)
		ctx.Set("self_wx_unionid", r.Unionid)

		ctx.Next()

	}
}

type Code2SessionResult struct {
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
	Openid     string `json:"openid"`
	Unionid    string `json:"unionid"`
	SessionKey string `json:"session_key"`
}

func Code2Session(code string, config util.Config) (r Code2SessionResult, err error) {
	var params url.Values
	params.Set("appid", config.WXAppid)
	params.Set("secret", config.WXSecret)
	params.Set("js_code", code)
	params.Set("grant_type", "authorization_code")

	path, err := url.Parse("https://api.weixin.qq.com/sns/jscode2session")
	if err != nil {
		return
	}

	path.RawQuery = params.Encode()
	response, err := http.Get(path.String())
	if err != nil {
		return
	}

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &r)
	return
}
