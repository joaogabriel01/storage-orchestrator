package strategies

import (
	"context"
	"fmt"
	"sync"

	"github.com/joaogabriel01/storage-orchestrator/pkg/protocols"
)

type SequentialSaveStrategy[K any, V any] struct{}

func (s *SequentialSaveStrategy[K, V]) Save(ctx context.Context, query K, item V, units map[string]protocols.StorageUnit[K, V], targets []string, _ ...any) ([]string, error) {
	saved := make([]string, 0, len(units))
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, key := range targets {

		if ctx.Err() != nil {
			return saved, ctx.Err()
		}
		unit := units[key]
		err := unit.Save(ctx, query, item)
		if err != nil {
			cancel()
			err = fmt.Errorf("error saving unit %v: %v", key, err.Error())
			return saved, err
		}
		saved = append(saved, key)
	}
	return saved, nil
}

type ParallelSaveStrategy[K any, V any] struct{}

func (p *ParallelSaveStrategy[K, V]) Save(ctx context.Context, query K, item V, units map[string]protocols.StorageUnit[K, V], targets []string, _ ...any) ([]string, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	mu := sync.Mutex{}
	saved := make([]string, 0, len(units))
	errCh := make(chan error, len(units))

	for _, key := range targets {
		wg.Add(1)

		go func(key string, unit protocols.StorageUnit[K, V]) {
			defer wg.Done()

			if ctx.Err() != nil {
				errCh <- fmt.Errorf("context finalized: %w", ctx.Err())
				return
			}

			if err := unit.Save(ctx, query, item); err != nil {
				cancel()
				err = fmt.Errorf("error saving unit %v: %v", key, err.Error())
				errCh <- err
				return
			}

			mu.Lock()
			saved = append(saved, key)
			mu.Unlock()
		}(key, units[key])
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			return saved, err
		}
	}

	return saved, nil
}
