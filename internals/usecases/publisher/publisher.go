package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"github-com/edarha/uploadfile-test/internals/usecases/entities"
	"github-com/edarha/uploadfile-test/internals/util"
	"time"

	"go.uber.org/zap"
)

type Publisher struct {
	Natsjs util.NatsJetstream
	Logger *zap.Logger
}

func (p *Publisher) UploadPublish(msgData *entities.MsgData) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	data, err := json.Marshal(msgData)
	if err != nil {
		return fmt.Errorf("UploadPublish: cannot marshal msgData, error: %s", err.Error())
	}
	id, err := p.Natsjs.PublishAsyncContext(ctx, util.UploadSubject, data)
	if err != nil {
		return fmt.Errorf("UploadPublish: Publish msg occur error, error: %s", err.Error())
	}

	p.Logger.Info("UploadPublish: Publish msg success", zap.String("msg_id", id))
	return nil
}
