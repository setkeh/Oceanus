package bucket

import (
	"context"
	"io"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Bucket struct {
	Ctx             context.Context
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	useSSL          bool
	Client          *minio.Client
	BucketName      string
	Location        string
}

func (b *Bucket) Init() {
	b.Ctx = context.Background()
	b.Endpoint = os.Getenv("BUCKET_ENDPOINT")
	b.AccessKeyID = os.Getenv("BUCKET_ACCESS_KEY")
	b.SecretAccessKey = os.Getenv("BUCKET_SECRET_KEY")
	b.useSSL = true
	b.BucketName = os.Getenv("BUCKET_NAME")
	b.Location = os.Getenv("BUCKET_LOCATION")

	// Initialize minio Client object.
	// TODO FIX THE ERROR HANDLING
	minioClient, _ := minio.New(b.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(b.AccessKeyID, b.SecretAccessKey, ""),
		Secure: b.useSSL,
	})

	b.Client = minioClient
}

func (b *Bucket) CreateBucket() error {
	err := b.Client.MakeBucket(b.Ctx, b.BucketName, minio.MakeBucketOptions{Region: b.Location})

	return err
}

func (b *Bucket) CheckBucketExists() bool {
	exists, errBucketExists := b.Client.BucketExists(b.Ctx, b.BucketName)
	if errBucketExists == nil && exists {
		return true
	}

	return false
}

func (b *Bucket) UploadImage(objectName string, file io.Reader, size int64) (int64, error) {
	userMetaData := map[string]string{"x-amz-acl": "public-read"}
	info, err := b.Client.PutObject(b.Ctx, b.BucketName, objectName, file, size, minio.PutObjectOptions{ContentType: "image/png", UserMetadata: userMetaData})
	return info.Size, err
}
