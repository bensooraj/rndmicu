package s3engine

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/joho/godotenv"
)

type TASK_TYPE int

const (
	UPLOAD_TASK TASK_TYPE = iota
	DELETE_TASK
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

	S3                     *s3.S3
	Uploader               *s3manager.Uploader
	UploadChunkSizeInBytes int64
	Concurrency            int

	JobChannels []chan *AudioFileJob
}

func NewS3Engine(opts NewS3EngineOptions) *S3Engine {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	s3e := S3Engine{
		AccessKeyID:            os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretAccessKey:        os.Getenv("AWS_SECRET_ACCESS_KEY"),
		DefaultRegion:          os.Getenv("AWS_DEFAULT_REGION"),
		ErrorChannel:           opts.ErrorChannel,
		UploadChunkSizeInBytes: opts.UploadChunkSizeInBytes,
		Concurrency:            opts.Concurrency,
		Ctx:                    opts.Ctx,
		BucketName:             opts.BucketName,
	}

	return &s3e
}

func (s3e *S3Engine) Init() error {

	awsConfig := aws.
		NewConfig().
		WithCredentials(credentials.NewStaticCredentials(
			s3e.AccessKeyID,
			s3e.SecretAccessKey,
			"",
		)).
		WithRegion(s3e.DefaultRegion)

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return err
	}

	s3e.S3 = s3.New(sess)

	uploader := s3manager.NewUploader(sess, func(u *s3manager.Uploader) {
		u.Concurrency = s3e.Concurrency
		u.PartSize = s3e.UploadChunkSizeInBytes // 10 * 1024 * 1024 // The minimum/default allowed part size is 5MB
	})

	s3e.Uploader = uploader

	return nil
}

func (s3e *S3Engine) StartWorkers() {
	s3e.JobChannels = make([]chan *AudioFileJob, s3e.Concurrency)
	for i := 0; i < s3e.Concurrency; i++ {
		s3e.JobChannels[i] = make(chan *AudioFileJob)
		go s3e.Worker(s3e.JobChannels[i])
	}
}

func (s3e *S3Engine) QueueJob(job *AudioFileJob) {
	s3e.JobChannels[taskCounter%len(s3e.JobChannels)] <- job
}

func (s3e *S3Engine) StopWorkers() {
	for _, ch := range s3e.JobChannels {
		close(ch)
	}
}

func (s3e *S3Engine) Worker(jobChannels chan *AudioFileJob) {
	n, _ := rand.Int(rand.Reader, big.NewInt(3))
	flushTimeout := time.After(time.Duration(1+n.Int64()) * time.Minute)
	for {
		select {
		case <-s3e.Ctx.Done():
			log.Println("")
			return

		case job, ok := <-jobChannels:
			if !ok {
				log.Println("jobChannels closed")
				return
			}

			var err error
			switch job.TaskType {
			case UPLOAD_TASK:
				log.Println("PutObjectInput.Bucket: s3e.BucketName", s3e.BucketName)
				log.Println("PutObjectInput.Bucket: Key", fmt.Sprintf("%s/%s", job.KeyName, job.Filename))
				_, err = s3e.Uploader.Upload(&s3manager.UploadInput{
					Bucket:      aws.String(s3e.BucketName),
					Key:         aws.String(fmt.Sprintf("%s/%s", job.KeyName, job.Filename)),
					Body:        job.File,
					ContentType: aws.String(job.ContentType),
				})
				if err != nil {
					s3e.ErrorChannel <- err
					log.Println("Error uploading", err)
					continue
				}
			case DELETE_TASK:
				_, err = s3e.S3.DeleteObject(&s3.DeleteObjectInput{
					Bucket: aws.String(s3e.BucketName),
					Key:    aws.String(fmt.Sprintf("%s/%s", job.KeyName, job.Filename)),
				})
				if err != nil {
					s3e.ErrorChannel <- err
					continue
				}
			}
		case f := <-flushTimeout:
			go s3e.Worker(jobChannels)
			log.Println("Flushing and re-spawn at", f.String())
			return
		}
	}
}
