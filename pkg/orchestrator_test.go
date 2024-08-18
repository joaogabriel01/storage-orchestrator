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

func TestOrchestrator(t *testing.T) {
	orchestrator := NewOrchestrator[string, string](map[string]protocols.StorageUnit[string, string]{}, uint(protocols.Cache))

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
	t.Run("it saves an item with success", func(t *testing.T) {
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
	})

	t.Run("it saves an item with unit error", func(t *testing.T) {
		mock1 := test.NewUnitMock()
		mock2 := test.NewUnitMock()

		orchestrator.AddUnit("mock1", mock1)
		orchestrator.AddUnit("mock2", mock2)

		expectedErr := errors.New("unit1 error")
		mock1.On("Save", "saved", mock.Anything).Return(expectedErr)
		mock2.On("Save", "saved", mock.Anything).Return(nil)

		saved, err := orchestrator.Save("saved")
		assert.ErrorIs(t, err, expectedErr)

		assert.ElementsMatch(t, saved, []string{"mock2"})

		mock1.AssertExpectations(t)
	})

	t.Run("it saves an item with context error", func(t *testing.T) {
		mock1 := test.NewUnitMock()
		mock2 := test.NewUnitMock()

		orchestrator.AddUnit("mock1", mock1)
		orchestrator.AddUnit("mock2", mock2)

		saved, err := orchestrator.Save("saved", func(opts *protocols.Options) {
			ctx, cancel := context.WithCancel(context.Background())
			opts.Context = ctx
			cancel()
		})

		assert.ErrorIs(t, err, context.Canceled)

		assert.ElementsMatch(t, saved, []string{})

	})
}
