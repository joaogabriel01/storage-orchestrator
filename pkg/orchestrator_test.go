package pkg

import (
	"testing"

	"github.com/joaogabriel01/storage-orchestrator/pkg/protocols"
	"github.com/joaogabriel01/storage-orchestrator/pkg/test"
	"github.com/stretchr/testify/assert"
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
	t.Run("it saves an item", func(t *testing.T) {
	})
}
