package test

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

func (u *UnitMock) Save(item string, ctx context.Context) error {
	args := u.Called(item, ctx)
	return args.Error(0)
}

func (u *UnitMock) Get(query string, ctx context.Context) (string, error) {
	args := u.Called(query, ctx)
	return args.String(0), args.Error(1)
}

func (u *UnitMock) Delete(query string, ctx context.Context) error {
	args := u.Called(query, ctx)
	return args.Error(0)
}

var _ protocols.StorageUnit[string, string] = (*UnitMock)(nil)
