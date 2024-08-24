package pkg

import (
	"context"
	"fmt"
	"sync"

	"github.com/joaogabriel01/storage-orchestrator/pkg/protocols"
)

type Orchestrator[K any, V any] struct {
	mu    sync.RWMutex
	units map[string]protocols.StorageUnit[K, V]
}

func (o *Orchestrator[K, V]) Save(query K, item V, opts ...protocols.SaveOptionsFunc) ([]string, error) {
	opt := o.defaultSaveOptions()

	for _, fn := range opts {
		fn(&opt)
	}

	var saved []string
	var err error

	// i dont't see other ways of insertion so I didn't use polymorphism
	switch {
	case opt.HowWillItSave == protocols.Parallel:
		saved, err = o.saveInParallel(query, item, opt.Targets, opt.Context)
	case opt.HowWillItSave == protocols.Sequential:
		saved, err = o.saveInSequence(query, item, opt.Targets, opt.Context)
	}

	return saved, err

}

func (o *Orchestrator[K, V]) saveInParallel(query K, item V, targets []string, ctx context.Context) ([]string, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	mu := sync.Mutex{}
	saved := make([]string, 0, len(o.units))
	errCh := make(chan error, len(o.units))

	for _, key := range targets {
		wg.Add(1)

		go func(key string, unit protocols.StorageUnit[K, V]) {
			defer wg.Done()

			if ctx.Err() != nil {
				errCh <- fmt.Errorf("context finalized: %w", ctx.Err())
				return
			}

			if err := unit.Save(query, item, ctx); err != nil {
				cancel()
				errCh <- err
				return
			}

			mu.Lock()
			saved = append(saved, key)
			mu.Unlock()
		}(key, o.units[key])
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

func (o *Orchestrator[K, V]) saveInSequence(query K, item V, targets []string, ctx context.Context) ([]string, error) {
	saved := make([]string, 0, len(o.units))
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, key := range targets {

		if ctx.Err() != nil {
			return saved, ctx.Err()
		}
		unit := o.units[key]
		err := unit.Save(query, item, ctx)
		if err != nil {
			cancel()
			return saved, err
		}
		saved = append(saved, key)
	}
	return saved, nil
}

func (o *Orchestrator[K, V]) Get(query K, opts ...protocols.GetOptionsFunc) (V, error) {
	opt := o.defaultGetOptions()
	var object V
	var err error
	for _, fn := range opts {
		fn(&opt)
	}

	switch {
	case opt.HowWillItGet == protocols.Cache:
		if len(opt.Targets) < 1 {
			return object, fmt.Errorf("unspecified order")
		}
		object, err = o.getInCache(query, opt.Targets, opt.Context)
	case opt.HowWillItGet == protocols.Race:
	}

	return object, err
}

func (o *Orchestrator[K, V]) getInCache(query K, orders []string, ctx context.Context) (value V, err error) {
	notExistIn := make([]string, 0)

	defer func() {
		if len(notExistIn) == len(orders) {
			err = fmt.Errorf("no unit returned")
			return
		}

		if len(notExistIn) > 0 {
			o.Save(query, value, func(so *protocols.SaveOptions) {
				so.Targets = notExistIn
			})
		}
		//TODO handle errors
	}()

	for _, order := range orders {
		value, err = o.units[order].Get(query, ctx)

		if err == nil {
			break
		}
		notExistIn = append(notExistIn, order)
	}

	return value, err
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

func (o *Orchestrator[K, V]) defaultSaveOptions() protocols.SaveOptions {
	ctx := context.Background()

	targets := make([]string, 0, len(o.units))
	for key := range o.units {
		targets = append(targets, key)
	}
	options := protocols.SaveOptions{
		Context:       ctx,
		HowWillItSave: protocols.Sequential,
		Targets:       targets,
	}
	return options
}

func (o *Orchestrator[K, V]) defaultGetOptions() protocols.GetOptions {
	ctx := context.Background()

	options := protocols.GetOptions{
		Context:      ctx,
		HowWillItGet: protocols.Cache,
		Targets:      []string{},
	}
	return options
}

func NewOrchestrator[K any, V any](units map[string]protocols.StorageUnit[K, V]) Orchestrator[K, V] {
	return Orchestrator[K, V]{
		units: units,
	}
}
