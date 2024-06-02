package storage

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type AwsS3Config struct {
	Region string
	Id     string
	Secret string
	Token  string
}

func NewAwsConnection(config *AwsS3Config) (*s3.S3, error) {

	s3Config := &aws.Config{
		Region:      aws.String(config.Region),
		Credentials: credentials.NewStaticCredentials(config.Id, config.Secret, config.Token),
	}

	s3Session, err := session.NewSession(s3Config)
	if err != nil {
		fmt.Println("Error creating AWS session:", err)
		return nil, err
	}
	s3Client := s3.New(s3Session)

	return s3Client, nil
	
}
