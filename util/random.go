package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomOpenid() string {
	return RandomString(20)
}

func RandomSessionKey() string {
	return RandomString(15)
}

func RandomGender() string {
	genders := []string{"0", "1", "2"}
	n := len(genders)
	return genders[rand.Intn(n)]
}

func RandomNicname() string {
	return RandomString(6)
}

func RandomAvatarUrl() string {
	path := "https://avatar.tobi.sh/tobiaslins.svg?text="
	name := strings.ToUpper(RandomString(2))
	return path + name
}
