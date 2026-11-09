package testFileTransfering

import (
	"Kaban/internal/DomainLevel"
	"context"
)

type testDelete struct {
	h *[]byte
}

func (t testDelete) Deleter() error {
	return nil
}

func (t testDelete) SetName(s string) DomainLevel.MakerDeleter {
	return t
}

func (t testDelete) SetSize(i int64) DomainLevel.MakerDeleter {
	return t
}

func (t testDelete) Make(ctx context.Context) (DomainLevel.Delete, error) {
	return t, nil
}

type DeleteError struct {
	err error
}

func NewDeleteError(err error) DeleteError {
	return DeleteError{err: err}
}

func (d DeleteError) Deleter() error {
	return d.err
}

func (d DeleteError) SetName(s string) DomainLevel.MakerDeleter {
	return d
}

func (d DeleteError) SetSize(i int64) DomainLevel.MakerDeleter {
	return d
}

func (d DeleteError) Make(ctx context.Context) (DomainLevel.Delete, error) {
	return d, nil
}
