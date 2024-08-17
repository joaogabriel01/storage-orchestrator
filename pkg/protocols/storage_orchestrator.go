package protocols

import "context"

type Options interface {
	GetTypeOrchestrator() uint
	GetContext() context.Context
}

type StorageOrchestrator[K any, V any] interface {
	Save(item V, opt ...Options) error
	Get(query K, opt ...Options) (V, error)
	Delete(query K, opt ...Options) error
	Sync(from string, to []string, opt ...Options) error
}

type StorageUnit[K any, V any] interface {
	Save(item V, ctx context.Context) error
	Get(query K, ctx context.Context) (V, error)
	Delete(query K, ctx context.Context) error
}
