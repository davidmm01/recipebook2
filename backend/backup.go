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

var backupBucketName string

// InitBackups starts the background backup goroutine
func InitBackups(bucket string) {
	backupBucketName = bucket
	if backupBucketName == "" {
		log.Println("Backups disabled: no backup bucket configured")
		return
	}

	go runBackupLoop()
}

// runBackupLoop runs the backup process every hour
func runBackupLoop() {
	// Run first backup after a short delay to let the server start
	time.Sleep(1 * time.Minute)

	for {
		if err := performBackup(context.Background()); err != nil {
			log.Printf("Backup failed: %v", err)
		}

		// Wait 6 hours before next backup
		time.Sleep(6 * time.Hour)
	}
}

// performBackup creates a backup if the database has changed
func performBackup(ctx context.Context) error {
	log.Println("Starting backup check...")

	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create storage client: %w", err)
	}
	defer client.Close()

	// Calculate hash of current database
	currentHash, err := calculateDBHash()
	if err != nil {
		return fmt.Errorf("failed to calculate database hash: %w", err)
	}
	log.Printf("Current database hash: %s", currentHash)

	// Get hash of last backup
	_, lastHash, err := getLastBackupHash(ctx, client)
	if err == nil && lastHash == currentHash {
		log.Printf("Database unchanged (hash: %s), skipping backup", currentHash)
		return nil
	}

	// Database has changed or no previous backup exists, create new backup
	timestamp := time.Now().UTC().Format("2006-01-02-15-04-05")
	backupName := fmt.Sprintf("recipes-backup-%s.db", timestamp)

	log.Printf("Creating backup: %s", backupName)

	// Read database file
	dbMutex.RLock()
	data, err := os.ReadFile(localDBPath)
	dbMutex.RUnlock()
	if err != nil {
		return fmt.Errorf("failed to read database file: %w", err)
	}

	// Upload backup
	backupObj := client.Bucket(backupBucketName).Object(backupName)
	writer := backupObj.NewWriter(ctx)
	writer.ContentType = "application/x-sqlite3"
	writer.Metadata = map[string]string{
		"hash":      currentHash,
		"timestamp": timestamp,
		"source":    dbFileName,
	}

	if _, err := writer.Write(data); err != nil {
		writer.Close()
		return fmt.Errorf("failed to write backup: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close backup writer: %w", err)
	}

	log.Printf("Backup created successfully: %s (hash: %s, size: %d bytes)", backupName, currentHash, len(data))
	return nil
}

// calculateDBHash calculates SHA256 hash of the current database file
func calculateDBHash() (string, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	f, err := os.Open(localDBPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// getLastBackupHash retrieves the hash of the most recent backup
func getLastBackupHash(ctx context.Context, client *storage.Client) (string, string, error) {
	query := &storage.Query{Prefix: "recipes-backup-"}

	var lastObj *storage.ObjectAttrs
	it := client.Bucket(backupBucketName).Objects(ctx, query)

	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return "", "", err
		}

		if lastObj == nil || attrs.Created.After(lastObj.Created) {
			lastObj = attrs
		}
	}

	if lastObj == nil {
		return "", "", fmt.Errorf("no previous backups found")
	}

	hash, ok := lastObj.Metadata["hash"]
	if !ok {
		return lastObj.Name, "", fmt.Errorf("no hash in metadata")
	}

	return lastObj.Name, hash, nil
}
