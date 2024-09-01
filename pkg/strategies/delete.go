package strategies

import (
	"context"
	"fmt"

	"github.com/joaogabriel01/storage-orchestrator/pkg/protocols"
)

type SequentialDeleteStrategy[K any, V any] struct{}

func (s *SequentialDeleteStrategy[K, V]) Delete(ctx context.Context, query K, units map[string]protocols.StorageUnit[K, V], targets []string, _ ...any) error {
	for _, key := range targets {
		unit := units[key]
		err := unit.Delete(ctx, query)
		if err != nil {
			return fmt.Errorf("error deleting in unit %v: %v", key, err.Error())
		}
	}
	return nil
}

var _ protocols.DeleteStrategy[any, any] = (*SequentialDeleteStrategy[any, any])(nil)
