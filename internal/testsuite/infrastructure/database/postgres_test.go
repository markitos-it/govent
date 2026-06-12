package database_test

import (
	"context"
	"testing"
	"time"

	"govent/internal/domain/types"
	"govent/internal/infrastructure/database"

	"govent/internal/testsuite/infrastructure/testdb"
	internal_test "govent/internal/testsuite/internal"

	"github.com/stretchr/testify/require"
)

func TestEventCreate(t *testing.T) {
	var event = internal_test.NewRandomEvent()
	err := testdb.GetRepository().Create(context.TODO(), event)
	require.NoError(t, err)

	var result types.Event
	err = testdb.GetDB().First(&result, "id = ?", event.Id).Error
	require.NoError(t, err)
	require.Equal(t, event.Id, result.Id)
	require.Equal(t, event.Name, result.Name)
	require.WithinDuration(t, event.CreatedAt, result.CreatedAt, time.Second)
	require.WithinDuration(t, event.UpdatedAt, result.UpdatedAt, time.Second)

	testdb.GetDB().Delete(&result)
}

func TestEventDelete(t *testing.T) {
	var event = internal_test.NewRandomEvent()
	_ = testdb.GetRepository().Create(context.TODO(), event)

	repository := database.NewEventPostgresRepository(testdb.GetDB())

	id, _ := types.NewSharedId(event.Id)
	err := repository.Delete(context.TODO(), id)
	require.NoError(t, err)
}

func TestEventOne(t *testing.T) {
	var event = internal_test.NewRandomEvent()
	_ = testdb.GetRepository().Create(context.TODO(), event)

	repository := database.NewEventPostgresRepository(testdb.GetDB())
	id, _ := types.NewSharedId(event.Id)

	result, err := repository.One(context.TODO(), id)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, event.Id, result.Id)
	require.Equal(t, event.Name, result.Name)

	err = repository.Delete(context.TODO(), id)
	require.NoError(t, err)
}
