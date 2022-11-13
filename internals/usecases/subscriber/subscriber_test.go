package subscriber

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github-com/edarha/uploadfile-test/internals/logs"
	"github-com/edarha/uploadfile-test/internals/usecases/entities"
	"github-com/edarha/uploadfile-test/internals/usecases/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TODO: check unit test
func TestUploadHandler(t *testing.T) {
	fileName := "fileName.json"
	tmpDir := os.TempDir()
	path := filepath.Join(tmpDir, fileName)
	_, err := os.Create(path)
	assert.NoError(t, err)

	msg := entities.MsgData{
		FileName: fileName,
		Path:     path,
	}
	data, err := json.Marshal(msg)
	assert.NoError(t, err)

	t.Run("Success: process msg success", func(t *testing.T) {
		mFileStore := &mocks.FileStore{}
		s := Subscriber{
			FileStore: mFileStore,
			Logger:    logs.NewZapLogger("test"),
		}
		mFileStore.On("UploadFile", mock.Anything, msg.FileName, msg.Path).Return(nil)
		err = s.UploadHandler(context.Background(), data)

		assert.NoError(t, err)
	})

	t.Run("Fail: process msg occur error", func(t *testing.T) {
		mFileStore := &mocks.FileStore{}
		s := Subscriber{
			FileStore: mFileStore,
		}
		mFileStore.On("UploadFile", mock.Anything, msg.FileName, msg.Path).Return(fmt.Errorf("cannot upload file"))
		err = s.UploadHandler(context.Background(), data)

		assert.Error(t, err)
		assert.Equal(t, "UploadHandler: UploadFile occur error, fileName: fileName.json, err: cannot upload file", err.Error())
	})

}
