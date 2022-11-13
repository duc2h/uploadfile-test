package main

import (
	"github-com/edarha/uploadfile-test/internals/logs"
	"github-com/edarha/uploadfile-test/internals/transport"
	"github-com/edarha/uploadfile-test/internals/usecases/mocks"
	"github-com/edarha/uploadfile-test/internals/usecases/publisher"
	"github-com/edarha/uploadfile-test/internals/usecases/subscriber"
	"github-com/edarha/uploadfile-test/internals/util"
	"time"

	limits "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func main() {
	logger := logs.NewZapLogger("development")
	r := gin.Default()
	// TODO: get config
	// init nats
	natsjs, err := util.ConnectNats(logger, util.NatsConf{
		Url:           "nats://nats:4223",
		UserName:      "admin",
		Password:      "admin",
		MaxReconnect:  10,
		ReconnectWait: time.Second,
	})

	if err != nil {
		logger.Panic("Cannot connect to nats", zap.Error(err))
	}

	// create stream
	err = natsjs.CreateStream(util.UploadStream, util.UploadSubject)
	if err != nil {
		logger.Panic("Cannot create stream", zap.Error(err))
	}

	// create subscriber
	mFileStore := &mocks.FileStore{}
	mFileStore.On("UploadFile", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	s := subscriber.Subscriber{
		Natsjs:    natsjs,
		Logger:    logger,
		FileStore: mFileStore,
	}

	err = s.UploadSubscription()
	if err != nil {
		logger.Panic("Cannot create subscription", zap.Error(err))
	}

	p := &publisher.Publisher{
		Natsjs: natsjs,
		Logger: logger,
	}
	router := transport.Router{
		P:      p,
		Logger: logger,
	}

	r.Use(limits.RequestSizeLimiter(10000)) // limit 10KB
	r.POST("/user/batch", router.UserBatch())

	r.Run(":8080")

	// TODO: gracefulshutdown
}
