package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"xiebeitech.com/mini-cloud-api/token"
	"xiebeitech.com/mini-cloud-api/util"
)

type AuthCode struct {
	Code string `form:"code" binding:"required"`
}

func AuthenticateMiddleware(config util.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var code AuthCode
		err := ctx.ShouldBindQuery(&code)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "code must be provided", "err": err.Error()})
			return
		}

		r, err := Code2Session(code.Code, config)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": err.Error()})
			return
		}

		if r.ErrCode != 0 {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": r.ErrCode, "message": r.ErrMsg})
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
	params := url.Values{}
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

const (
	authorizationHeaderKey  = "Authorization"
	authorizationTypeBearer = "bearer"
	authorizationKey        = "authorization_payload"
)

func JWTAuthenticateMiddleware(maker token.Maker) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err})
			return
		}

		accessToken := fields[1]
		payload, err := maker.VerifyToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.Set(authorizationKey, payload)
		c.Next()
	}
}
