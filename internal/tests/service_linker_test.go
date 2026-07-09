package tests

import (
	"os"
	"reflect"
	"testing"

	"vessel.dev/vessel/internal/services"
	"vessel.dev/vessel/internal/store"
	"vessel.dev/vessel/internal/types"
)

func setupServiceLinkerTestStore(t *testing.T) (*store.Store, func()) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "vessel-servicelinker-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	s, err := store.NewStore(tmpDir)
	if err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("NewStore failed: %v", err)
	}

	cleanup := func() {
		s.Close()
		os.RemoveAll(tmpDir)
	}
	return s, cleanup
}

func TestServiceLinker_EmptyProjectID(t *testing.T) {
	s, cleanup := setupServiceLinkerTestStore(t)
	defer cleanup()

	linker := services.NewServiceLinker(s)
	envMap, err := linker.GetLinkedEnvironmentVariables("")
	if err != nil {
		t.Fatalf("expected nil error for empty projectID, got %v", err)
	}
	if len(envMap) != 0 {
		t.Errorf("expected empty map for empty projectID, got %v", envMap)
	}
}

func TestServiceLinker_NoLinkedServices(t *testing.T) {
	s, cleanup := setupServiceLinkerTestStore(t)
	defer cleanup()

	linker := services.NewServiceLinker(s)
	envMap, err := linker.GetLinkedEnvironmentVariables("proj-nonexistent")
	if err != nil {
		t.Fatalf("expected nil error for project with no linked services, got %v", err)
	}
	if len(envMap) != 0 {
		t.Errorf("expected empty map when no services linked, got %v", envMap)
	}
}

func TestServiceLinker_TableDrivenEngines(t *testing.T) {
	s, cleanup := setupServiceLinkerTestStore(t)
	defer cleanup()

	tests := []struct {
		name         string
		setupDB      *types.DatabaseConfig
		setupStorage *types.StorageConfig
		expectedEnvs map[string]string
	}{
		{
			name: "PostgreSQL Engine",
			setupDB: &types.DatabaseConfig{
				ProjectID:    "proj-pg",
				Name:         "my-pg",
				Engine:       "postgres",
				Username:     "pguser",
				Password:     "secretpassword",
				DatabaseName: "appdb",
				InternalDNS:  "postgres-my-pg.vessel-net",
			},
			expectedEnvs: map[string]string{
				"DATABASE_URL":      "postgresql://pguser:secretpassword@postgres-my-pg.vessel-net:5432/appdb",
				"POSTGRES_URL":      "postgresql://pguser:secretpassword@postgres-my-pg.vessel-net:5432/appdb",
				"POSTGRES_HOST":     "postgres-my-pg.vessel-net",
				"POSTGRES_PORT":     "5432",
				"POSTGRES_USER":     "pguser",
				"POSTGRES_PASSWORD": "secretpassword",
				"POSTGRES_DB":       "appdb",
			},
		},
		{
			name: "MySQL Engine",
			setupDB: &types.DatabaseConfig{
				ProjectID:    "proj-mysql",
				Name:         "my-mysql",
				Engine:       "mysql",
				Username:     "root",
				Password:     "mysqlpass",
				DatabaseName: "shopdb",
				InternalDNS:  "mysql-my-mysql.vessel-net",
			},
			expectedEnvs: map[string]string{
				"DATABASE_URL":   "root:mysqlpass@tcp(mysql-my-mysql.vessel-net:3306)/shopdb",
				"MYSQL_URL":      "root:mysqlpass@tcp(mysql-my-mysql.vessel-net:3306)/shopdb",
				"MYSQL_HOST":     "mysql-my-mysql.vessel-net",
				"MYSQL_PORT":     "3306",
				"MYSQL_USER":     "root",
				"MYSQL_PASSWORD": "mysqlpass",
				"MYSQL_DATABASE": "shopdb",
			},
		},
		{
			name: "Redis Engine",
			setupDB: &types.DatabaseConfig{
				ProjectID:   "proj-redis",
				Name:        "my-redis",
				Engine:      "redis",
				Password:    "redisauth",
				InternalDNS: "redis-my-redis.vessel-net",
			},
			expectedEnvs: map[string]string{
				"REDIS_URL":      "redis://:redisauth@redis-my-redis.vessel-net:6379",
				"REDIS_HOST":     "redis-my-redis.vessel-net",
				"REDIS_PORT":     "6379",
				"REDIS_PASSWORD": "redisauth",
			},
		},
		{
			name: "MongoDB Engine",
			setupDB: &types.DatabaseConfig{
				ProjectID:    "proj-mongo",
				Name:         "my-mongo",
				Engine:       "mongodb",
				Username:     "mongoadmin",
				Password:     "mongosecret",
				DatabaseName: "analytics",
				InternalDNS:  "mongodb-my-mongo.vessel-net",
			},
			expectedEnvs: map[string]string{
				"DATABASE_URL":   "mongodb://mongoadmin:mongosecret@mongodb-my-mongo.vessel-net:27017/analytics?authSource=admin",
				"MONGO_URL":      "mongodb://mongoadmin:mongosecret@mongodb-my-mongo.vessel-net:27017/analytics?authSource=admin",
				"MONGO_HOST":     "mongodb-my-mongo.vessel-net",
				"MONGO_PORT":     "27017",
				"MONGO_USER":     "mongoadmin",
				"MONGO_PASSWORD": "mongosecret",
				"MONGO_DB":       "analytics",
			},
		},
		{
			name: "MinIO Storage Engine",
			setupStorage: &types.StorageConfig{
				ProjectID:   "proj-s3",
				Name:        "my-s3",
				Type:        "minio",
				AccessKey:   "minioadmin",
				SecretKey:   "miniosecretkey",
				BucketName:  "uploads",
				InternalDNS: "minio-my-s3.vessel-net",
			},
			expectedEnvs: map[string]string{
				"S3_ENDPOINT":       "http://minio-my-s3.vessel-net:9000",
				"S3_ACCESS_KEY":     "minioadmin",
				"S3_SECRET_KEY":     "miniosecretkey",
				"S3_BUCKET":         "uploads",
				"MINIO_URL":         "http://minio-my-s3.vessel-net:9000",
				"MINIO_CONSOLE_URL": "http://minio-my-s3.vessel-net:9001",
			},
		},
		{
			name: "Unsupported/Unknown Engine Skipped",
			setupDB: &types.DatabaseConfig{
				ProjectID:    "proj-unknown",
				Name:         "my-unknown",
				Engine:       "cassandra",
				Username:     "admin",
				Password:     "pass",
				DatabaseName: "db",
				InternalDNS:  "cassandra.vessel-net",
			},
			expectedEnvs: map[string]string{},
		},
	}

	linker := services.NewServiceLinker(s)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectID := ""
			if tt.setupDB != nil {
				projectID = tt.setupDB.ProjectID
				if err := s.CreateDatabase(tt.setupDB); err != nil {
					t.Fatalf("CreateDatabase failed: %v", err)
				}
			}
			if tt.setupStorage != nil {
				projectID = tt.setupStorage.ProjectID
				if err := s.CreateStorage(tt.setupStorage); err != nil {
					t.Fatalf("CreateStorage failed: %v", err)
				}
			}

			envMap, err := linker.GetLinkedEnvironmentVariables(projectID)
			if err != nil {
				t.Fatalf("GetLinkedEnvironmentVariables failed: %v", err)
			}

			if !reflect.DeepEqual(envMap, tt.expectedEnvs) {
				t.Errorf("mismatch in env variables.\nGot:      %v\nExpected: %v", envMap, tt.expectedEnvs)
			}
		})
	}
}

