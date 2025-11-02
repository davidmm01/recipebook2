package main

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
)

const (
	maxImageSize     = 10 << 20 // 10MB
	imagePathPrefix  = "recipe-images"
	allowedImageExts = ".jpg,.jpeg,.png,.gif,.webp"

	maxIconSize     = 2 << 20 // 2MB
	iconPathPrefix  = "recipe-icons"
	allowedIconExts = ".jpg,.jpeg,.png,.svg,.webp"
)

// UploadImageToGCS uploads an image file to Google Cloud Storage
func UploadImageToGCS(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader, recipeID string) (string, error) {
	// Validate file size
	if fileHeader.Size > maxImageSize {
		return "", fmt.Errorf("file size exceeds maximum of 10MB")
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !strings.Contains(allowedImageExts, ext) {
		return "", fmt.Errorf("invalid file type: %s (allowed: %s)", ext, allowedImageExts)
	}

	// Generate unique filename
	filename := fmt.Sprintf("%s/%s/%s%s", imagePathPrefix, recipeID, uuid.New().String(), ext)

	// Create GCS client
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create storage client: %w", err)
	}
	defer client.Close()

	// Create GCS writer
	wc := client.Bucket(bucketName).Object(filename).NewWriter(ctx)
	wc.ContentType = getContentType(ext)

	// Copy file to GCS
	if _, err = io.Copy(wc, file); err != nil {
		wc.Close()
		return "", fmt.Errorf("failed to write file to storage: %w", err)
	}

	// Close writer (uploads the file)
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("failed to close storage writer: %w", err)
	}

	// Return public URL
	// Format: https://storage.googleapis.com/{bucket}/{object}
	imageURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, filename)
	return imageURL, nil
}

// UploadIconToGCS uploads an icon file to Google Cloud Storage
func UploadIconToGCS(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader) (string, string, error) {
	// Validate file size
	if fileHeader.Size > maxIconSize {
		return "", "", fmt.Errorf("file size exceeds maximum of 2MB")
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !strings.Contains(allowedIconExts, ext) {
		return "", "", fmt.Errorf("invalid file type: %s (allowed: %s)", ext, allowedIconExts)
	}

	// Generate unique filename
	uniqueFilename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	objectPath := fmt.Sprintf("%s/%s", iconPathPrefix, uniqueFilename)

	// Create GCS client
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", "", fmt.Errorf("failed to create storage client: %w", err)
	}
	defer client.Close()

	// Create GCS writer
	wc := client.Bucket(bucketName).Object(objectPath).NewWriter(ctx)
	wc.ContentType = getContentType(ext)

	// Copy file to GCS
	if _, err = io.Copy(wc, file); err != nil {
		wc.Close()
		return "", "", fmt.Errorf("failed to write file to storage: %w", err)
	}

	// Close writer (uploads the file)
	if err := wc.Close(); err != nil {
		return "", "", fmt.Errorf("failed to close storage writer: %w", err)
	}

	// Return both the filename and public URL
	iconURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectPath)
	return uniqueFilename, iconURL, nil
}

// DeleteImageFromGCS deletes an image from Google Cloud Storage
func DeleteImageFromGCS(ctx context.Context, imageURL string) error {
	// Extract object path from URL
	// URL format: https://storage.googleapis.com/{bucket}/{object}
	parts := strings.Split(imageURL, "/")
	if len(parts) < 5 {
		return fmt.Errorf("invalid image URL format")
	}

	// Get object path (everything after bucket name)
	objectPath := strings.Join(parts[4:], "/")

	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create storage client: %w", err)
	}
	defer client.Close()

	// Delete object
	if err := client.Bucket(bucketName).Object(objectPath).Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete image from storage: %w", err)
	}

	return nil
}

// getContentType returns the MIME type for an image extension
func getContentType(ext string) string {
	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".svg":
		return "image/svg+xml"
	default:
		return "application/octet-stream"
	}
}
