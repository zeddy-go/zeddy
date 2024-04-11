package hashx

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHash(t *testing.T) {
	t.Run("sha256", func(t *testing.T) {
		hash, err := Sha256("123456")
		require.NoError(t, err)
		require.True(t, CheckSha256(hash, "123456"))

		hash, err = Sha256("123456", WithKey([]byte("789")))
		require.NoError(t, err)
		require.True(t, CheckSha256(hash, "123456", WithKey([]byte("789"))))
	})

	t.Run("sha1", func(t *testing.T) {
		hash, err := Sha1("123456")
		require.NoError(t, err)
		require.True(t, CheckSha1(hash, "123456"))

		hash, err = Sha1("123456", WithKey([]byte("789")))
		require.NoError(t, err)
		require.True(t, CheckSha1(hash, "123456", WithKey([]byte("789"))))
	})

	t.Run("md5", func(t *testing.T) {
		hash, err := Md5("123456")
		require.NoError(t, err)
		require.True(t, CheckMd5(hash, "123456"))

		hash, err = Md5("123456", WithKey([]byte("789")))
		require.NoError(t, err)
		require.True(t, CheckMd5(hash, "123456", WithKey([]byte("789"))))
	})
}