func TestServiceLinker_CombinedServices(t *testing.T) {
	s, cleanup := setupServiceLinkerTestStore(t)
	defer cleanup()

	projectID := "proj-combined"

	pg := &types.DatabaseConfig{
		ProjectID:    projectID,
		Name:         "pg-db",
		Engine:       "postgres",
		Username:     "usr",
		Password:     "pwd",
		DatabaseName: "maindb",
		InternalDNS:  "pg.vessel-net",
	}
	if err := s.CreateDatabase(pg); err != nil {
		t.Fatalf("CreateDatabase pg failed: %v", err)
	}

	redis := &types.DatabaseConfig{
		ProjectID:   projectID,
		Name:        "cache",
		Engine:      "redis",
		Password:    "redispass",
		InternalDNS: "redis.vessel-net",
	}
	if err := s.CreateDatabase(redis); err != nil {
		t.Fatalf("CreateDatabase redis failed: %v", err)
	}

	s3 := &types.StorageConfig{
		ProjectID:   projectID,
		Name:        "assets",
		Type:        "minio",
		AccessKey:   "ak",
		SecretKey:   "sk",
		BucketName:  "media",
		InternalDNS: "minio.vessel-net",
	}
	if err := s.CreateStorage(s3); err != nil {
		t.Fatalf("CreateStorage s3 failed: %v", err)
	}

	linker := services.NewServiceLinker(s)
	envMap, err := linker.GetLinkedEnvironmentVariables(projectID)
	if err != nil {
		t.Fatalf("GetLinkedEnvironmentVariables failed: %v", err)
	}

	// Verify Postgres keys exist and are correct
	if envMap["POSTGRES_URL"] != "postgresql://usr:pwd@pg.vessel-net:5432/maindb" {
		t.Errorf("unexpected POSTGRES_URL: %s", envMap["POSTGRES_URL"])
	}
	// Verify Redis keys exist and are correct
	if envMap["REDIS_URL"] != "redis://:redispass@redis.vessel-net:6379" {
		t.Errorf("unexpected REDIS_URL: %s", envMap["REDIS_URL"])
	}
	// Verify MinIO S3 keys exist and are correct
	if envMap["S3_BUCKET"] != "media" || envMap["MINIO_URL"] != "http://minio.vessel-net:9000" {
		t.Errorf("unexpected S3 keys: bucket=%s, url=%s", envMap["S3_BUCKET"], envMap["MINIO_URL"])
	}
}
