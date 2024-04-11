package database

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDsn(t *testing.T) {
	var dsn DSN = "mysql://root:toor@(127.0.0.1:3306)/main?charset=utf8&parseTime=True&loc=Local"
	require.Equal(t, Type("mysql"), dsn.Type())
	require.Equal(t, "mysql://root:toor@(127.0.0.1:3306)/main?charset=utf8&loc=Local&parseTime=True", dsn.Encode())
	require.Equal(t, "root:toor@(127.0.0.1:3306)/main?charset=utf8&loc=Local&parseTime=True", dsn.RemoveSchema())
	require.Equal(t, []string{"mysql", "root:toor@(127.0.0.1:3306)/main?charset=utf8&loc=Local&parseTime=True"}, dsn.split())
}
