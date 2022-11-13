package subscriber

import (
	"context"
)

type FileStore interface {
	UploadFile(ctx context.Context, fileName, path string) error
	Close() error
}
