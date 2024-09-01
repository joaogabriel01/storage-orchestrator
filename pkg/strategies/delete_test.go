package strategies

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var sequentialDeleteStrategy SequentialDeleteStrategy[string, string]

func deleteSequentialSetup() {
	sequentialDeleteStrategy = SequentialDeleteStrategy[string, string]{}
	initialSetup()
}

func TestSequentialDelete(t *testing.T) {

	t.Run("should reach all storage units when none returns error", func(t *testing.T) {
		deleteSequentialSetup()
		ctx := context.Background()

		mock1.On("Delete", "query", mock.Anything).Return(nil)
		mock2.On("Delete", "query", mock.Anything).Return(nil)

		err := sequentialDeleteStrategy.Delete(ctx, "query", units, targets)
		assert.NoError(t, err)

		mock1.AssertNumberOfCalls(t, "Delete", 1)
		mock2.AssertNumberOfCalls(t, "Delete", 1)

		mock1.AssertExpectations(t)
		mock2.AssertExpectations(t)

	})

	t.Run("should return error when one unit return errors too", func(t *testing.T) {
		deleteSequentialSetup()
		ctx := context.Background()
		mock1.On("Delete", "query", mock.Anything).Return(fmt.Errorf("mock1 error"))

		err := sequentialDeleteStrategy.Delete(ctx, "query", units, targets)
		assert.ErrorContains(t, err, "mock1 error")

		mock1.AssertNumberOfCalls(t, "Delete", 1)
		mock2.AssertNumberOfCalls(t, "Delete", 0)

		mock1.AssertExpectations(t)
		mock2.AssertExpectations(t)
	})
}
