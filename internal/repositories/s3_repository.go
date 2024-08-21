package repositories

import (
	"Backend_golang_project/infrastructure/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/sirupsen/logrus"
	"io"
)

type S3RepositoryInterface interface {
	UploadFile(bucket, key string, body io.ReadSeeker) error
}

type S3Repository struct {
	service    *s3.S3
	s3Uploader *s3manager.Uploader
}

func NewS3Repository(logger *logrus.Logger, cf *config.Config) (S3RepositoryInterface, error) {
	awsConfig := &aws.Config{
		Region: aws.String(cf.S3Config.Region),
		Credentials: credentials.NewStaticCredentials(
			cf.S3Config.AwsId,
			cf.S3Config.AwsKey,
			"",
		),
		Endpoint:         aws.String(cf.S3Config.Endpoint),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}
	session, err := session.NewSession(awsConfig)
	if err != nil {
		logger.Error("cannot create s3 client")
	}
	return &S3Repository{
		service:    s3.New(session),
		s3Uploader: s3manager.NewUploader(session),
	}, err
}

func (r *S3Repository) UploadFile(bucket, key string, body io.ReadSeeker) error {
	_, err := r.service.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   body,
	})
	return err
}
