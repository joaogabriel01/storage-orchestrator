package pkg

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/joaogabriel01/storage-orchestrator/pkg/protocols"
	"github.com/joaogabriel01/storage-orchestrator/pkg/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupOrchestrator() Orchestrator[string, string] {
	return NewOrchestrator[string, string](map[string]protocols.StorageUnit[string, string]{}, []string{})

}

var orchestrator Orchestrator[string, string]

func setup() {
	orchestrator = setupOrchestrator()
}

func TestOrchestratorUnitOperations(t *testing.T) {

	t.Run("it makes the storage unit operations", func(t *testing.T) {
		setup()
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

	t.Run("should save standard order without error", func(t *testing.T) {
		setup()
		mock1 := test.NewUnitMock()
		mock2 := test.NewUnitMock()

		orchestrator.AddUnit("mock1", mock1)
		orchestrator.AddUnit("mock2", mock2)
		err := orchestrator.SetStandardOrder("mock1", "mock2")
		assert.NoError(t, err)
	})

	t.Run("should return error when there are no units", func(t *testing.T) {
		setup()

		mock1 := test.NewUnitMock()

		orchestrator.AddUnit("mock1", mock1)
		err := orchestrator.SetStandardOrder("mock1", "mock2")
		assert.ErrorContains(t, err, "this unit does not exist: mock2")

	})
}

func TestOrchestradorSave(t *testing.T) {

	t.Run("it saves an item with success - sequential mode", func(t *testing.T) {
		setup()
		mock1 := test.NewUnitMock()
		mock2 := test.NewUnitMock()

		orchestrator.AddUnit("mock1", mock1)
		orchestrator.AddUnit("mock2", mock2)

		mock1.On("Save", "query", "saved", mock.Anything).Return(nil)
		mock2.On("Save", "query", "saved", mock.Anything).Return(nil)

		orchestrator.SetStandardOrder("mock1", "mock2")
		saved, err := orchestrator.Save("query", "saved")
		assert.NoError(t, err)
		assert.ElementsMatch(t, saved, []string{"mock1", "mock2"})

		mock1.AssertExpectations(t)
		mock2.AssertExpectations(t)
	})

	t.Run("it saves an item with success in just one mock - sequential mode", func(t *testing.T) {
		setup()
		mock1 := test.NewUnitMock()
		mock2 := test.NewUnitMock()

		orchestrator.AddUnit("mock1", mock1)
		orchestrator.AddUnit("mock2", mock2)

		mock2.On("Save", "query", "saved", mock.Anything).Return(nil)

		saved, err := orchestrator.Save("query", "saved", func(opts *protocols.SaveOptions) {
			opts.Targets = []string{"mock2"}
		})
		assert.NoError(t, err)
		assert.ElementsMatch(t, saved, []string{"mock2"})

		mock1.AssertNotCalled(t, "Save")
		mock2.AssertExpectations(t)
	})

	t.Run("it saves an item with success - parallel mode", func(t *testing.T) {
		setup()
		mock1 := test.NewUnitMock()
		mock2 := test.NewUnitMock()

		orchestrator.AddUnit("mock1", mock1)
		orchestrator.AddUnit("mock2", mock2)

		mock1.On("Save", "query", "saved", mock.Anything).Return(nil)
		mock2.On("Save", "query", "saved", mock.Anything).Return(nil)
		orchestrator.SetStandardOrder("mock1", "mock2")

		saved, err := orchestrator.Save("query", "saved", func(opts *protocols.SaveOptions) {
			opts.HowWillItSave = protocols.Parallel
		})
		assert.NoError(t, err)
		assert.ElementsMatch(t, saved, []string{"mock1", "mock2"})

		mock1.AssertExpectations(t)
		mock2.AssertExpectations(t)
	})

	t.Run("it saves an item with success in just one mock - parallel mode", func(t *testing.T) {
		setup()
		mock1 := test.NewUnitMock()
		mock2 := test.NewUnitMock()

		orchestrator.AddUnit("mock1", mock1)
		orchestrator.AddUnit("mock2", mock2)

		mock1.On("Save", "query", "saved", mock.Anything).Return(nil)
		mock2.On("Save", "query", "saved", mock.Anything).Return(nil)

		saved, err := orchestrator.Save("query", "saved", func(opts *protocols.SaveOptions) {
			opts.HowWillItSave = protocols.Parallel
			opts.Targets = []string{"mock2"}
		})
		assert.NoError(t, err)
		assert.ElementsMatch(t, saved, []string{"mock2"})

		mock1.AssertNotCalled(t, "Save")
		mock2.AssertExpectations(t)
	})

	t.Run("it saves an item with unit error - sequential mode", func(t *testing.T) {
		setup()
		mock1 := test.NewUnitMock()
		mock2 := test.NewUnitMock()

		orchestrator.AddUnit("mock1", mock1)
		orchestrator.AddUnit("mock2", mock2)

		mock1.On("Save", "query", "saved", mock.Anything).Return(errors.New("unit1 error"))
		mock2.On("Save", "query", "saved", mock.Anything).Return(nil)

		saved, err := orchestrator.Save("query", "saved", func(opts *protocols.SaveOptions) {
			opts.HowWillItSave = protocols.Sequential
			opts.Targets = []string{
				"mock1",
				"mock2",
			}
		})

		assert.ErrorContains(t, err, "error saving unit mock1: unit1 error")

		assert.ElementsMatch(t, saved, []string{})

		mock1.AssertExpectations(t)
	})

	t.Run("it saves an item with unit error - parallel mode", func(t *testing.T) {
		setup()
		mock1 := test.NewUnitMock()
		mock2 := test.NewUnitMock()

		orchestrator.AddUnit("mock1", mock1)
		orchestrator.AddUnit("mock2", mock2)
		orchestrator.SetStandardOrder("mock1", "mock2")

		mock1.On("Save", "query", "saved", mock.Anything).Return(errors.New("unit1 error"))
		mock2.On("Save", "query", "saved", mock.Anything).Return(nil)

		saved, err := orchestrator.Save("query", "saved", func(opts *protocols.SaveOptions) {
			opts.HowWillItSave = protocols.Parallel
		})

		assert.ErrorContains(t, err, "error saving unit mock1: unit1 error")

		assert.ElementsMatch(t, saved, []string{"mock2"})

		mock1.AssertExpectations(t)
	})

	t.Run("it saves an item with context error - sequential mode", func(t *testing.T) {
		setup()
		mock1 := test.NewUnitMock()
		mock2 := test.NewUnitMock()

		orchestrator.AddUnit("mock1", mock1)
		orchestrator.AddUnit("mock2", mock2)
		orchestrator.SetStandardOrder("mock1", "mock2")

		saved, err := orchestrator.Save("query", "saved", func(opts *protocols.SaveOptions) {
			ctx, cancel := context.WithCancel(context.Background())
			opts.Context = ctx
			cancel()
		})

		assert.ErrorIs(t, err, context.Canceled)

		assert.ElementsMatch(t, saved, []string{})

	})

	t.Run("it saves an item with context error - parallel mode", func(t *testing.T) {
		setup()
		mock1 := test.NewUnitMock()
		mock2 := test.NewUnitMock()

		orchestrator.AddUnit("mock1", mock1)
		orchestrator.AddUnit("mock2", mock2)
		orchestrator.SetStandardOrder("mock1", "mock2")

		saved, err := orchestrator.Save("query", "saved", func(opts *protocols.SaveOptions) {
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

	t.Run("it receives an error when trying to get a Cache item without passing order", func(t *testing.T) {
		setup()
		mock1 := test.NewUnitMock()
		mock2 := test.NewUnitMock()

		orchestrator.AddUnit("mock1", mock1)
		orchestrator.AddUnit("mock2", mock2)

		_, err := orchestrator.Get("caughtTest")

		assert.Errorf(t, err, "unspecified order")

	})

	t.Run("it receives a valid value from the first unit when passed order and is of cache type", func(t *testing.T) {
		setup()
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

		mock1.AssertExpectations(t)
	})

	t.Run("it receives a valid value from the second unit when passed order and is of cache type", func(t *testing.T) {
		setup()
		mock1 := test.NewUnitMock()
		mock2 := test.NewUnitMock()

		orchestrator.AddUnit("mock1", mock1)
		orchestrator.AddUnit("mock2", mock2)

		expectedValue := "worked"
		mock1.On("Save", "caughtTest", "worked", mock.Anything).Return(nil)
		mock1.On("Get", "caughtTest", mock.Anything).Return("", fmt.Errorf("value not found"))

		mock2.On("Get", "caughtTest", mock.Anything).Return(expectedValue, nil)

		value, err := orchestrator.Get("caughtTest", func(opts *protocols.GetOptions) {
			opts.Targets = []string{"mock1", "mock2"}
		})

		assert.Equal(t, expectedValue, value)
		assert.NoError(t, err)

		mock1.AssertExpectations(t)
		mock2.AssertExpectations(t)
	})

	t.Run("it doesn't execute the save method of the unit when none has data", func(t *testing.T) {
		setup()
		mock1 := test.NewUnitMock()
		mock2 := test.NewUnitMock()

		orchestrator.AddUnit("mock1", mock1)
		orchestrator.AddUnit("mock2", mock2)

		expectedValue := ""
		expectedErr := "no unit returned"

		mock1.On("Get", "caughtTest", mock.Anything).Return("", fmt.Errorf("value not found"))

		mock2.On("Get", "caughtTest", mock.Anything).Return("", fmt.Errorf("value not found"))

		value, err := orchestrator.Get("caughtTest", func(opts *protocols.GetOptions) {
			opts.Targets = []string{"mock1", "mock2"}
		})

		assert.Equal(t, expectedValue, value)
		assert.ErrorContains(t, err, expectedErr)

		mock1.AssertExpectations(t)
		mock2.AssertExpectations(t)
	})

	t.Run("should return an unit error when it doesn't have it in the first caches and it gives an error when saving", func(t *testing.T) {
		setup()
		mock1 := test.NewUnitMock()
		mock2 := test.NewUnitMock()

		orchestrator.AddUnit("mock1", mock1)
		orchestrator.AddUnit("mock2", mock2)

		mock1.On("Get", "caughtTest", mock.Anything).Return("", fmt.Errorf("value not found"))
		mock1.On("Save", "caughtTest", "value", mock.Anything).Return(fmt.Errorf("didnt'save mock1"))

		mock2.On("Get", "caughtTest", mock.Anything).Return("value", nil)

		value, err := orchestrator.Get("caughtTest", func(opts *protocols.GetOptions) {
			opts.Targets = []string{"mock1", "mock2"}
		})

		assert.Equal(t, "value", value)
		assert.ErrorContains(t, err, "didnt'save mock1")

		mock1.AssertExpectations(t)
		mock2.AssertExpectations(t)
	})
}

func TestOrchestratorDelete(t *testing.T) {

	t.Run("should reach all storage units when none returns error", func(t *testing.T) {
		setup()
		mock1 := test.NewUnitMock()
		mock2 := test.NewUnitMock()

		orchestrator.AddUnit("mock1", mock1)
		orchestrator.AddUnit("mock2", mock2)

		mock1.On("Delete", "query", mock.Anything).Return(nil)
		mock2.On("Delete", "query", mock.Anything).Return(nil)

		err := orchestrator.Delete("query", func(opt *protocols.DeleteOptions) {
			opt.Targets = []string{
				"mock1",
				"mock2",
			}
		})

		assert.NoError(t, err)
		mock1.AssertNumberOfCalls(t, "Delete", 1)
		mock2.AssertNumberOfCalls(t, "Delete", 1)

		mock1.AssertExpectations(t)
		mock2.AssertExpectations(t)

	})
}
