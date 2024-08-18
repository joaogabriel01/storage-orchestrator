package protocols

import "context"

type TypeGetOptions uint
type TypeSaveOptions uint

const (
	Sequential TypeSaveOptions = iota
	Parallel
)

const (
	Cache TypeGetOptions = iota
	Race
)

type StorageOrchestrator[K any, V any] interface {
	Save(item V, opt ...SaveOptionsFunc) ([]string, error)
	Get(query K, opt ...OptionsFunc) (V, error)
	Delete(query K, opt ...OptionsFunc) error
	Sync(from string, to []string, opt ...OptionsFunc) error

	AddUnit(storage StorageUnit[K, V], storageName string) error
	GetUnits() (map[string]StorageUnit[K, V], error)
	GetUnit(string) (StorageUnit[K, V], error)
}

type SaveOptionsFunc func(*SaveOptions)

type OptionsFunc func(*Options)

type SaveOptions struct {
	Context       context.Context
	HowWillItSave TypeSaveOptions
}

type GetOptions struct {
	Context      context.Context
	HowWillItGet TypeGetOptions
}

type Options struct {
	Context context.Context
}

type StorageUnit[K any, V any] interface {
	Save(item V, ctx context.Context) error
	Get(query K, ctx context.Context) (V, error)
	Delete(query K, ctx context.Context) error
}
