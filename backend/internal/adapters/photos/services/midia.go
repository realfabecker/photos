package services

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/realfabecker/photos/internal/core/ports"
	"time"
)

type S3MidiaSigner struct {
	presignCliente *s3.PresignClient
	bucketName     string
}

func NewS3MidiaSigner(bucketName string, client *s3.Client) ports.MidiaSigner {
	return &S3MidiaSigner{
		bucketName:     bucketName,
		presignCliente: s3.NewPresignClient(client),
	}
}

func (s S3MidiaSigner) GetObjectUrl(name string, lifetime int64) (string, error) {
	request, err := s.presignCliente.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String("midia/" + name),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(lifetime * int64(time.Second))
	})
	if err != nil {
		return "", err
	}
	return request.URL, nil
}

func (s S3MidiaSigner) PutObjectUrl(name string, lifetime int64) (string, error) {
	request, err := s.presignCliente.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String("midia/" + name),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(lifetime * int64(time.Second))
	})
	if err != nil {
		return "", err
	}
	return request.URL, nil
}
