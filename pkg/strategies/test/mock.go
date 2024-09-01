package strategies_mock

import (
	"context"

	"github.com/joaogabriel01/storage-orchestrator/pkg/protocols"
	"github.com/stretchr/testify/mock"
)

type MockSaveStrategy struct {
	mock.Mock
}

func (u *MockSaveStrategy) Save(ctx context.Context, query string, item string, units map[string]protocols.StorageUnit[string, string], targets []string, auxiliary ...any) ([]string, error) {
	args := u.Called(ctx, query, item, units, targets, auxiliary)
	return args.Get(0).([]string), args.Error(1)
}

var _ protocols.SaveStrategy[string, string] = (*MockSaveStrategy)(nil)

type MockGetStrategy struct {
	mock.Mock
}

func (m *MockGetStrategy) Get(ctx context.Context, query string, units map[string]protocols.StorageUnit[string, string], targets []string, auxiliary ...any) (string, error) {
	args := m.Called(ctx, query, units, targets, auxiliary)
	return args.Get(0).(string), args.Error(1)
}

var _ protocols.GetStrategy[string, string] = (*MockGetStrategy)(nil)

type MockDeleteStrategy struct {
	mock.Mock
}

func (m *MockDeleteStrategy) Delete(ctx context.Context, query string, units map[string]protocols.StorageUnit[string, string], targets []string, _ ...any) error {
	args := m.Called(ctx, query, units, targets)
	return args.Error(0)
}

var _ protocols.DeleteStrategy[string, string] = (*MockDeleteStrategy)(nil)
