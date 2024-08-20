package pkg

import (
	"context"
	"errors"
	"testing"

	"github.com/joaogabriel01/storage-orchestrator/pkg/protocols"
	"github.com/joaogabriel01/storage-orchestrator/pkg/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestOrchestratorUnitOperations(t *testing.T) {
	orchestrator := NewOrchestrator[string, string](map[string]protocols.StorageUnit[string, string]{})

	t.Run("it makes the storage unit operations", func(t *testing.T) {
		mock1 := test.NewUnitMock()
		mock2 := test.NewUnitMock()

		orchestrator.AddUnit("mock1", mock1)
		orchestrator.AddUnit("mock2", mock2)

		mocks, err := orchestrator.GetUnits()
		assert.NoError(t, err)

		assert.Same(t, mock1, mocks["mock1"])
		assert.Same(t, mock2, mocks["mock2"])

		mock1Received, err := orchestrator.GetUnit("mock1")
		assert.NoError(t, err)
		assert.Same(t, mock1, mock1Received)

		mockNotFound, err := orchestrator.GetUnit("non-existentMock")
		assert.EqualError(t, err, "unit not found")
		assert.Equal(t, nil, mockNotFound)

	})
}

func TestOrchestradorSave(t *testing.T) {
	orchestrator := NewOrchestrator[string, string](map[string]protocols.StorageUnit[string, string]{})

	t.Run("it saves an item with success - sequential mode", func(t *testing.T) {
		mock1 := test.NewUnitMock()
		mock2 := test.NewUnitMock()

		orchestrator.AddUnit("mock1", mock1)
		orchestrator.AddUnit("mock2", mock2)

		mock1.On("Save", "saved", mock.Anything).Return(nil)
		mock2.On("Save", "saved", mock.Anything).Return(nil)

		saved, err := orchestrator.Save("saved")
		assert.NoError(t, err)
		assert.ElementsMatch(t, saved, []string{"mock1", "mock2"})

		mock1.AssertExpectations(t)
		mock2.AssertExpectations(t)
	})

	t.Run("it saves an item with success in just one mock - sequential mode", func(t *testing.T) {
		mock1 := test.NewUnitMock()
		mock2 := test.NewUnitMock()

		orchestrator.AddUnit("mock1", mock1)
		orchestrator.AddUnit("mock2", mock2)

		mock2.On("Save", "saved", mock.Anything).Return(nil)

		saved, err := orchestrator.Save("saved", func(opts *protocols.SaveOptions) {
			opts.Targets = []string{"mock2"}
		})
		assert.NoError(t, err)
		assert.ElementsMatch(t, saved, []string{"mock2"})

		mock1.AssertNotCalled(t, "Save")
		mock2.AssertExpectations(t)
	})

	t.Run("it saves an item with success - parallel mode", func(t *testing.T) {
		mock1 := test.NewUnitMock()
		mock2 := test.NewUnitMock()

		orchestrator.AddUnit("mock1", mock1)
		orchestrator.AddUnit("mock2", mock2)

		mock1.On("Save", "saved", mock.Anything).Return(nil)
		mock2.On("Save", "saved", mock.Anything).Return(nil)

		saved, err := orchestrator.Save("saved", func(opts *protocols.SaveOptions) {
			opts.HowWillItSave = protocols.Parallel
		})
		assert.NoError(t, err)
		assert.ElementsMatch(t, saved, []string{"mock1", "mock2"})

		mock1.AssertExpectations(t)
		mock2.AssertExpectations(t)
	})

	t.Run("it saves an item with success in just one mock - parallel mode", func(t *testing.T) {
		mock1 := test.NewUnitMock()
		mock2 := test.NewUnitMock()

		orchestrator.AddUnit("mock1", mock1)
		orchestrator.AddUnit("mock2", mock2)

		mock1.On("Save", "saved", mock.Anything).Return(nil)
		mock2.On("Save", "saved", mock.Anything).Return(nil)

		saved, err := orchestrator.Save("saved", func(opts *protocols.SaveOptions) {
			opts.HowWillItSave = protocols.Parallel
			opts.Targets = []string{"mock2"}
		})
		assert.NoError(t, err)
		assert.ElementsMatch(t, saved, []string{"mock2"})

		mock1.AssertNotCalled(t, "Save")
		mock2.AssertExpectations(t)
	})

	t.Run("it saves an item with unit error - sequential mode", func(t *testing.T) {
		mock1 := test.NewUnitMock()
		mock2 := test.NewUnitMock()

		orchestrator.AddUnit("mock1", mock1)
		orchestrator.AddUnit("mock2", mock2)

		expectedErr := errors.New("unit1 error")
		mock1.On("Save", "saved", mock.Anything).Return(expectedErr)
		mock2.On("Save", "saved", mock.Anything).Return(nil)

		saved, err := orchestrator.Save("saved", func(opts *protocols.SaveOptions) {
			opts.HowWillItSave = protocols.Sequential
			opts.Targets = []string{
				"mock1",
				"mock2",
			}
		})

		assert.ErrorIs(t, err, expectedErr)

		assert.ElementsMatch(t, saved, []string{})

		mock1.AssertExpectations(t)
	})

	t.Run("it saves an item with unit error - parallel mode", func(t *testing.T) {
		mock1 := test.NewUnitMock()
		mock2 := test.NewUnitMock()

		orchestrator.AddUnit("mock1", mock1)
		orchestrator.AddUnit("mock2", mock2)

		expectedErr := errors.New("unit1 error")
		mock1.On("Save", "saved", mock.Anything).Return(expectedErr)
		mock2.On("Save", "saved", mock.Anything).Return(nil)

		saved, err := orchestrator.Save("saved", func(opts *protocols.SaveOptions) {
			opts.HowWillItSave = protocols.Parallel
		})

		assert.ErrorIs(t, err, expectedErr)

		assert.ElementsMatch(t, saved, []string{"mock2"})

		mock1.AssertExpectations(t)
	})

	t.Run("it saves an item with context error - sequential mode", func(t *testing.T) {
		mock1 := test.NewUnitMock()
		mock2 := test.NewUnitMock()

		orchestrator.AddUnit("mock1", mock1)
		orchestrator.AddUnit("mock2", mock2)

		saved, err := orchestrator.Save("saved", func(opts *protocols.SaveOptions) {
			ctx, cancel := context.WithCancel(context.Background())
			opts.Context = ctx
			cancel()
		})

		assert.ErrorIs(t, err, context.Canceled)

		assert.ElementsMatch(t, saved, []string{})

	})

	t.Run("it saves an item with context error - parallel mode", func(t *testing.T) {
		mock1 := test.NewUnitMock()
		mock2 := test.NewUnitMock()

		orchestrator.AddUnit("mock1", mock1)
		orchestrator.AddUnit("mock2", mock2)

		saved, err := orchestrator.Save("saved", func(opts *protocols.SaveOptions) {
			ctx, cancel := context.WithCancel(context.Background())
			opts.Context = ctx
			cancel()

			opts.HowWillItSave = protocols.Parallel
		})

		assert.ErrorIs(t, err, context.Canceled)

		assert.ElementsMatch(t, saved, []string{})

	})

}

func TestOrchestratorGet(t *testing.T) {
	orchestrator := NewOrchestrator[string, string](map[string]protocols.StorageUnit[string, string]{})

	t.Run("it receives an error when trying to get a Cache item without passing order", func(t *testing.T) {
		mock1 := test.NewUnitMock()
		mock2 := test.NewUnitMock()

		orchestrator.AddUnit("mock1", mock1)
		orchestrator.AddUnit("mock2", mock2)

		_, err := orchestrator.Get("caughtTest")

		assert.Errorf(t, err, "unspecified order")

	})

	t.Run("it receives a valid value from the first unit when passed order and is of cache type", func(t *testing.T) {
		mock1 := test.NewUnitMock()
		mock2 := test.NewUnitMock()

		orchestrator.AddUnit("mock1", mock1)
		orchestrator.AddUnit("mock2", mock2)

		expectedValue := "worked"
		mock1.On("Get", "caughtTest", mock.Anything).Return(expectedValue, nil)

		value, err := orchestrator.Get("caughtTest", func(opts *protocols.GetOptions) {
			opts.Targets = []string{"mock1", "mock2"}
		})

		assert.Equal(t, expectedValue, value)
		assert.NoError(t, err)
	})
}
