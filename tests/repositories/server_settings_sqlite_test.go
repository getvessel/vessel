package repositories_test

import (
	"context"
	"database/sql"
	"testing"
	"vessl.dev/vessl/internal/repositories"

	_ "modernc.org/sqlite"
)

func TestSettingsRepositoryCreatesAndUpdatesDefaults(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:?_pragma=foreign_keys(ON)")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	if err := repositories.RunMigrations(db); err != nil {
		t.Fatalf("run migrations: %v", err)
	}

	repo := repositories.NewSettingsRepo(db)
	settings, err := repo.GetServerSettings(context.Background())
	if err != nil {
		t.Fatalf("get default settings: %v", err)
	}
	if settings.ID != "global" {
		t.Fatalf("expected global settings, got %q", settings.ID)
	}

	settings.SiteName = "Vessl Test"
	settings.RegistrationEnabled = false
	if err := repo.UpdateServerSettings(context.Background(), settings); err != nil {
		t.Fatalf("update settings: %v", err)
	}

	updated, err := repo.GetServerSettings(context.Background())
	if err != nil {
		t.Fatalf("get updated settings: %v", err)
	}
	if updated.SiteName != "Vessl Test" {
		t.Fatalf("expected updated site name, got %q", updated.SiteName)
	}
	if updated.RegistrationEnabled {
		t.Fatal("expected registration to be disabled")
	}
}
