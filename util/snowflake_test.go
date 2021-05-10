package util

import (
	"errors"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSnowflake(t *testing.T) {
	var s []int64
	c := make(chan int64)
	sf := NewSnowflake(10)

	go Generate(c, sf)
	go Generate(c, sf)
	go Generate(c, sf)
	go Generate(c, sf)

	go func() {
		for {
			s = append(s, <-c)
		}
	}()

	time.Sleep(1e9)
	err := IsRepeat(s)
	require.NoError(t, err)
	require.NotEmpty(t, s)
	require.Greater(t, len(s), 10000)
}

func Generate(c chan int64, sf *Snowflake) {
	for i := 0; i < 20000; i++ {
		id, err := sf.NextID()
		if err != nil {
			log.Fatalln("cannot new id:", err)
		}
		c <- id
	}
}

func IsRepeat(s []int64) error {
	for i := 0; i < len(s); i++ {
		for n := i + 1; n < len(s); n++ {
			if s[i] == s[n] {
				return errors.New("is repeat")
			}
		}
	}
	return nil
}
