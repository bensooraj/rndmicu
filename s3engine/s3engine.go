package s3engine

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/joho/godotenv"
)

type S3Engine struct {
	// Creds
	AccessKeyID     string
	SecretAccessKey string
	DefaultRegion   string
	BucketName      string

	// Communication
	ErrorChannel chan error
	Ctx          context.Context

	Uploader               *s3manager.Uploader
	UploadChunkSizeInBytes int64
	Concurrency            int

	UploadJobChannels []chan *AudioFileUploadJob
}

func NewS3Engine(opts NewS3EngineOptions) *S3Engine {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}
	var s3 S3Engine
	s3 = S3Engine{
		AccessKeyID:            os.Getenv("DB_USERNAME"),
		SecretAccessKey:        os.Getenv("DB_USERNAME"),
		DefaultRegion:          os.Getenv("DB_USERNAME"),
		ErrorChannel:           opts.ErrorChannel,
		UploadChunkSizeInBytes: opts.UploadChunkSizeInBytes,
		Concurrency:            opts.Concurrency,
	}

	return &s3
}

func (s3 *S3Engine) Init() error {

	awsConfig := aws.
		NewConfig().
		WithCredentials(credentials.NewStaticCredentials(
			s3.AccessKeyID,
			s3.SecretAccessKey,
			"",
		)).
		WithRegion(s3.DefaultRegion)

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return err
	}

	uploader := s3manager.NewUploader(sess, func(u *s3manager.Uploader) {
		u.Concurrency = s3.Concurrency
		u.PartSize = s3.UploadChunkSizeInBytes //10 * 1024 * 1024 // The minimum/default allowed part size is 5MB
	})

	s3.Uploader = uploader

	return nil
}

func (s3 *S3Engine) StartUploadWorkers() {
	for i := 0; i < s3.Concurrency; i++ {
		s3.UploadJobChannels[i] = make(chan *AudioFileUploadJob)
		go s3.UploadWorker(s3.UploadJobChannels[i])
	}
}

func (s3 *S3Engine) QueueUploadJob(job *AudioFileUploadJob) {
	s3.UploadJobChannels[taskCounter%len(s3.UploadJobChannels)] <- job
}

func (s3 *S3Engine) StopUploadWorkers() {
	for _, ch := range s3.UploadJobChannels {
		close(ch)
	}
}

func (s3 *S3Engine) UploadWorker(uploadJobChannels chan *AudioFileUploadJob) {
	flushTimeout := time.After(time.Duration(1+rand.Intn(3)) * time.Minute)
	for {
		select {
		case <-s3.Ctx.Done():
			log.Println("")
			return

		case job, ok := <-uploadJobChannels:
			if !ok {
				log.Println("uploadJobChannels closed")
				return
			}
			_, err := s3.Uploader.Upload(&s3manager.UploadInput{
				Bucket:      aws.String(s3.BucketName),
				Key:         aws.String(fmt.Sprintf("%s/%s", job.KeyName, job.Filename)),
				Body:        job.File,
				ContentType: aws.String(job.ContentType),
			})
			if err != nil {
				s3.ErrorChannel <- err
				log.Println("Error uploading")
				continue
			}
		case f := <-flushTimeout:
			go s3.UploadWorker(uploadJobChannels)
			log.Println("Flushing and re-spawn at", f.String())
			return
		}
	}
}
