package ginx

import (
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/utils/tests"
	"testing"
)

func getDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(&tests.DummyDialector{})
	require.NoError(t, err)
	db.DryRun = true
	return db
}

func TestPage(t *testing.T) {
	type test struct{}
	t.Run("only_limit", func(t *testing.T) {
		db := getDB(t)
		p := Page{
			Size: 15,
		}

		db, err := p.Apply(db)
		require.NoError(t, err)

		db.Find(&test{})

		require.Equal(t, "SELECT * FROM `tests` LIMIT 15", db.Statement.SQL.String())
	})

	t.Run("normal", func(t *testing.T) {
		db := getDB(t)
		p := Page{
			Page: 2,
			Size: 15,
		}

		db, err := p.Apply(db)
		require.NoError(t, err)

		db.Find(&test{})

		require.Equal(t, "SELECT * FROM `tests` LIMIT 15 OFFSET 15", db.Statement.SQL.String())
	})
}

func TestFilters(t *testing.T) {
	type test struct{}
	f := Filters{
		m: map[string]string{
			"a": "!1",
			"b": "~2",
			"c": ">3",
			"d": "<4",
			"e": ">=5",
			"f": "<=6",
			"g": "7",
		},
	}
	result := f.ParseAll()
	require.Contains(t, result, []any{"a", "!=", "1"})
	require.Contains(t, result, []any{"b", "like", "%2%"})
	require.Contains(t, result, []any{"c", ">", "3"})
	require.Contains(t, result, []any{"d", "<", "4"})
	require.Contains(t, result, []any{"e", ">=", "5"})
	require.Contains(t, result, []any{"f", "<=", "6"})
	require.Contains(t, result, []any{"g", "7"})
}

func TestSorts(t *testing.T) {
	type test struct{}
	s := Sorts{
		m: map[string]string{
			"a": "desc",
			"b": "asc",
		},
	}
	db := getDB(t)
	db, err := s.Apply(db)
	require.NoError(t, err)
	db.Find(&test{})
	require.Contains(t, db.Statement.SQL.String(), "SELECT * FROM `tests` ORDER BY")
	require.Contains(t, db.Statement.SQL.String(), "b ASC")
	require.Contains(t, db.Statement.SQL.String(), "a DESC")
}
