package transport

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github-com/edarha/uploadfile-test/internals/usecases/entities"
	"github-com/edarha/uploadfile-test/internals/usecases/publisher"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Router struct {
	P      *publisher.Publisher
	Logger *zap.Logger
}

// Serve api /user/batch
// Create a file with uuid.
// Publish msg to nats with `upload.send` subject.
func (r *Router) UserBatch() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			r.Logger.Error("Cannot read the payload", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Cannot read the payload",
			})
			return
		}

		fileName := fmt.Sprintf("%s.json", uuid.New())
		tmpDir := os.TempDir()
		path := filepath.Join(tmpDir, fileName)

		err = os.WriteFile(path, data, 0666)
		if err != nil {
			r.Logger.Error("Cannot write a file", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Cannot write a file",
			})
			return
		}

		// create msg data
		msgData := &entities.MsgData{
			FileName: fileName,
			Path:     path,
		}

		// publish msg to nats
		err = r.P.UploadPublish(msgData)
		if err != nil {
			r.Logger.Error("Cannot publish a message", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Cannot publish a message",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Upload data success",
		})
	}
}
