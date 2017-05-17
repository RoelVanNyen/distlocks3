package distlocks3

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	uuid "github.com/satori/go.uuid"
)

type TimeObjects []*s3.ObjectVersion

func (p TimeObjects) Len() int {
	return len(p)
}

// Define compare
func (p TimeObjects) Less(i, j int) bool {
	return (*p[i]).LastModified.Before(*p[j].LastModified)
}

// Define swap over an array
func (p TimeObjects) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func PutLockS3(bucket string, prefix string, region string) string {
	cfg := aws.NewConfig().WithRegion(region)
	svc := s3.New(session.Must(session.NewSession()), cfg)

	returnV, err := svc.PutObject(&s3.PutObjectInput{
		Body:   strings.NewReader(""),
		Bucket: &bucket,
		Key:    &prefix,
	})
	if err != nil {
		log.Printf(err.Error() + "\n")
		return ""
	} else {
		return *returnV.VersionId
	}
}

func DeleteLockS3(bucket string, prefix string, region string) {
	cfg := aws.NewConfig().WithRegion(region)
	svc := s3.New(session.Must(session.NewSession()), cfg)
	params := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(prefix),
	}
	_, err := svc.DeleteObject(params)
	if err != nil {
		fmt.Printf(err.Error())
	}
}

func GetOldestVersion(bucket string, prefix string, region string) s3.ObjectVersion {
	cfg := aws.NewConfig().WithRegion(region)
	svc := s3.New(session.Must(session.NewSession()), cfg)

	listParams := &s3.ListObjectVersionsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	}

	resp, _ := svc.ListObjectVersions(listParams)

	sort.Sort(TimeObjects(resp.Versions))
	for _, key := range resp.Versions {

		if *key.IsLatest == true {
			return *key
		}
	}
	fmt.Printf("Cannot find a last key !!!\n")
	return s3.ObjectVersion{}
}

func AquireLock(bucketName string, prefix string, region string) string {
	u1 := uuid.NewV4()
	fmt.Printf("LockID: %v\n", u1.String())

	rPrefix := fmt.Sprintf("%v/locks/%v.lock", prefix, u1.String())

	uploadedVersion := PutLockS3(bucketName, rPrefix, region)

	for true {
		versionPrefix := fmt.Sprintf("%v/locks", prefix)
		version := GetOldestVersion(bucketName, versionPrefix, region)
		if *version.VersionId == uploadedVersion {
			fmt.Printf("Aquired lock \n")
			return rPrefix
		} else {
			time.Sleep(time.Second * 10)
		}

	}
	return ""
}

func ReleaseLock(bucket string, prefix string, region string) {
	fmt.Printf("Released Lock \n")
	DeleteLockS3(bucket, prefix, region)
}
