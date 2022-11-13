package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github-com/edarha/uploadfile-test/internals/logs"
	"github-com/edarha/uploadfile-test/internals/usecases/mocks"
	"github-com/edarha/uploadfile-test/internals/usecases/publisher"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type responseBody struct {
	Message string
}

func TestUserBatch(t *testing.T) {
	data := "body-test"
	reqBody, err := json.Marshal(data)
	assert.NoError(t, err)

	t.Run("Success: the status is 200", func(t *testing.T) {
		_, w, err := userBatchPublishSuccess(bytes.NewBuffer(reqBody))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		body, err := ioutil.ReadAll(w.Body)
		assert.NoError(t, err)

		res := &responseBody{}
		err = json.Unmarshal(body, res)
		assert.NoError(t, err)
		assert.Equal(t, "Upload data success", res.Message)
	})

	t.Run("Fail: the status is 400, Cannot publish a message", func(t *testing.T) {
		_, w, err := userBatchPublishFail(bytes.NewBuffer(reqBody))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		body, err := ioutil.ReadAll(w.Body)
		assert.NoError(t, err)

		res := &responseBody{}
		err = json.Unmarshal(body, res)
		assert.NoError(t, err)
		assert.Equal(t, "Cannot publish a message", res.Message)
	})

	t.Run("Fail: the status is 413, Payload is over 10KB", func(t *testing.T) {
		data, err := os.ReadFile("../../files/payload-heavy.json")
		assert.NoError(t, err)

		_, w, err := userBatchPublishFail(bytes.NewBuffer(data))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusRequestEntityTooLarge, w.Code)

		body, err := ioutil.ReadAll(w.Body)
		assert.NoError(t, err)

		res := &responseBody{}
		err = json.Unmarshal(body, res)
		assert.NoError(t, err)
		assert.Equal(t, "Payload is over 10KB", res.Message)
	})
}

func userBatchPublishSuccess(body *bytes.Buffer) (*http.Request, *httptest.ResponseRecorder, error) {
	logger := logs.NewZapLogger("test")
	mNats := new(mocks.NatsJetstream)
	mNats.On("PublishAsyncContext", mock.Anything, mock.Anything, mock.Anything).Return("", nil)
	p := &publisher.Publisher{
		Natsjs: mNats,
		Logger: logger,
	}
	r := gin.New()
	router := Router{
		P:      p,
		Logger: logger,
	}

	r.Use(router.CheckLimitPayload(10000))
	r.POST("/", router.UserBatch())
	req, err := http.NewRequest(http.MethodPost, "/", body)
	if err != nil {
		return req, httptest.NewRecorder(), err
	}

	req.Header.Set("Content-type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return req, w, nil
}

func userBatchPublishFail(body *bytes.Buffer) (*http.Request, *httptest.ResponseRecorder, error) {
	logger := logs.NewZapLogger("test")
	mNats := new(mocks.NatsJetstream)
	mNats.On("PublishAsyncContext", mock.Anything, mock.Anything, mock.Anything).Return("", fmt.Errorf("Cannot publish msg"))

	p := &publisher.Publisher{
		Natsjs: mNats,
		Logger: logger,
	}
	r := gin.New()
	router := Router{
		P:      p,
		Logger: logger,
	}

	r.Use(router.CheckLimitPayload(10000))
	r.POST("/", router.UserBatch())
	req, err := http.NewRequest(http.MethodPost, "/", body)
	if err != nil {
		return req, httptest.NewRecorder(), err
	}

	req.Header.Set("Content-type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return req, w, nil
}
