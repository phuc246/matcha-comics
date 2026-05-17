package handlers

import (
	"bytes"
	"context"
	"fmt"
	"matcha-comic-backend/models"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

// newR2Client creates an S3-compatible client for Cloudflare R2
func newR2Client() (*s3.Client, error) {
	accountID := os.Getenv("R2_ACCOUNT_ID")
	accessKey := os.Getenv("R2_ACCESS_KEY_ID")
	secretKey := os.Getenv("R2_SECRET_ACCESS_KEY")

	if accountID == "" || accessKey == "" || secretKey == "" {
		return nil, fmt.Errorf("R2 credentials not configured")
	}

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

	return s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	}), nil
}

// UploadImages uploads WebP images to Cloudflare R2 and returns public URLs
// Falls back to local storage if R2 is not configured
func (h *Handler) UploadImages(c *gin.Context) {
	if err := c.Request.ParseMultipartForm(100 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request too large"})
		return
	}

	storyID := c.PostForm("story_id")
	chapterNum := c.PostForm("chapter_number")
	files := c.Request.MultipartForm.File["images"]

	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No images provided"})
		return
	}

	// Try R2 first, fallback to local
	r2Client, err := newR2Client()
	useR2 := err == nil

	bucketName := os.Getenv("R2_BUCKET_NAME")
	publicURL := os.Getenv("R2_PUBLIC_URL")

	type fileInfo struct {
		url  string
		size int64
	}
	uploadedFiles := make([]fileInfo, 0, len(files))
	urls := make([]string, 0, len(files))

	for i, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			continue
		}

		// Read file bytes
		buf := make([]byte, fileHeader.Size)
		file.Read(buf)
		file.Close()

		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		if ext == "" {
			ext = ".webp"
		}

		timestamp := time.Now().UnixNano()
		filename := fmt.Sprintf("%04d_%d%s", i+1, timestamp, ext)
		// R2 key path: stories/{story_id}/chapters/{chapter_num}/001_xxx.webp
		key := fmt.Sprintf("stories/%s/chapters/%s/%s", storyID, chapterNum, filename)

		currentURL := ""
		if useR2 {
			// Upload to Cloudflare R2
			contentType := "image/webp"
			switch ext {
			case ".png":
				contentType = "image/png"
			case ".jpg", ".jpeg":
				contentType = "image/jpeg"
			}

			_, uploadErr := r2Client.PutObject(context.Background(), &s3.PutObjectInput{
				Bucket:      aws.String(bucketName),
				Key:         aws.String(key),
				Body:        bytes.NewReader(buf),
				ContentType: aws.String(contentType),
			})

			if uploadErr == nil {
				// Public R2 URL
				currentURL = fmt.Sprintf("%s/%s", strings.TrimRight(publicURL, "/"), key)
			}
		}

		if currentURL == "" {
			// Fallback: save locally
			localDir := filepath.Join("uploads", "stories", storyID, "chapters", chapterNum)
			os.MkdirAll(localDir, 0755)
			localPath := filepath.Join(localDir, filename)
			if writeErr := os.WriteFile(localPath, buf, 0644); writeErr == nil {
				currentURL = fmt.Sprintf("/uploads/stories/%s/chapters/%s/%s", storyID, chapterNum, filename)
			}
		}

		if currentURL != "" {
			urls = append(urls, currentURL)
			uploadedFiles = append(uploadedFiles, fileInfo{url: currentURL, size: fileHeader.Size})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"urls":    urls,
		"count":   len(urls),
		"storage": map[bool]string{true: "cloudflare-r2", false: "local"}[useR2],
	})

	// Async record to DB to not slow down response
	go func(items []fileInfo, isR2 bool) {
		provider := "local"
		if isR2 {
			provider = "r2"
		}
		for _, item := range items {
			h.DB.Create(&models.Media{
				URL:      item.url,
				Size:     item.size,
				Provider: provider,
				FilePath: item.url,
			})
		}
	}(uploadedFiles, useR2)
}

// ListMedia returns all uploaded files
func (h *Handler) ListMedia(c *gin.Context) {
	var media []models.Media
	h.DB.Order("created_at desc").Find(&media)
	c.JSON(http.StatusOK, media)
}

// DeleteMedia deletes a record and its physical file
func (h *Handler) DeleteMedia(c *gin.Context) {
	id := c.Param("id")
	var media models.Media
	if err := h.DB.First(&media, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Media not found"})
		return
	}

	// 1. Delete physical file
	if media.Provider == "r2" {
		r2Client, err := newR2Client()
		if err == nil {
			// Extract key from URL
			publicURL := os.Getenv("R2_PUBLIC_URL")
			key := strings.TrimPrefix(media.URL, strings.TrimRight(publicURL, "/")+"/")
			r2Client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
				Bucket: aws.String(os.Getenv("R2_BUCKET_NAME")),
				Key:    aws.String(key),
			})
		}
	} else {
		// Delete local file
		localPath := strings.TrimPrefix(media.URL, "/uploads")
		os.Remove(filepath.Join("uploads", localPath))
	}

	// 2. Delete DB record
	h.DB.Delete(&media)
	c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}

// GetStorageStats returns the total storage used and current counts
func (h *Handler) GetStorageStats(c *gin.Context) {
	var stats struct {
		TotalSize  int64 `json:"totalSize"`
		TotalFiles int64 `json:"totalFiles"`
	}

	h.DB.Model(&models.Media{}).Select("SUM(size) as total_size, COUNT(*) as total_files").Scan(&stats)

	c.JSON(http.StatusOK, gin.H{
		"totalSize":  stats.TotalSize,
		"totalFiles": stats.TotalFiles,
		"limitBytes": 10 * 1024 * 1024 * 1024, // 10 GB
		"freeTier":   "10 GB",
	})
}
