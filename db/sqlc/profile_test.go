package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"xiebeitech.com/mini-cloud-api/util"
)

func NewProfile() (CreateProfileParams, error) {
	user, _ := NewUser()
	arg := CreateProfileParams{
		ID:        util.RandomOpenid(),
		UserID:    user.ID,
		Nickname:  util.RandomNicname(),
		AvatarUrl: util.RandomAvatarUrl(),
		Gender:    util.RandomGender(),
	}
	err := testQueries.CreateProfile(context.Background(), arg)
	return arg, err
}

func TestCreateProfile(t *testing.T) {
	_, err := NewProfile()
	require.NoError(t, err)
}

func TestGetProfile(t *testing.T) {
	arg, err := NewProfile()
	require.NoError(t, err)

	p, err := testQueries.GetProfile(context.Background(), arg.UserID)
	require.NoError(t, err)
	require.NotEmpty(t, p)

	require.Equal(t, arg.Nickname, p.Nickname)
	require.Equal(t, arg.AvatarUrl, p.AvatarUrl)
}

func TestGetProfileIDByUserId(t *testing.T) {
	arg, err := NewProfile()
	require.NoError(t, err)

	id, err := testQueries.GetProfileIDByUserId(context.Background(), arg.UserID)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	require.Equal(t, arg.ID, id)
}

func TestUpdateProfile(t *testing.T) {
	arg, err := NewProfile()
	require.NoError(t, err)

	profile := UpdateProfileParams{
		Nickname:  util.RandomNicname(),
		AvatarUrl: util.RandomAvatarUrl(),
		Gender:    util.RandomGender(),
		UserID:    arg.UserID,
	}

	err = testQueries.UpdateProfile(context.Background(), profile)
	require.NoError(t, err)
}
