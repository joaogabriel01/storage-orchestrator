package strategies

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var sequentialSaveStrategy SequentialSaveStrategy[string, string]
var parallelSaveStrategy ParallelSaveStrategy[string, string]

func sequentialSaveSetup() {
	sequentialSaveStrategy = SequentialSaveStrategy[string, string]{}
	initialSetup()
}

func parallelSaveSetup() {
	parallelSaveStrategy = ParallelSaveStrategy[string, string]{}
	initialSetup()
}

func TestSequentialSave(t *testing.T) {

	t.Run("should return all units and without error", func(t *testing.T) {
		sequentialSaveSetup()
		ctx := context.Background()

		mock1.On("Save", "query", "worked", mock.Anything).Return(nil)
		mock2.On("Save", "query", "worked", mock.Anything).Return(nil)

		saved, err := sequentialSaveStrategy.Save(ctx, "query", "worked", units, targets)
		assert.NoError(t, err)
		assert.ElementsMatch(t, saved, []string{"mock1", "mock2"})

		mock1.AssertExpectations(t)
		mock2.AssertExpectations(t)
	})

	t.Run("should return error when any unit fails", func(t *testing.T) {
		sequentialSaveSetup()
		ctx := context.Background()

		mock1.On("Save", "query", "saved", mock.Anything).Return(fmt.Errorf("unit1 error"))
		mock2.AssertNotCalled(t, "Save")

		saved, err := sequentialSaveStrategy.Save(ctx, "query", "saved", units, targets)

		assert.ErrorContains(t, err, "error saving unit mock1: unit1 error")

		assert.ElementsMatch(t, saved, []string{})

		mock1.AssertExpectations(t)
		mock2.AssertExpectations(t)
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
		sequentialSaveSetup()
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		cancel()

		saved, err := sequentialSaveStrategy.Save(ctx, "query", "saved", units, targets)

		assert.ErrorIs(t, err, context.Canceled)

		assert.ElementsMatch(t, saved, []string{})

	})

}

func TestParallelSave(t *testing.T) {

	t.Run("should return all units and without error", func(t *testing.T) {
		parallelSaveSetup()
		ctx := context.Background()

		mock1.On("Save", "query", "worked", mock.Anything).Return(nil)
		mock2.On("Save", "query", "worked", mock.Anything).Return(nil)

		saved, err := parallelSaveStrategy.Save(ctx, "query", "worked", units, targets)
		assert.NoError(t, err)
		assert.ElementsMatch(t, saved, []string{"mock1", "mock2"})

		mock1.AssertExpectations(t)
		mock2.AssertExpectations(t)
	})

	t.Run("should return error when any unit fails but how it is parallel, it should return the saved units", func(t *testing.T) {
		parallelSaveSetup()
		ctx := context.Background()

		mock1.On("Save", "query", "saved", mock.Anything).Return(fmt.Errorf("unit1 error"))
		mock2.On("Save", "query", "saved", mock.Anything).Return(nil)

		saved, err := parallelSaveStrategy.Save(ctx, "query", "saved", units, targets)

		assert.ErrorContains(t, err, "error saving unit mock1: unit1 error")

		assert.ElementsMatch(t, saved, []string{"mock2"})

		mock1.AssertExpectations(t)
		mock2.AssertExpectations(t)
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
		parallelSaveSetup()
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		cancel()

		saved, err := parallelSaveStrategy.Save(ctx, "query", "saved", units, targets)

		assert.ErrorIs(t, err, context.Canceled)

		assert.ElementsMatch(t, saved, []string{})

	})

}
