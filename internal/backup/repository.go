package backup

import "context"

type Repository interface {
	CreateConfig(ctx context.Context, cfg *BackupConfig) error
	GetConfigByID(ctx context.Context, id string) (*BackupConfig, error)
	ListConfigsByProject(ctx context.Context, projectID string) ([]*BackupConfig, error)
	ListAllActiveConfigs(ctx context.Context) ([]*BackupConfig, error)
	DeleteConfig(ctx context.Context, id, projectID string) error

	CreateRecord(ctx context.Context, rec *BackupRecord) error
	ListRecordsByConfig(ctx context.Context, backupConfigID string) ([]*BackupRecord, error)

	CreateS3Destination(ctx context.Context, dest *S3Destination) error
	ListS3Destinations(ctx context.Context, projectID string) ([]*S3Destination, error)
	GetS3Destination(ctx context.Context, id string) (*S3Destination, error)
	DeleteS3Destination(ctx context.Context, id, projectID string) error
}
