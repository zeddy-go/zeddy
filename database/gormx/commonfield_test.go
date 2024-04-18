package gormx

import (
	"github.com/sony/sonyflake"
	"github.com/stretchr/testify/require"
	"github.com/zeddy-go/zeddy/container"
	"testing"
)

func TestSnowflakeID_BeforeCreate(t *testing.T) {
	err := container.Bind[*sonyflake.Sonyflake](func() (*sonyflake.Sonyflake, error) {
		return sonyflake.New(sonyflake.Settings{})
	})
	require.NoError(t, err)
	s := SnowflakeID{}
	err = s.BeforeCreate(nil)
	require.NoError(t, err)
	require.NotEmpty(t, s.ID)
}
