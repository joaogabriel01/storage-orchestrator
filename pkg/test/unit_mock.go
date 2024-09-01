package unit_test

import (
	"context"

	"github.com/joaogabriel01/storage-orchestrator/pkg/protocols"
	"github.com/stretchr/testify/mock"
)

type UnitMock struct {
	mock.Mock
}

func NewUnitMock() *UnitMock {
	return &UnitMock{}
}

func (u *UnitMock) Save(ctx context.Context, query string, item string) error {
	args := u.Called(query, item, ctx)
	return args.Error(0)
}

func (u *UnitMock) Get(ctx context.Context, query string) (string, error) {
	args := u.Called(query, ctx)
	return args.String(0), args.Error(1)
}

func (u *UnitMock) Delete(ctx context.Context, query string) error {
	args := u.Called(query, ctx)
	return args.Error(0)
}

var _ protocols.StorageUnit[string, string] = (*UnitMock)(nil)
