package s3engine

import (
	"context"
	"io"
)

var (
	taskCounter int = 0
)

type AudioFileJob struct {
	TaskType    TASK_TYPE
	File        io.Reader
	Filename    string
	Size        int64
	ContentType string
	KeyName     string
}

type NewS3EngineOptions struct {
	ErrorChannel           chan error
	Ctx                    context.Context
	UploadChunkSizeInBytes int64
	Concurrency            int

	BucketName string
}
