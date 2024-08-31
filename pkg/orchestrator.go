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
		saved, err = o.saveInParallel(opt.Context, query, item, opt.Targets)
	case opt.HowWillItSave == protocols.Sequential:
		saved, err = o.saveInSequence(opt.Context, query, item, opt.Targets)
	}

	return saved, err

}

func (o *Orchestrator[K, V]) saveInParallel(ctx context.Context, query K, item V, targets []string) ([]string, error) {
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

			if err := unit.Save(ctx, query, item); err != nil {
				cancel()
				err = fmt.Errorf("error saving unit %v: %v", key, err.Error())
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

func (o *Orchestrator[K, V]) saveInSequence(ctx context.Context, query K, item V, targets []string) ([]string, error) {
	saved := make([]string, 0, len(o.units))
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, key := range targets {

		if ctx.Err() != nil {
			return saved, ctx.Err()
		}
		unit := o.units[key]
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
		object, err = o.getInCache(opt.Context, query, opt.Targets)
	case opt.HowWillItGet == protocols.Race:
	}

	return object, err
}

func (o *Orchestrator[K, V]) Delete(query K, opt ...protocols.OptionsFunc) error {
	return nil
}

func (o *Orchestrator[K, V]) getInCache(ctx context.Context, query K, orders []string) (value V, returnErr error) {
	notExistIn := make([]string, 0)

	defer func() {
		returnErr = o.addMissingElements(query, value, orders, notExistIn)
	}()

	for _, order := range orders {
		value, returnErr = o.units[order].Get(ctx, query)

		if returnErr == nil {
			break
		}
		notExistIn = append(notExistIn, order)
	}

	return value, returnErr
}

func (o *Orchestrator[K, V]) addMissingElements(query K, value V, orders []string, missing []string) error {
	if len(missing) == len(orders) {
		err := fmt.Errorf("no unit returned")
		return err
	}

	if len(missing) > 0 {
		_, err := o.Save(query, value, func(so *protocols.SaveOptions) {
			so.Targets = missing
		})
		if err != nil {
			return fmt.Errorf("err saving to units: %v ", err.Error())

		}
	}
	return nil
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

var _ protocols.StorageOrchestrator[any, any] = (*Orchestrator[any, any])(nil)
