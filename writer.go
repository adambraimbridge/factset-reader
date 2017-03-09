package main

import (
	"path"
	"time"

	log "github.com/Sirupsen/logrus"
	"strings"
)

type Writer interface {
	Write(src string, localFileName string, s3FileName string, archive string) error
}

type S3Writer struct {
	s3Client S3Client
}

func NewWriter(config s3Config) (Writer, error) {
	s3, err := NewS3Client(config)
	return &S3Writer{s3Client: s3}, err
}

func (s3w *S3Writer) Write(src string, localFileName string, s3FileName string, archive string) error {
	log.Infof("Writing file [%s]\n", s3FileName)
	s3ResFilePath := s3w.getS3ResFilePath(s3FileName, archive)
	p := path.Join(src, localFileName)
	n, err := s3w.s3Client.PutObject(s3ResFilePath, p)
	if err != nil {
		return err
	}
	log.Infof("Uploaded file [%s] of size [%d] successfully", s3ResFilePath, n)
	return nil
}

func (s3w *S3Writer) getS3ResFilePath(s3FileName string, archive string) string {
	var resFilePath string
	if archive == "" {
		return archive
	}

	if strings.Contains(archive, "full") {
		resFilePath = "Weekly/" + time.Now().Format("2006-01-02") + "/" + s3FileName
	} else {
		resFilePath = "Daily/" + time.Now().Format("2006-01-02") + "/" + s3FileName
	}

	return resFilePath
}
