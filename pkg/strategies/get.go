package strategies

import (
	"context"
	"fmt"

	"github.com/joaogabriel01/storage-orchestrator/pkg/protocols"
)

type CacheGetStrategy[K any, V any] struct{}

func (c *CacheGetStrategy[K, V]) Get(ctx context.Context, query K, units map[string]protocols.StorageUnit[K, V], targets []string, auxiliary ...any) (value V, returnErr error) {
	var notExistIn []string

	if len(auxiliary) != 1 {
		return value, fmt.Errorf("save function not found")
	}

	saveFunction, ok := auxiliary[0].(protocols.SaveStrategy[K, V])
	if !ok {
		return value, fmt.Errorf("save function check did not work")
	}

	defer func() {
		returnErr = c.addMissingElements(ctx, query, value, targets, units, notExistIn, saveFunction)
	}()

	for _, target := range targets {
		unit := units[target]
		value, err := unit.Get(ctx, query)
		if err == nil {
			return value, nil
		}
		notExistIn = append(notExistIn, target)
	}

	return value, returnErr
}

func (c *CacheGetStrategy[K, V]) addMissingElements(ctx context.Context, query K, value V, orders []string, units map[string]protocols.StorageUnit[K, V], missing []string, saveFunction protocols.SaveStrategy[K, V]) error {
	if len(missing) == len(orders) {
		err := fmt.Errorf("no unit returned")
		return err
	}

	if len(missing) > 0 {
		_, err := saveFunction.Save(ctx, query, value, units, missing)
		if err != nil {
			return fmt.Errorf("err saving to units: %v ", err.Error())

		}
	}
	return nil
}

var _ protocols.GetStrategy[any, any] = (*CacheGetStrategy[any, any])(nil)
