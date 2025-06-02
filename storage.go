package object_storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/logging"
	"gopkg.in/yaml.v3"
)

type objectStorage struct {
	s3_client    *s3.Client
	baseEndpoint string
	bucket       string
	metaFilename string
}

func New(ctx context.Context, cfg *Config) (ObjectStorage, error) {
	var credentialsProviderFunc aws.CredentialsProviderFunc = func(context.Context) (aws.Credentials, error) {
		return aws.Credentials{
			AccessKeyID:     cfg.AccessKey,
			SecretAccessKey: cfg.SecretKey,
		}, nil
	}

	s3_config, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(cfg.Region),
		config.WithBaseEndpoint(cfg.BaseEndpoint),
		config.WithCredentialsProvider(credentialsProviderFunc),
		config.WithLogger(logging.Nop{}),
	)
	if err != nil {
		return nil, err
	}

	s3_client := s3.NewFromConfig(s3_config)

	return &objectStorage{
		s3_client:    s3_client,
		baseEndpoint: cfg.BaseEndpoint,
		bucket:       cfg.Bucket,
		metaFilename: cfg.MetaFilename,
	}, nil
}

func (s *objectStorage) GetObject(ctx context.Context, path string) (*Object, error) {
	url, err := url.JoinPath(s.baseEndpoint, s.bucket, path)
	if err != nil {
		return nil, err
	}

	object := Object{
		Path:     path,
		Meta:     make(map[string]any),
		Contains: make([]string, 0),
		URL:      url,
	}

	// single file
	if len(path) != 0 && path[len(path)-1] != '/' {
		_, err := s.s3_client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    &path,
		})
		if err != nil {
			return nil, err
		}

		return &object, nil
	}

	// directory
	list_dir_files_resp, err := s.s3_client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: &s.bucket,
		Prefix: &path,
	})

	if err != nil {
		return nil, err
	}

	for _, file := range list_dir_files_resp.Contents {
		if file.Key == nil {
			continue
		}

		key := *file.Key
		if len(key) == 0 {
			continue
		}

		filename := key[len(path):]
		if len(filename) == 0 {
			continue
		}

		parts := strings.Split(filename, "/")

		if len(parts) == 2 && len(parts[1]) == 0 {
			object.Contains = append(object.Contains, filename)
			continue
		}

		if len(parts) != 1 {
			continue
		}

		if filename == s.metaFilename {
			resp, err := s.s3_client.GetObject(ctx, &s3.GetObjectInput{
				Bucket: aws.String(s.bucket),
				Key:    &key,
			})
			if err != nil {
				return nil, err
			}

			d := yaml.NewDecoder(resp.Body)
			err = d.Decode(&object.Meta)
			if err != nil {
				return nil, err
			}
			continue
		}

		object.Contains = append(object.Contains, filename)
	}

	return &object, nil
}

func (s *objectStorage) GetObjectBody(ctx context.Context, obj *Object) (io.ReadCloser, error) {
	if obj == nil {
		return nil, fmt.Errorf("Object is nil")
	}

	if len(obj.Path) == 0 || obj.Path[len(obj.Path)-1] == '/' {
		return nil, fmt.Errorf("Object is directory.")
	}
	resp, err := s.s3_client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    &obj.Path,
	})
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
