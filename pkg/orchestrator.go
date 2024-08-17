package pkg

import (
	"fmt"

	"github.com/joaogabriel01/storage-orchestrator/pkg/protocols"
)

type Orchestrator[K any, V any] struct {
	units            map[string]protocols.StorageUnit[K, V]
	typeOrchestrator uint
}

func (o *Orchestrator[K, V]) AddUnit(storageName string, storage protocols.StorageUnit[K, V]) error {
	o.units[storageName] = storage
	return nil
}

func (o *Orchestrator[K, V]) GetUnits() (map[string]protocols.StorageUnit[K, V], error) {
	return o.units, nil
}

func (o *Orchestrator[K, V]) GetUnit(unitName string) (protocols.StorageUnit[K, V], error) {
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
