package main

import (
	"context"
	"fmt"
	"log"
	"matcha-comic-backend/models"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	godotenv.Load(".env")
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	accountID := os.Getenv("R2_ACCOUNT_ID")
	accessKey := os.Getenv("R2_ACCESS_KEY_ID")
	secretKey := os.Getenv("R2_SECRET_ACCESS_KEY")
	bucketName := os.Getenv("R2_BUCKET_NAME")
	publicURL := os.Getenv("R2_PUBLIC_URL")

	endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID)
	cfg := aws.Config{
		Credentials: credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		Region:      "auto",
		EndpointResolverWithOptions: aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: endpoint}, nil
			},
		),
	}
	r2Client := s3.NewFromConfig(cfg, func(o *s3.Options) { o.UsePathStyle = true })

	var media []models.Media
	db.Find(&media)

	fmt.Printf("Found %d media records to sync...\n", len(media))

	for _, m := range media {
		if m.Provider == "r2" {
			key := strings.TrimPrefix(m.URL, strings.TrimRight(publicURL, "/")+"/")
			head, err := r2Client.HeadObject(context.Background(), &s3.HeadObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(key),
			})
			if err == nil && head.ContentLength != nil {
				db.Model(&m).Update("size", *head.ContentLength)
				fmt.Printf("Updated %s: %d bytes\n", key, *head.ContentLength)
			} else {
				db.Model(&m).Update("size", 0)
				fmt.Printf("Failed to head %s: %v\n", key, err)
			}
		} else {
			// Local file
			localPath := strings.TrimPrefix(m.URL, "/uploads")
			info, err := os.Stat("uploads" + localPath)
			if err == nil {
				db.Model(&m).Update("size", info.Size())
				fmt.Printf("Updated local %s: %d bytes\n", localPath, info.Size())
			}
		}
	}
	fmt.Println("Sync completed!")
}
