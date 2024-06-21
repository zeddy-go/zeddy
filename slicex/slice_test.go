package slicex

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestContains(t *testing.T) {
	t.Run("find slice item in string", func(t *testing.T) {
		data := []string{" or ", " and "}
		found := Contains("id = 1 or no = 2", data)
		require.False(t, found)
	})
}

func TestMapBy(t *testing.T) {
	type Item struct {
		A int
	}
	s := []Item{
		{
			A: 1,
		},
		{
			A: 2,
		},
		{
			A: 3,
		},
	}
	require.Equal(t, map[int]Item{1: s[0], 2: s[1], 3: s[2]}, MapBy[int]("A", s))
}

func TestGetFieldSlice(t *testing.T) {
	type Item2 struct {
		A int
	}
	type Item1 struct {
		Item *Item2
	}
	s := []Item1{
		{
			Item: &Item2{
				A: 1,
			},
		},
		{
			Item: &Item2{
				A: 2,
			},
		},
		{
			Item: &Item2{
				A: 3,
			},
		},
	}

	require.Equal(t, []int{1, 2, 3}, GetFieldSlice[int]("Item.A", s))
}
