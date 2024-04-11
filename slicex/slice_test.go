package slicex

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestContains(t *testing.T) {
	t.Run("find slice item in string", func(t *testing.T) {
		data := []string{" or ", " and "}
		found := Contains("id = 1 or no = 2", data)
		require.True(t, found)
	})
}
