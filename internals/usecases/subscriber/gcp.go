package subscriber

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github-com/edarha/uploadfile-test/internals/util"

	"cloud.google.com/go/storage"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

type GCPClient struct {
	client *storage.Client
	logger *zap.Logger
	conf   *util.Config
}

func NewGCPClient(ctx context.Context, log *zap.Logger, conf *util.Config) (*GCPClient, error) {
	log.Info("Init GCP client")
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(conf.GoogleCredentialPath))
	if err != nil {
		return nil, fmt.Errorf("NewGCPClient: Init gcp client failed: %s", err.Error())
	}

	gcpClient := &GCPClient{
		client: client,
		logger: log,
		conf:   conf,
	}

	return gcpClient, nil
}

func (g *GCPClient) UploadFile(ctx context.Context, fileName string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*g.conf.Timeout)
	defer cancel()

	pathFile := fmt.Sprintf("/files/%s", fileName)
	f, err := os.Open(pathFile)
	if err != nil {
		return err
	}
	defer f.Close()

	wc := g.client.Bucket(g.conf.GoogleBucket).Object(g.conf.ObjectPath + fileName).NewWriter(ctx)
	if _, err := io.Copy(wc, f); err != nil {
		return err
	}

	if err := wc.Close(); err != nil {
		return err
	}

	g.logger.Info("UploadFile: upload file success", zap.String("object_path", g.conf.ObjectPath), zap.String("object", fileName))
	return nil
}

func (g *GCPClient) Close() error {
	return g.client.Close()
}
