package gormx

import (
	"github.com/stretchr/testify/require"
	"github.com/zeddy-go/zeddy/database"
	"gorm.io/gorm"
	"gorm.io/gorm/utils/tests"
	"testing"
)

func TestCondition(t *testing.T) {
	db, err := gorm.Open(tests.DummyDialector{})
	db.DryRun = true
	type testModel struct{}
	require.NoError(t, err)
	t.Run("normal", func(t *testing.T) {
		db, err := applyCondition(db, []any{"id", 1})
		require.NoError(t, err)
		db.Find(&testModel{})
		require.Equal(t, "SELECT * FROM `test_models` WHERE `id` = ?", db.Statement.SQL.String())
		require.Equal(t, 1, db.Statement.Vars[0])
	})

	t.Run("compare", func(t *testing.T) {
		db, err := applyCondition(db, []any{"id", ">", 1})
		require.NoError(t, err)
		db.Find(&testModel{})
		require.Equal(t, "SELECT * FROM `test_models` WHERE `id` > (?)", db.Statement.SQL.String())
		require.Equal(t, 1, db.Statement.Vars[0])
	})

	t.Run("like", func(t *testing.T) {
		db, err := applyCondition(db, []any{"id", "like", "1"})
		require.NoError(t, err)
		db.Find(&testModel{})
		require.Equal(t, "SELECT * FROM `test_models` WHERE `id` LIKE (?)", db.Statement.SQL.String())
		require.Equal(t, "%1%", db.Statement.Vars[0])
	})

	t.Run("and", func(t *testing.T) {
		db, err := applyCondition(db, []database.Condition{
			{"id", 1},
			{"no", "2"},
		}...)
		require.NoError(t, err)
		db.Find(&testModel{})
		require.Equal(t, "SELECT * FROM `test_models` WHERE `id` = ? AND `no` = ?", db.Statement.SQL.String())
		require.Equal(t, []any{1, "2"}, db.Statement.Vars)
	})

	t.Run("or", func(t *testing.T) {
		db, err := applyCondition(db, []any{"id = ? or no = ?", 1, "2"})
		require.NoError(t, err)
		db.Find(&testModel{})
		require.Equal(t, "SELECT * FROM `test_models` WHERE id = ? or no = ?", db.Statement.SQL.String())
		require.Equal(t, []any{1, "2"}, db.Statement.Vars)
	})

	t.Run("or_and", func(t *testing.T) {
		db, err := applyCondition(db, []database.Condition{
			{"id = ? or no = ?", 1, "2"},
			{"id != ? or no != ?", 1, "2"},
		}...)
		require.NoError(t, err)
		db.Find(&testModel{})
		require.Equal(t, "SELECT * FROM `test_models` WHERE (id = ? or no = ?) AND (id != ? or no != ?)", db.Statement.SQL.String())
		require.Equal(t, []any{1, "2", 1, "2"}, db.Statement.Vars)
	})
}
