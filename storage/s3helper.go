package storage

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/KevinFagan/discord-memories/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

// Sync will ensure to create all the necessary "folders" defined within
// the memories config file. If the "folders" already exist, the creation process will be skipped
func Sync(service *s3.S3, config config.Config, bucket string) error {
	for _, v := range config.Options {
		if exists, err := ObjectExists(service, bucket, v.Path); err != nil && !exists {
			logrus.Infof("initializing %s", v.Path)
			_, err := service.PutObject(&s3.PutObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(v.Path),
				Body:   nil,
				ACL:    aws.String("private"),
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// ObjectExsits checks if an object exists within a specified bucket
// If it does not exist, an error will be returned
func ObjectExists(service *s3.S3, bucket, key string) (bool, error) {
	_, err := service.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		if s3err, ok := err.(awserr.Error); ok && s3err.Code() == s3.ErrCodeNoSuchKey {
			return false, nil // Object does not exist
		}
		return false, err // Other error
	}
	return true, nil // Object exists
}

// GetRandomObjectUnderPrefix will retrieve a random object under a known prefix. If no objects
// are found under the prefix, an error will be returned. Otherwise, the S3 object and its name will be returned
func GetRandomObjectUnderPrefix(service *s3.S3, bucket, prefix string) (*s3.GetObjectOutput, string, error) {
	objectList, err := service.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		return nil, "", err
	}

	// Removing the parent (prefix) from list of contents
	for i := range objectList.Contents {
		if *objectList.Contents[i].Key == prefix {
			objectList.Contents = append(objectList.Contents[:i], objectList.Contents[i+1:]...)
			break
		}
	}
	if len(objectList.Contents) == 0 {
		return nil, "", fmt.Errorf("no objects found under prefix %s", prefix)
	}

	// Selecting a random object
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := r.Intn(len(objectList.Contents))
	o, err := service.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    objectList.Contents[n].Key,
	})
	if err != nil {
		return nil, "", err
	}

	return o, *objectList.Contents[n].Key, err
}

// UploadObject will upload an object to the specified bucket
func UploadObject(service *s3.S3, bucket, prefix string, attachment discordgo.MessageAttachment) error {
	body, err := getFileBytes(attachment.URL)
	if err != nil {
		return err
	}

	_, err = service.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(fmt.Sprintf("%s/%s", prefix, attachment.Filename)),
		Body:          bytes.NewReader(body),
		ContentLength: aws.Int64(int64(len(body))),
		ContentType:   aws.String(attachment.ContentType),
	})

	return err
}

func getFileBytes(url string) ([]byte, error) {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the bytes from the response body
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
