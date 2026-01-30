package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

func main() {
	dbBucket := os.Getenv("DB_BUCKET_NAME")
	backupsBucket := os.Getenv("BACKUPS_BUCKET_NAME")

	if dbBucket == "" || backupsBucket == "" {
		log.Fatal("DB_BUCKET_NAME and BACKUPS_BUCKET_NAME environment variables are required")
	}

	ctx := context.Background()
	if err := run(ctx, dbBucket, backupsBucket); err != nil {
		log.Fatalf("Backup failed: %v", err)
	}
}

func run(ctx context.Context, dbBucket, backupsBucket string) error {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("creating storage client: %w", err)
	}
	defer client.Close()

	// Download recipes.db to a temp file
	tmpFile, err := os.CreateTemp("", "recipes-*.db")
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	log.Printf("Downloading recipes.db from gs://%s/recipes.db", dbBucket)
	reader, err := client.Bucket(dbBucket).Object("recipes.db").NewReader(ctx)
	if err != nil {
		return fmt.Errorf("opening recipes.db from bucket: %w", err)
	}

	if _, err := io.Copy(tmpFile, reader); err != nil {
		reader.Close()
		return fmt.Errorf("downloading recipes.db: %w", err)
	}
	reader.Close()

	// Calculate SHA256 hash
	if _, err := tmpFile.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("seeking temp file: %w", err)
	}

	hasher := sha256.New()
	if _, err := io.Copy(hasher, tmpFile); err != nil {
		return fmt.Errorf("hashing temp file: %w", err)
	}
	currentHash := hex.EncodeToString(hasher.Sum(nil))
	log.Printf("Current database hash: %s", currentHash)

	// Check latest backup hash
	lastHash, err := getLatestBackupHash(ctx, client, backupsBucket)
	if err == nil && lastHash == currentHash {
		log.Printf("Database unchanged (hash: %s), skipping backup", currentHash)
		return nil
	}
	if err != nil {
		log.Printf("No previous backup found or error reading it: %v", err)
	}

	// Upload new backup
	timestamp := time.Now().UTC().Format("2006-01-02-15-04-05")
	backupName := fmt.Sprintf("recipes-backup-%s.db", timestamp)

	if _, err := tmpFile.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("seeking temp file for upload: %w", err)
	}

	log.Printf("Uploading backup: gs://%s/%s", backupsBucket, backupName)
	writer := client.Bucket(backupsBucket).Object(backupName).NewWriter(ctx)
	writer.ContentType = "application/x-sqlite3"
	writer.Metadata = map[string]string{
		"hash":      currentHash,
		"timestamp": timestamp,
		"source":    "recipes.db",
	}

	written, err := io.Copy(writer, tmpFile)
	if err != nil {
		writer.Close()
		return fmt.Errorf("writing backup: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("closing backup writer: %w", err)
	}

	log.Printf("Backup created successfully: %s (hash: %s, size: %d bytes)", backupName, currentHash, written)
	return nil
}

func getLatestBackupHash(ctx context.Context, client *storage.Client, bucket string) (string, error) {
	query := &storage.Query{Prefix: "recipes-backup-"}

	var latest *storage.ObjectAttrs
	it := client.Bucket(bucket).Objects(ctx, query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return "", fmt.Errorf("listing backups: %w", err)
		}
		if latest == nil || attrs.Created.After(latest.Created) {
			latest = attrs
		}
	}

	if latest == nil {
		return "", fmt.Errorf("no previous backups found")
	}

	hash, ok := latest.Metadata["hash"]
	if !ok {
		return "", fmt.Errorf("latest backup %s has no hash metadata", latest.Name)
	}

	return hash, nil
}
