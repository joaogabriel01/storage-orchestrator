package strategies

import (
	"context"
	"fmt"
	"testing"

	strategies_mock "github.com/joaogabriel01/storage-orchestrator/pkg/strategies/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var cacheGetStrategy CacheGetStrategy[string, string]
var saveMock *strategies_mock.MockSaveStrategy

func cacheGetSetup() {
	cacheGetStrategy = CacheGetStrategy[string, string]{}
	saveMock = &strategies_mock.MockSaveStrategy{}
	initialSetup()
}

func TestGet(t *testing.T) {

	t.Run("should return error when the save function is not passed", func(t *testing.T) {
		cacheGetSetup()
		ctx := context.Background()

		_, err := cacheGetStrategy.Get(ctx, "query", units, targets)

		assert.ErrorContains(t, err, "save function not found")

	})

	t.Run("should return error when the save function is not the correct type", func(t *testing.T) {
		cacheGetSetup()
		ctx := context.Background()

		_, err := cacheGetStrategy.Get(ctx, "query", units, targets, "not the correct type")

		assert.ErrorContains(t, err, "save function check did not work")

	})

	t.Run("should return a valid value from the first unit when passed order and is of cache type", func(t *testing.T) {
		cacheGetSetup()
		ctx := context.Background()

		mock1.On("Get", "query", mock.Anything).Return("worked", nil)

		value, err := cacheGetStrategy.Get(ctx, "query", units, targets, saveMock)

		assert.NoError(t, err)
		assert.Equal(t, "worked", value)

		mock1.AssertExpectations(t)
	})

	t.Run("should return a valid value even when it is not the first unit", func(t *testing.T) {
		cacheGetSetup()
		ctx := context.Background()
		saveMock.On("Save", mock.Anything, "query", "worked", units, []string{"mock1"}, mock.Anything).Return([]string{"mock1"}, nil)

		mock1.On("Get", "query", mock.Anything).Return("", fmt.Errorf("value not found"))

		mock2.On("Get", "query", mock.Anything).Return("worked", nil)

		value, err := cacheGetStrategy.Get(ctx, "query", units, targets, saveMock)

		assert.Equal(t, "worked", value)
		assert.NoError(t, err)

		mock1.AssertExpectations(t)
		mock2.AssertExpectations(t)
	})

	t.Run("should return an error and does not execute the save method when it was not found in any unit", func(t *testing.T) {
		cacheGetSetup()
		ctx := context.Background()

		mock1.On("Get", "query", mock.Anything).Return("", fmt.Errorf("value not found"))

		mock2.On("Get", "query", mock.Anything).Return("", fmt.Errorf("value not found"))

		value, err := cacheGetStrategy.Get(ctx, "query", units, targets, saveMock)

		assert.Equal(t, "", value)
		assert.ErrorContains(t, err, "no unit returned")

		mock1.AssertExpectations(t)
		mock2.AssertExpectations(t)
	})

	t.Run("should return an error when the save method return an error", func(t *testing.T) {
		cacheGetSetup()
		ctx := context.Background()
		saveMock.On("Save", mock.Anything, "query", "worked", units, []string{"mock1"}, mock.Anything).Return([]string{}, fmt.Errorf("didnt'save mock1"))

		mock1.On("Get", "query", mock.Anything).Return("", fmt.Errorf("value not found"))

		mock2.On("Get", "query", mock.Anything).Return("worked", nil)

		value, err := cacheGetStrategy.Get(ctx, "query", units, targets, saveMock)

		assert.Equal(t, "worked", value)
		assert.ErrorContains(t, err, "didnt'save mock1")

		mock1.AssertExpectations(t)
		mock2.AssertExpectations(t)
	})

}
