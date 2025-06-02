package object_storage

import (
	"context"
	"io"
)

type ObjectStorage interface {
	GetObject(ctx context.Context, path string) (*Object, error)
	GetObjectBody(ctx context.Context, obj *Object) (io.ReadCloser, error)
}
