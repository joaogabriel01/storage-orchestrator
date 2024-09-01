package pkg

import (
	"fmt"
	"sync"

	"github.com/joaogabriel01/storage-orchestrator/pkg/protocols"
	"github.com/joaogabriel01/storage-orchestrator/pkg/strategies"
)

type Orchestrator[K any, V any] struct {
	mu               sync.RWMutex
	units            map[string]protocols.StorageUnit[K, V]
	standardOrder    []string
	saveStrategies   []protocols.SaveStrategy[K, V]
	getStrategies    []protocols.GetStrategy[K, V]
	deleteStrategies []protocols.DeleteStrategy[K, V]
}

func (o *Orchestrator[K, V]) Save(query K, item V, opts ...protocols.SaveOptionsFunc) ([]string, error) {
	opt := o.defaultSaveOptions()
	for _, fn := range opts {
		fn(&opt)
	}

	return o.saveStrategies[opt.HowWillItSave].Save(opt.Context, query, item, o.units, opt.Targets)

}

func (o *Orchestrator[K, V]) Get(query K, opts ...protocols.GetOptionsFunc) (V, error) {
	opt := o.defaultGetOptions()
	for _, fn := range opts {
		fn(&opt)
	}

	return o.getStrategies[opt.HowWillItGet].Get(opt.Context, query, o.units, opt.Targets, o.saveStrategies[protocols.Sequential])
}

func (o *Orchestrator[K, V]) Delete(query K, opts ...protocols.DeleteOptionsFunc) error {
	opt := o.defaultDeleteOptions()

	for _, fn := range opts {
		fn(&opt)
	}

	return o.deleteStrategies[opt.HowWillItDelete].Delete(opt.Context, query, o.units, opt.Targets)
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

func (o *Orchestrator[K, V]) SetStandardOrder(targets ...string) error {
	order := make([]string, len(targets))
	for c, target := range targets {
		_, ok := o.units[target]
		if !ok {
			return fmt.Errorf("this unit does not exist: %v", target)
		}
		order[c] = target
	}
	o.standardOrder = order
	return nil
}

func NewOrchestrator[K any, V any](units map[string]protocols.StorageUnit[K, V], standardOrder []string) Orchestrator[K, V] {
	var saveStragies []protocols.SaveStrategy[K, V]
	var getStrategies []protocols.GetStrategy[K, V]
	var deleteStrategies []protocols.DeleteStrategy[K, V]

	sequentialSave := strategies.SequentialSaveStrategy[K, V]{}
	parallelSave := strategies.ParallelSaveStrategy[K, V]{}
	saveStragies = append(saveStragies, &sequentialSave, &parallelSave)

	cacheGet := strategies.CacheGetStrategy[K, V]{}
	getStrategies = append(getStrategies, &cacheGet)

	sequentialDelete := strategies.SequentialDeleteStrategy[K, V]{}
	deleteStrategies = append(deleteStrategies, &sequentialDelete)
	return Orchestrator[K, V]{
		units:            units,
		standardOrder:    standardOrder,
		saveStrategies:   saveStragies,
		getStrategies:    getStrategies,
		deleteStrategies: deleteStrategies,
	}
}

var _ protocols.StorageOrchestrator[any, any] = (*Orchestrator[any, any])(nil)
