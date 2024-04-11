package tool

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestOr(t *testing.T) {
	require.Equal(t, 1, Or(0, 1))
	require.Equal(t, "1", Or("", "1"))
}
