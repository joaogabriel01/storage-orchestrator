package protocols

import "context"

type TypeGetOptions uint
type TypeSaveOptions uint
type TypeDeleteOptions uint

const (
	Sequential TypeSaveOptions = iota
	Parallel
)

const (
	SequentialDelete TypeDeleteOptions = iota
)

const (
	Cache TypeGetOptions = iota
	Race
)

type StorageOrchestrator[K any, V any] interface {
	Save(query K, item V, opt ...SaveOptionsFunc) ([]string, error)
	Get(query K, opt ...GetOptionsFunc) (V, error)
	Delete(query K, opt ...DeleteOptionsFunc) error

	AddUnit(storageName string, storage StorageUnit[K, V]) error
	GetUnits() (map[string]StorageUnit[K, V], error)
	GetUnit(string) (StorageUnit[K, V], error)

	SetStandardOrder(targets ...string) error
}

type SaveOptionsFunc func(*SaveOptions)

type GetOptionsFunc func(*GetOptions)

type DeleteOptionsFunc func(*DeleteOptions)

type SaveOptions struct {
	Context       context.Context
	HowWillItSave TypeSaveOptions
	Targets       []string
}

type GetOptions struct {
	Context      context.Context
	HowWillItGet TypeGetOptions
	Targets      []string
}

type DeleteOptions struct {
	Context         context.Context
	Targets         []string
	HowWillItDelete TypeDeleteOptions
}

type StorageUnit[K any, V any] interface {
	Save(ctx context.Context, query K, item V) error
	Get(ctx context.Context, query K) (V, error)
	Delete(ctx context.Context, query K) error
}

type SaveStrategy[K any, V any] interface {
	Save(ctx context.Context, query K, item V, units map[string]StorageUnit[K, V], targets []string, auxiliary ...any) ([]string, error)
}

type GetStrategy[K any, V any] interface {
	Get(ctx context.Context, query K, units map[string]StorageUnit[K, V], targets []string, auxiliary ...any) (V, error)
}

type DeleteStrategy[K any, V any] interface {
	Delete(ctx context.Context, query K, units map[string]StorageUnit[K, V], targets []string, auxiliary ...any) error
}
