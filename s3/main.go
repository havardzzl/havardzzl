package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type MyS3Client struct {
	endpoint  string
	accessKey string
	secretKey string
	s3sess    *session.Session
	client    *s3.S3
}

// NewMyS3Client
// endpoint 是 ip:port 或 domain:port 的形式，不加 http 或 https
func NewMyS3Client(endpoint string, accessKey string, secretKey string, signingRegion string) *MyS3Client {
	useSSL := strings.HasSuffix(endpoint, ":443")
	theSess := createS3Session(endpoint, accessKey, secretKey, signingRegion, useSSL)
	return &MyS3Client{
		endpoint:  endpoint,
		accessKey: accessKey,
		secretKey: secretKey,
		s3sess:    theSess,
		client:    s3.New(theSess),
	}
}

type credProvider struct {
	AccessKey string
	SecretKey string
}

func (m credProvider) Retrieve() (credentials.Value, error) {
	return credentials.Value{
		AccessKeyID:     m.AccessKey,
		SecretAccessKey: m.SecretKey,
	}, nil
}

func (m credProvider) IsExpired() bool {
	return false
}

func createS3Session(endpoint string, accessKey string, secretKey string, signingRegion string, useSSL bool) *session.Session {
	disableSSL := !useSSL
	forcePathStyle := true
	cp := credProvider{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}

	myCustomResolver := func(service, region string, optFns ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error) {
		if service == endpoints.S3ServiceID {
			return endpoints.ResolvedEndpoint{
				URL:           endpoint,
				SigningRegion: signingRegion,
			}, nil
		}

		return endpoints.DefaultResolver().EndpointFor(service, region, optFns...)
	}

	s3Session := session.Must(session.NewSession(&aws.Config{
		Region:           aws.String("us-east-1"),
		Endpoint:         aws.String(endpoint),
		EndpointResolver: endpoints.ResolverFunc(myCustomResolver),
		Credentials:      credentials.NewCredentials(cp),
		DisableSSL:       &disableSSL,
		S3ForcePathStyle: &forcePathStyle,
	}))

	return s3Session
}

// PutObject
// 上传文件
//
// bucketName 要上传到哪个 bucket
// key 上传之后对象的 key
// localFileName 本地文件路径，绝对相对都行
func (ms3 *MyS3Client) PutObject(bucketName string, key string, localFileName string) error {
	file, err := os.Open(localFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	tempBkName := bucketName
	tempKey := key
	putObjectInput := s3.PutObjectInput{
		Bucket: &tempBkName,
		Key:    &tempKey,
		Body:   file,
	}
	_, err = ms3.client.PutObject(&putObjectInput)
	if err != nil {
		return err
	}

	return nil
}

// PutObjectMultipart
// 分段上传文件
//
// bucketName 要上传到哪个 bucket
// key 上传之后对象的 key
// localFileName 本地文件路径，绝对相对都行
// partSize 分段大小，单位字节
func (ms3 *MyS3Client) PutObjectMultipart(bucketName string, key string, localFileName string, partSize int64) error {
	file, err := os.Open(localFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	tempBkName := bucketName
	tempKey := key

	uploader := s3manager.NewUploader(ms3.s3sess, func(u *s3manager.Uploader) {
		u.PartSize = partSize
	})

	var ui s3manager.UploadInput
	ui.Bucket = &tempBkName
	ui.Key = &tempKey
	ui.Body = file
	_, err = uploader.Upload(&ui)
	if err != nil {
		return err
	}

	return nil
}

func (ms3 *MyS3Client) GetObject(bucket string, key string) {
	var oi s3.GetObjectInput
	oi.Bucket = &bucket
	oi.Key = &key

	oo, err := ms3.client.GetObject(&oi)
	if err != nil {
		fmt.Println(err.Error())
	}

	data, err := io.ReadAll(oo.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(string(data))
}

func (ms3 *MyS3Client) DeleteObject(bucket string, key string) error {
	do := s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}

	_, err := ms3.client.DeleteObject(&do)

	return err
}

func (ms3 *MyS3Client) ListObject(bucket string) {
	listObjectInput := s3.ListObjectsInput{Bucket: &bucket}
	listObjectsOutput, err := ms3.client.ListObjects(&listObjectInput)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		for _, obj := range listObjectsOutput.Contents {
			fmt.Println(*obj.Key)
		}
	}
}

func main() {
	ms3 := NewMyS3Client("s3-tcloud-online.kanzhun.tech", "SPJH33VX40G4ZLVKG22K", "Cjd885dYPeMJMXtS4hrn1QISOKyG3pfl0bvQyjk2", "default")

	args := os.Args[1:]
	if len(args) < 2 {
		fmt.Println("usage: s3-client <bucket> <key> <localFileName>")
		return
	}
	var err error
	switch args[0] {
	case "put":
		err = ms3.PutObject(args[1], args[2], args[3])
	case "get":
		ms3.GetObject(args[1], args[2])
	case "delete":
		err = ms3.DeleteObject(args[1], args[2])
	case "list":
		ms3.ListObject(args[1])
	case "multipart":
		err = ms3.PutObjectMultipart(args[1], args[2], args[3], 100*1024*1024)
	default:
		fmt.Println("usage: s3-client <bucket> <key> <localFileName>")
	}
	fmt.Println("err: ", err)
}
