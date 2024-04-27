package main

import (
	"context"
	"io"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	AccessKey = "xzzYekjqzUJ0dE8Ou16y"
	SecretKey = "sH8LbWc0tszIqNEYv9QDcC2BryTynwd8j1xfD0UK"
	Endpoint  = "192.168.56.101:9000"
)

func main() {
	client, err := minio.New(Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(AccessKey, SecretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}
	// bucketName := "testbucket"
	// location := "us-east-1"

	// ctx := context.Background()
	// err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{
	// 	Region: location,
	// })
	// if err != nil {
	// 	exists, errBucketExists := client.BucketExists(ctx, bucketName)
	// 	if errBucketExists == nil && exists {
	// 		log.Printf("We already own %s\n", bucketName)
	// 	} else {
	// 		log.Fatalln(err)
	// 	}
	// } else {
	// 	log.Printf("Successfully created %s\n", bucketName)
	// }

	// objectName := "testdata"
	// filePath := "./testdata"
	// contentType := "application/octet-stream"

	// info, err := client.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)

	bucketName := "testbucket"
	objectName := "testdata"
	ctx := context.Background()

	object, err := client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{
		ServerSideEncryption: nil,
		Checksum:             true,
	})
	if err != nil {
		log.Fatalln(err)
	}
	bs := make([]byte, 1024)
	n, err := object.Read(bs)
	if err != nil && err != io.EOF {
		log.Fatalln(err)
	}
	log.Println(string(bs[:n]))
}
