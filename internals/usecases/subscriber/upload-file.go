package subscriber

import (
	"context"
)

type FileStore interface {
	UploadFile(ctx context.Context, fileName string) error
	Close() error
}
