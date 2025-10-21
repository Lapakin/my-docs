package minio

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/lapotkin/file-storage/pkg/file-storage/repository"
)

type RepoManager struct {
	bucketName string
}

func NewRepoManager() *RepoManager {
	return &RepoManager{}
}

func (rm *RepoManager) NewObjectStorageRepo(s3Client *s3.Client) repository.ObjectStorageRepository {
	return NewObjectStorageRepository(s3Client)
}
