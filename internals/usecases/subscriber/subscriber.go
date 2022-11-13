package subscriber

import (
	"context"
	"encoding/json"
	"fmt"
	"github-com/edarha/uploadfile-test/internals/usecases/entities"
	"github-com/edarha/uploadfile-test/internals/util"
	"os"
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type Subscriber struct {
	Natsjs    util.NatsJetstream
	Logger    *zap.Logger
	FileStore FileStore
}

// Create consumer in nats
// It is used to receive msg from nats than processing the msg base on our logic.
func (s *Subscriber) UploadSubscription() error {
	err := s.Natsjs.QueueSubscribe(util.UploadSubject, util.UploadQueue, s.UploadHandler,
		nats.ManualAck(),
		nats.MaxDeliver(10),
		nats.AckWait(30*time.Second))

	if err != nil {
		return fmt.Errorf(fmt.Sprintf("UploadSubscription: create consumer occur error, err: %s", err.Error()))
	}

	s.Logger.Info("Create UploadSubscription successfully")
	return nil
}

// This function to process the message from nats.
// It will call gcp or s3 to upload data.
// Then remove that file.
func (s *Subscriber) UploadHandler(ctx context.Context, data []byte) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	var msgData entities.MsgData
	err := json.Unmarshal(data, &msgData)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("UploadHandler: cannot unmarshal data, err: %s", err.Error()))
	}

	err = s.FileStore.UploadFile(ctx, msgData.FileName, msgData.Path)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("UploadHandler: UploadFile occur error, fileName: %s, err: %s", msgData.FileName, err.Error()))
	}

	err = os.Remove(msgData.Path)
	if err != nil {
		s.Logger.Error("Remove file failed", zap.String("filePath", msgData.Path))
	}
	return nil
}
