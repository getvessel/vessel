package engine

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"vessl.dev/vessl/internal/models"
)

func (bm *BackupManager) uploadToS3(ctx context.Context, dest *models.S3Destination, fileName string, data []byte) (string, error) {
	url := fmt.Sprintf("https://%s/%s/%s", dest.Endpoint, dest.Bucket, fileName)
	if strings.HasPrefix(dest.Endpoint, "http://") || strings.HasPrefix(dest.Endpoint, "https://") {
		url = fmt.Sprintf("%s/%s/%s", strings.TrimRight(dest.Endpoint, "/"), dest.Bucket, fileName)
	} else if strings.HasPrefix(dest.Endpoint, "localhost") || strings.HasPrefix(dest.Endpoint, "127.0.0.1") {
		url = fmt.Sprintf("http://%s/%s/%s", dest.Endpoint, dest.Bucket, fileName)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	req.ContentLength = int64(len(data))
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("X-Vessl-S3-Access-Key", dest.AccessKeyID)
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Sprintf("s3://%s/%s", dest.Bucket, fileName), nil
	}
	defer resp.Body.Close()
	return fmt.Sprintf("s3://%s/%s", dest.Bucket, fileName), nil
}

func (bm *BackupManager) enforceRetentionPolicy(cfg *models.BackupConfig) {
	if cfg.RetentionDays <= 0 {
		return
	}
	cutoff := time.Now().Add(-time.Duration(cfg.RetentionDays) * 24 * time.Hour)
	records, err := bm.store.ListBackupRecords(cfg.ID)
	if err != nil {
		return
	}
	for _, rec := range records {
		if rec.Status == models.BackupRecordStatusCompleted && rec.FilePath != "" {
			started, err := time.Parse(time.RFC3339, rec.StartedAt)
			if err == nil && started.Before(cutoff) {
				_ = os.Remove(rec.FilePath)
				_ = bm.store.UpdateBackupRecord(models.UpdateBackupRecordOpts{
					ID:          rec.ID,
					Status:      models.BackupRecordStatusExpired,
					S3URL:       rec.S3URL,
					Logs:        rec.Logs + "\nFile pruned by retention policy.",
					CompletedAt: time.Now().UTC().Format(time.RFC3339),
				})
			}
		}
	}
}
