package publisher

import (
	"fmt"
	"github-com/edarha/uploadfile-test/internals/logs"
	"github-com/edarha/uploadfile-test/internals/usecases/entities"
	"github-com/edarha/uploadfile-test/internals/usecases/mocks"
	"github-com/edarha/uploadfile-test/internals/util"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUploadPublish(t *testing.T) {
	logger := logs.NewZapLogger("test")
	msgData := &entities.MsgData{
		FileName: "fileName",
		Path:     "path",
	}

	t.Run("Success: publish msg to nats success", func(t *testing.T) {
		mNats := new(mocks.NatsJetstream)
		mNats.On("PublishAsyncContext", mock.Anything, util.UploadSubject, mock.Anything).Return("12345", nil)
		p := Publisher{
			Natsjs: mNats,
			Logger: logger,
		}
		err := p.UploadPublish(msgData)
		assert.NoError(t, err)
	})

	t.Run("Fail: publish msg to nats fail", func(t *testing.T) {
		mNats := new(mocks.NatsJetstream)
		mNats.On("PublishAsyncContext", mock.Anything, util.UploadSubject, mock.Anything).Return("", fmt.Errorf("test"))
		p := Publisher{
			Natsjs: mNats,
			Logger: logger,
		}
		err := p.UploadPublish(msgData)
		assert.Error(t, err)
		assert.Equal(t, "UploadPublish: Publish msg occur error, error: test", err.Error())
	})

}
