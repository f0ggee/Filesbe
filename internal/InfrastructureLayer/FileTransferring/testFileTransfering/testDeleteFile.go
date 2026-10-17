package testFileTransfering

import (
	"Kaban/internal/DomainLevel"
	"context"
)

type TestDelete struct {
}

func (t TestDelete) SetName(s string) DomainLevel.MakerDeleter {
	return t
}

func (t TestDelete) SetSize(i int64) DomainLevel.MakerDeleter {
	return t
}

func (t TestDelete) Make(ctx context.Context) (DomainLevel.Delete, error) {
	return t, nil
}

func NewTestDelete() *TestDelete {
	return &TestDelete{}
}

func (t TestDelete) Deleter() error {
	return nil
}
