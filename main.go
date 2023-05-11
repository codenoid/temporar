package main

import (
	"context"
	"embed"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	gonanoid "github.com/matoous/go-nanoid/v2"

	_ "github.com/joho/godotenv/autoload"
)

//go:embed all:frontend/build
var f embed.FS

var s3client *s3.Client

func init() {
	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: os.Getenv("S3_ENDPOINT_URL"),
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(os.Getenv("S3_ACCESS_KEY"), os.Getenv("S3_SECRET_KEY"), "")),
	)
	if err != nil {
		log.Fatal(err)
	}

	s3client = s3.NewFromConfig(cfg)
}

func main() {

	r := gin.Default()

	r.POST("/", func(c *gin.Context) {

		fileID, err := gonanoid.Generate("abcdefghijklmnopqrstuvwxyz1234567890", 6)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		defer file.Close()

		expAt := time.Now().Add(time.Hour * 1)
		_, err = s3client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
			Key:    aws.String("temporar/" + fileID),
			Body:   file,
			Metadata: map[string]string{
				"file_name": header.Filename,
			},
			Expires: aws.Time(expAt),
		})
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
			return
		}

		scheme := "https"
		if c.Request.TLS == nil {
			scheme = "http"
		}
		c.JSON(http.StatusOK, gin.H{
			"success":    true,
			"id":         fileID,
			"expired_at": expAt.Unix(),
			"name":       header.Filename,
			"link":       scheme + "://" + c.Request.Host + "/" + fileID,
		})
	})

	r.GET("/:file_id", func(c *gin.Context) {
		fileID := c.Param("file_id")

		obj, err := s3client.GetObject(context.TODO(), &s3.GetObjectInput{
			Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
			Key:    aws.String("temporar/" + fileID),
		})

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.DataFromReader(http.StatusOK, obj.ContentLength, "application/octet-stream", obj.Body, map[string]string{
			"Content-Disposition": "attachment; filename=" + obj.Metadata["file_name"],
		})

		s3client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
			Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
			Key:    aws.String("temporar/" + fileID),
		})
	})

	r.Use(svelteStaticLoader)

	bindAddr := ":8080"
	if os.Getenv("BIND_ADDR") != "" {
		bindAddr = os.Getenv("BIND_ADDR")
	}
	r.Run(bindAddr)
}
