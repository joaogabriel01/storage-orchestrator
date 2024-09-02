package pkg

import (
	"testing"

	"github.com/joaogabriel01/storage-orchestrator/pkg/protocols"
	strategies_mock "github.com/joaogabriel01/storage-orchestrator/pkg/strategies/test"
	unit_test "github.com/joaogabriel01/storage-orchestrator/pkg/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var orchestrator Orchestrator[string, string]
var mock1 *unit_test.UnitMock
var mock2 *unit_test.UnitMock
var saveStrategy strategies_mock.MockSaveStrategy
var getStrategy strategies_mock.MockGetStrategy
var deleteStrategy strategies_mock.MockDeleteStrategy

func setupOrchestrator() {
	units := make(map[string]protocols.StorageUnit[string, string])
	mock1 = unit_test.NewUnitMock()
	mock2 = unit_test.NewUnitMock()

	standardOrder := []string{"mock1", "mock2"}
	saveStrategy = strategies_mock.MockSaveStrategy{}
	getStrategy = strategies_mock.MockGetStrategy{}
	deleteStrategy = strategies_mock.MockDeleteStrategy{}

	var saveStrategiesTyped []protocols.SaveStrategy[string, string]
	var getStrategiesTyped []protocols.GetStrategy[string, string]
	var deleteStrategiesTyped []protocols.DeleteStrategy[string, string]

	saveStrategiesTyped = append(saveStrategiesTyped, &saveStrategy)
	getStrategiesTyped = append(getStrategiesTyped, &getStrategy)
	deleteStrategiesTyped = append(deleteStrategiesTyped, &deleteStrategy)

	orchestrator = NewOrchestratorWithParameters[string, string](units, standardOrder, saveStrategiesTyped, getStrategiesTyped, deleteStrategiesTyped)

	orchestrator.AddUnit("mock1", mock1)
	orchestrator.AddUnit("mock2", mock2)
}

func TestOrchestratorUnitOperations(t *testing.T) {

	t.Run("should add and get units without error", func(t *testing.T) {
		setupOrchestrator()

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
		setupOrchestrator()
		err := orchestrator.SetStandardOrder("mock2")
		assert.NoError(t, err)
	})

	t.Run("should return error when there are no units", func(t *testing.T) {
		setupOrchestrator()
		err := orchestrator.SetStandardOrder("mock1", "mock2", "mock3")
		assert.ErrorContains(t, err, "this unit does not exist: mock3")

	})
}

func TestOrchestratorStrategies(t *testing.T) {

	t.Run("should call save strategy without opt func", func(t *testing.T) {
		setupOrchestrator()
		saveStrategy.On("Save", mock.Anything, "query", "value", orchestrator.units, []string{"mock1", "mock2"}, mock.Anything).Return([]string{"mock1", "mock2"}, nil)
		saved, err := orchestrator.Save("query", "value")
		assert.NoError(t, err)
		assert.Equal(t, []string{"mock1", "mock2"}, saved)
		saveStrategy.AssertExpectations(t)

	})
	t.Run("should call save strategy with opt func", func(t *testing.T) {
		setupOrchestrator()
		saveStrategy.On("Save", mock.Anything, "query", "value", orchestrator.units, []string{"mock1"}, mock.Anything).Return([]string{"mock1"}, nil)
		saved, err := orchestrator.Save("query", "value", func(opt *protocols.SaveOptions) {
			opt.Targets = []string{"mock1"}
		})
		assert.NoError(t, err)
		assert.Equal(t, []string{"mock1"}, saved)
		saveStrategy.AssertExpectations(t)

	})

	t.Run("should call get strategy without opt func", func(t *testing.T) {
		setupOrchestrator()
		getStrategy.On("Get", mock.Anything, "query", orchestrator.units, []string{"mock1", "mock2"}, mock.Anything).Return("value", nil)
		value, err := orchestrator.Get("query")
		assert.NoError(t, err)
		assert.Equal(t, "value", value)
		getStrategy.AssertExpectations(t)
	})

	t.Run("should call get strategy with opt func", func(t *testing.T) {
		setupOrchestrator()
		getStrategy.On("Get", mock.Anything, "query", orchestrator.units, []string{"mock1"}, mock.Anything).Return("value", nil)
		value, err := orchestrator.Get("query", func(opt *protocols.GetOptions) {
			opt.Targets = []string{"mock1"}
		})
		assert.NoError(t, err)
		assert.Equal(t, "value", value)
		getStrategy.AssertExpectations(t)
	})

	t.Run("should call delete strategy without opt func", func(t *testing.T) {
		setupOrchestrator()

		deleteStrategy.On("Delete", mock.Anything, "query", orchestrator.units, []string{"mock1", "mock2"}, mock.Anything).Return(nil)
		err := orchestrator.Delete("query")
		assert.NoError(t, err)
		deleteStrategy.AssertExpectations(t)

	})

	t.Run("should call delete strategy with opt func", func(t *testing.T) {
		setupOrchestrator()

		deleteStrategy.On("Delete", mock.Anything, "query", orchestrator.units, []string{"mock1"}, mock.Anything).Return(nil)
		err := orchestrator.Delete("query", func(opt *protocols.DeleteOptions) {
			opt.Targets = []string{"mock1"}
		})
		assert.NoError(t, err)
		deleteStrategy.AssertExpectations(t)

	})
}
