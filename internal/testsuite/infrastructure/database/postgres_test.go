package database_test

import (
	"context"
	"os"
	"testing"
	"time"

	"go-vents/internal/domain/types"
	"go-vents/internal/infrastructure/database"

	"go-vents/internal/testsuite/infrastructure/testdb"
	internal_test "go-vents/internal/testsuite/internal"

	"github.com/stretchr/testify/require"
)

func TestEventCreate(t *testing.T) {
	if os.Getenv("DATABASE_DRIVER") == "spanner" {
		t.Skip("Skipping Postgres test when using Spanner driver")
	}
	var event = internal_test.NewRandomEvent()
	err := testdb.GetRepository().Create(context.TODO(), event)
	require.NoError(t, err)

	var result types.Event
	err = testdb.GetDB().First(&result, "id = ?", event.Id).Error
	require.NoError(t, err)
	require.Equal(t, event.Id, result.Id)
	require.Equal(t, event.Slug, result.Slug)
	require.WithinDuration(t, event.CreatedAt, result.CreatedAt, time.Second)
	require.WithinDuration(t, event.UpdatedAt, result.UpdatedAt, time.Second)

	testdb.GetDB().Delete(&result)
}

func TestEventDelete(t *testing.T) {
	if os.Getenv("DATABASE_DRIVER") == "spanner" {
		t.Skip("Skipping Postgres test when using Spanner driver")
	}
	var event = internal_test.NewRandomEvent()
	_ = testdb.GetRepository().Create(context.TODO(), event)

	repository := database.NewEventPostgresRepository(testdb.GetDB())

	id, _ := types.NewSharedId(event.Id)
	err := repository.Delete(context.TODO(), id)
	require.NoError(t, err)
}

func TestEventOne(t *testing.T) {
	if os.Getenv("DATABASE_DRIVER") == "spanner" {
		t.Skip("Skipping Postgres test when using Spanner driver")
	}
	var event = internal_test.NewRandomEvent()
	_ = testdb.GetRepository().Create(context.TODO(), event)

	repository := database.NewEventPostgresRepository(testdb.GetDB())
	id, _ := types.NewSharedId(event.Id)

	result, err := repository.One(context.TODO(), id)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, event.Id, result.Id)
	require.Equal(t, event.Slug, result.Slug)

	err = repository.Delete(context.TODO(), id)
	require.NoError(t, err)
}
