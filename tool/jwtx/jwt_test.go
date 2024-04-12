package jwtx

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-module/carbon/v2"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestExpire(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		key := "123"
		token, err := NewJwtTokenBuilder().WithKey(key).BuildToken()
		require.NoError(t, err)
		require.NotEmpty(t, token)

		tok, err := Parse(token, key)
		require.NoError(t, err)
		require.Equal(t, 0, len(tok.Claims.(jwt.MapClaims)))
	})

	t.Run("normal2", func(t *testing.T) {
		key := "123"
		token, err := NewJwtTokenBuilder().WithKey(key).WithClaim(JwtID, 1).BuildToken()
		require.NoError(t, err)
		require.NotEmpty(t, token)

		tok, err := Parse(token, key)
		require.NoError(t, err)
		require.Equal(t, 1, len(tok.Claims.(jwt.MapClaims)))
		require.Equal(t, float64(1), tok.Claims.(jwt.MapClaims)[JwtID])
	})

	t.Run("err1", func(t *testing.T) {
		key := "123"
		token, err := NewJwtTokenBuilder().WithKey(key).WithClaim(JwtID, 1).BuildToken()
		require.NoError(t, err)
		require.NotEmpty(t, token)

		_, err = Parse(token, "")
		require.Error(t, err)
		require.Contains(t, err.Error(), "signature is invalid")
	})

	t.Run("expire", func(t *testing.T) {
		key := "123"
		token, err := NewJwtTokenBuilder().WithKey(key).WithClaim(JwtExpiresAt, carbon.Now().SubSeconds(1).Timestamp()).BuildToken()
		require.NoError(t, err)
		require.NotEmpty(t, token)

		_, err = Parse(token, "123")
		require.Error(t, err)
		require.Contains(t, err.Error(), "token is expired")
	})
}
