package json

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestJson(t *testing.T) {
	j := `{"time":"2024-02-26 16:03:40"}`

	a := struct {
		T Time `json:"time"`
	}{}

	err := json.Unmarshal([]byte(j), &a)
	require.NoError(t, err)

	tm, err := time.ParseInLocation(time.DateTime, "2024-02-26 16:03:40", time.Local)
	require.NoError(t, err)

	require.True(t, a.T.Equal(tm))

	js, err := json.Marshal(a)
	require.NoError(t, err)

	require.Equal(t, j, string(js))
}

func TestNull(t *testing.T) {
	j := `{"time":""}`

	a := struct {
		T *Time `json:"time"`
	}{}

	err := json.Unmarshal([]byte(j), &a)
	require.NoError(t, err)

	fmt.Printf("%+v\n", a)
}
