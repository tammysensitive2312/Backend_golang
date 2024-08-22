package repositories

import (
	"Backend_golang_project/infrastructure/config"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/sirupsen/logrus"
	"io"
)

type S3RepositoryInterface interface {
	UploadFile(ctx context.Context, bucket, key string, body *io.PipeReader) error
}

type S3Repository struct {
	s3Uploader *manager.Uploader
	client     *s3.Client
}

func NewS3Repository(ctx context.Context, logger *logrus.Logger, cf *config.Config) (S3RepositoryInterface, error) {
	// Tạo AWS config
	awsCfg, err := awsConfig.LoadDefaultConfig(ctx,
		awsConfig.WithRegion(cf.S3Config.Region),
		awsConfig.WithCredentialsProvider(credentials.NewStaticCredentials(
			cf.S3Config.AwsId,
			cf.S3Config.AwsKey,
			""),
		),
	)
	if err != nil {
		logger.WithError(err).Error("Failed to load AWS config")
		return nil, err
	}

	// Tạo S3 client
	// ở đây em sử dụng baseEndpoint là một phiên bản cải tiến hiện đại
	// thay thế cho EndpointResolver
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.BaseEndpoint = aws.String(cf.S3Config.Endpoint)
	})

	// Tạo S3 uploader
	uploader := manager.NewUploader(client)

	return &S3Repository{
		client:     client,
		s3Uploader: uploader,
	}, nil
}

// UploadFile ở hàm này ban đầu em xử lý theo kiểu lưu trữ trước vào bộ nhớ rồi mới upload data lên s3
// nhưng sau khi nhận thấy vấn đề với file lớn thì em tìm hiểu cách khác để xử lý và sử dụng s3 manager và nâng cấp aws sdk go v2
func (r *S3Repository) UploadFile(ctx context.Context, bucket, key string, pipeReader *io.PipeReader) error {
	// Cấu hình uploader với các tùy chọn
	uploader := manager.NewUploader(r.client, func(u *manager.Uploader) {
		u.PartSize = 5 * 1024 * 1024
		u.Concurrency = 5
	})

	// Thực hiện upload
	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   pipeReader,
	})

	// Đóng pipeReader sau khi hoàn thành
	if pipeReader != nil {
		pipeReader.Close()
	}

	return err
}
