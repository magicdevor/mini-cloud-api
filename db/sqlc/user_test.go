package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"xiebeitech.com/mini-cloud-api/util"
)

func NewUser() (CreateUserParams, error) {
	arg := CreateUserParams{
		ID:          util.RandomOpenid(),
		Appid:       util.RandomOpenid(),
		Openid:      util.RandomOpenid(),
		SessionKey:  util.RandomSessionKey(),
		OpenidFrom:  util.RandomOpenid(),
		AppidFrom:   util.RandomOpenid(),
		Unionid:     util.RandomOpenid(),
		UnionidFrom: util.RandomOpenid(),
	}
	_, err := testQueries.CreateUser(context.Background(), arg)
	return arg, err
}

func TestCreateUser(t *testing.T) {
	arg := CreateUserParams{
		ID:          util.RandomOpenid(),
		Appid:       util.RandomOpenid(),
		Openid:      util.RandomOpenid(),
		SessionKey:  util.RandomSessionKey(),
		OpenidFrom:  util.RandomOpenid(),
		AppidFrom:   util.RandomOpenid(),
		Unionid:     util.RandomOpenid(),
		UnionidFrom: util.RandomOpenid(),
	}

	u, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, u)

	require.Equal(t, arg.ID, u.ID)
	require.Equal(t, arg.Appid, u.Appid)
	require.Equal(t, arg.AppidFrom, u.AppidFrom)
	require.Equal(t, arg.Openid, u.Openid)
	require.Equal(t, arg.SessionKey, u.SessionKey)
	require.Equal(t, arg.OpenidFrom, u.OpenidFrom)
	require.Equal(t, arg.Unionid, u.Unionid)
	require.Equal(t, arg.UnionidFrom, u.UnionidFrom)

	require.NotZero(t, u.CreatedAt)
	require.NotZero(t, u.UpdatedAt)
}

func TestGetUserIDByAppidAndOpenid(t *testing.T) {
	user, err := NewUser()
	require.NoError(t, err)

	arg := GetUserIDByAppidAndOpenidParams{
		Appid:  user.Appid,
		Openid: user.Openid,
	}

	id, err := testQueries.GetUserIDByAppidAndOpenid(context.Background(), arg)
	require.NoError(t, err)

	require.Equal(t, user.ID, id)
}

func TestGetUserOpenDataByID(t *testing.T) {
	user, err := NewUser()
	require.NoError(t, err)

	openData, err := testQueries.GetUserOpenDataByID(context.Background(), user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, openData)

	require.Equal(t, user.Openid, openData.Openid)
	require.Equal(t, user.SessionKey, openData.SessionKey)
}

func TestUpdateUser(t *testing.T) {
	user, err := NewUser()
	require.NoError(t, err)

	arg := UpdateUserParams{
		SessionKey: util.RandomSessionKey(),
		ID:         user.ID,
	}
	err = testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
}
