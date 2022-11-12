package subscriber

import (
	"context"
	"fmt"
	"github-com/edarha/uploadfile-test/internals/util"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"go.uber.org/zap"
)

type AWSClient struct {
	session *session.Session
	logger  *zap.Logger
	conf    *util.Config
}

func NewAWSClient(ctx context.Context, log *zap.Logger, conf *util.Config) (*AWSClient, error) {
	log.Info("Init AWS client")
	s3Sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(conf.AwsRegion),
		Credentials: credentials.NewStaticCredentials(conf.AwsAccessKey, conf.AwsSecretKey, conf.AwsToken),
	})
	if err != nil {
		return nil, fmt.Errorf("NewAWSClient: Init aws client failed: %s", err.Error())
	}

	awsClient := &AWSClient{
		session: s3Sess,
		logger:  log,
		conf:    conf,
	}

	return awsClient, nil
}

func (a *AWSClient) UploadFile(ctx context.Context, fileName string) error {
	_, cancel := context.WithTimeout(ctx, time.Second*a.conf.Timeout)
	defer cancel()

	pathName := fmt.Sprintf("/files/%s.json", fileName)
	f, err := os.Open(pathName)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = s3.New(a.session).PutObject(&s3.PutObjectInput{
		Bucket:             aws.String(a.conf.AwsBucket),
		Key:                aws.String(a.conf.AwsKey),
		ACL:                aws.String(a.conf.AwsACL),
		Body:               f,
		ContentDisposition: aws.String("attachment"),
	})

	return err
}

func (a *AWSClient) Close() error {
	return nil
}
