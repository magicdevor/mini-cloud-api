package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "xiebeitech.com/mini-cloud-api/db/sqlc"
	"xiebeitech.com/mini-cloud-api/token"
	"xiebeitech.com/mini-cloud-api/util"
)

func (s *Server) GetProfile(ctx *gin.Context) {
	payload := ctx.MustGet(authorizationHeaderKey).(*token.Payload)
	p, err := s.db.GetProfile(ctx, payload.Userame)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "success", "data": p})
}

type ProfileParams struct {
	// UserInfo      string `form:"userInfo"`
	RawData       string `form:"rawData"`
	Signature     string `form:"signature" binding:"required"`
	EncryptedData string `form:"encryptedData" binding:"required"`
	Iv            string `form:"iv" binding:"required"`
	CloudUID      string `form:"cloudUID"`
}

func (s *Server) CreateProfile(ctx *gin.Context) {
	payload := ctx.MustGet(authorizationKey).(*token.Payload)
	var params ProfileParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": err.Error()})
		return
	}
	openData, err := s.db.GetUserOpenDataByID(ctx, payload.Userame)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": err.Error()})
		return
	}
	decrypter := util.NewWXUserDataCrypt(s.config.WXAppid, openData.SessionKey)
	profile, err := decrypter.Decrypt(params.EncryptedData, params.Iv)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": err.Error()})
		return
	}
	id, _ := s.snowflake.NextID()
	arg := db.CreateProfileParams{
		ID:        strconv.FormatInt(id, 16),
		UserID:    payload.Userame,
		Nickname:  profile.NickName,
		AvatarUrl: profile.AvatarURL,
		Gender:    strconv.FormatInt(int64(profile.Gender), 10),
	}
	log.Println(profile)
	if err = s.db.CreateProfile(ctx, arg); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "success"})

}

func (s *Server) UpdateProfile(ctx *gin.Context) {

}
