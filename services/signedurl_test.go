package services

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3PresignGetObjectAPIImp struct{}

func (dt S3PresignGetObjectAPIImp) PresignGetObject(ctx context.Context,
	params *s3.GetObjectInput,
	optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error) {

	output := &v4.PresignedHTTPRequest{
		URL: fmt.Sprintf("https://%s.s3.ca-central-1.amazonaws.com/%s", *params.Bucket, *params.Key),
	}

	return output, nil
}

func TestPresignGetObject(t *testing.T) {
	api := &S3PresignGetObjectAPIImp{}

	getConfig(t)

	keyObj := "estimate/est-99.pdf"

	input := &s3.GetObjectInput{
		Bucket: aws.String(cfg.S3Bucket),
		Key:    aws.String(keyObj),
	}

	resp, err := GetPresignedURL(context.TODO(), *api, input)
	if err != nil {
		t.Fatalf("Expected null error received: %s", err)
	}

	expectedURL := fmt.Sprintf("https://%s.s3.ca-central-1.amazonaws.com/%s", cfg.S3Bucket, keyObj)
	if expectedURL != resp.URL {
		t.Fatalf("URL should be: %s, have: %s", expectedURL, resp.URL)
	}

	t.Logf("URL: %s", expectedURL)
}

func TestIntegCreateSignedURL(t *testing.T) {

	getConfig(t)

	fileObject := "estimate/est-1005.pdf"
	url, err := CreateSignedURL(cfg, fileObject)
	if err != nil {
		t.Fatalf("Expected null error received: %s", err)
	}

	t.Logf("url: %+v\n", url)
}
