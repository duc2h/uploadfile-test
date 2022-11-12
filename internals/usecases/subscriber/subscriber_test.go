package subscriber

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github-com/edarha/uploadfile-test/internals/usecases/entities"
	"github-com/edarha/uploadfile-test/internals/usecases/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TODO: check unit test
func TestUploadHandler(t *testing.T) {
	mFileStore := &mocks.FileStore{}
	s := Subscriber{
		FileStore: mFileStore,
	}
	msg := entities.MsgData{
		FileName: "fileName",
	}
	data, err := json.Marshal(msg)
	assert.NoError(t, err)

	t.Run("Fail: process msg occur error", func(t *testing.T) {
		mFileStore.On("UploadFile", mock.Anything, msg.FileName).Return(fmt.Errorf("cannot upload file"))
		err = s.UploadHandler(context.Background(), data)

		assert.Error(t, err)
		assert.Equal(t, "UploadHandler: UploadFile occur error, err: cannot upload file", err.Error())
	})

	// t.Run("Success: process msg success", func(t *testing.T) {
	// 	mFileStore.On("UploadFile", mock.Anything, msg.ObjectName, msg.PathFile).Return(nil)
	// 	err = s.UploadHandler(context.Background(), data)

	// 	assert.NoError(t, err)
	// })
}
