package pkg

import (
	"context"
	"fmt"
	"sync"

	"github.com/joaogabriel01/storage-orchestrator/pkg/protocols"
)

func defaultSaveOptions() protocols.SaveOptions {
	ctx := context.Background()
	options := protocols.SaveOptions{
		Context:       ctx,
		HowWillItSave: protocols.Sequential,
	}
	return options
}

type Orchestrator[K any, V any] struct {
	mu               sync.RWMutex
	units            map[string]protocols.StorageUnit[K, V]
	typeOrchestrator uint
}

func (o *Orchestrator[K, V]) Save(item V, opts ...protocols.SaveOptionsFunc) ([]string, error) {
	opt := defaultSaveOptions()

	for _, fn := range opts {
		fn(&opt)
	}

	var saved []string
	var err error

	// i dont't see other ways of insertion so I didn't use polymorphism
	switch {
	case opt.HowWillItSave == protocols.Parallel:
		saved, err = o.saveInParallel(item, opt.Context)
	case opt.HowWillItSave == protocols.Sequential:
		saved, err = o.saveInSequence(item, opt.Context)
	}

	return saved, err

}

func (o *Orchestrator[K, V]) saveInParallel(item V, ctx context.Context) ([]string, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	mu := sync.Mutex{}
	saved := make([]string, 0, len(o.units))
	errCh := make(chan error, len(o.units))

	for key, unit := range o.units {
		wg.Add(1)

		go func(key string, unit protocols.StorageUnit[K, V]) {
			defer wg.Done()

			if ctx.Err() != nil {
				errCh <- fmt.Errorf("context finalized: %w", ctx.Err())
				return
			}

			if err := unit.Save(item, ctx); err != nil {
				cancel()
				errCh <- err
				return
			}

			mu.Lock()
			saved = append(saved, key)
			mu.Unlock()
		}(key, unit)
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

func (o *Orchestrator[K, V]) saveInSequence(item V, ctx context.Context) ([]string, error) {
	saved := make([]string, 0, len(o.units))
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for key, unit := range o.units {

		if ctx.Err() != nil {
			return saved, ctx.Err()
		}

		err := unit.Save(item, ctx)
		if err != nil {
			cancel()
			return saved, err
		}
		saved = append(saved, key)
	}
	return saved, nil
}

func (o *Orchestrator[K, V]) AddUnit(storageName string, storage protocols.StorageUnit[K, V]) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.units[storageName] = storage
	return nil
}

func (o *Orchestrator[K, V]) GetUnits() (map[string]protocols.StorageUnit[K, V], error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	unitsCopy := make(map[string]protocols.StorageUnit[K, V])
	for k, v := range o.units {
		unitsCopy[k] = v
	}
	return unitsCopy, nil
}

func (o *Orchestrator[K, V]) GetUnit(unitName string) (protocols.StorageUnit[K, V], error) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	unit, exists := o.units[unitName]
	if !exists {
		return nil, fmt.Errorf("unit not found")
	}

	return unit, nil
}

func NewOrchestrator[K any, V any](units map[string]protocols.StorageUnit[K, V], typeOrchestrator uint) Orchestrator[K, V] {
	return Orchestrator[K, V]{
		units:            units,
		typeOrchestrator: typeOrchestrator,
	}
}
