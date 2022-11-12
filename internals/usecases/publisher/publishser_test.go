package publisher

// TODO: write test.
// func TestUploadPublish(t *testing.T) {
// 	mNats := &mocks.NatsJetstream{}
// 	l := logs.WithPrefix("test")
// 	p := Publisher{
// 		Natsjs: mNats,
// 		Logger: l,
// 	}

// 	msgData := &entities.MsgData{
// 		ObjectName: "objectName",
// 		PathFile:   "pathFile",
// 	}

// 	mNats.On("PublishAsyncContext", mock.Anything, util.UploadSubject, mock.Anything).Return("12345", nil)

// 	err := p.UploadPublish(msgData)
// 	assert.NoError(t, err)
// }
