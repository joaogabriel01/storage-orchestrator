package protocols

import "context"

type typeOrchestrator uint

const (
	Cache typeOrchestrator = iota
	Race
)

type StorageOrchestrator[K any, V any] interface {
	Save(item V, opt ...OptionsFunc) ([]string, error)
	Get(query K, opt ...OptionsFunc) (V, error)
	Delete(query K, opt ...OptionsFunc) error
	Sync(from string, to []string, opt ...OptionsFunc) error

	AddUnit(storage StorageUnit[K, V], storageName string) error
	GetUnits() (map[string]StorageUnit[K, V], error)
	GetUnit(string) (StorageUnit[K, V], error)
}

type OptionsFunc func(*Options)

type Options struct {
	Context context.Context
}

type StorageUnit[K any, V any] interface {
	Save(item V, ctx context.Context) error
	Get(query K, ctx context.Context) (V, error)
	Delete(query K, ctx context.Context) error
}
