package repos

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

// SeedAdminUser inserts a default admin user if CLOUD_ADMIN_EMAIL and CLOUD_ADMIN_PASSWORD are set.
func SeedAdminUser(db *sql.DB) error {
	adminEmail := os.Getenv("CLOUD_ADMIN_EMAIL")
	adminPassword := os.Getenv("CLOUD_ADMIN_PASSWORD")

	if adminEmail == "" || adminPassword == "" {
		return nil
	}

	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM cloud_users WHERE email = $1`, adminEmail).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check admin user existence: %w", err)
	}
	if count > 0 {
		log.Printf("Admin user %s already exists, skipping seed.", adminEmail)
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash admin password: %w", err)
	}

	id := "admin-" + adminEmail
	_, err = db.Exec(
		`INSERT INTO cloud_users (id, email, full_name, password_hash, role, email_verified) VALUES ($1, $2, $3, $4, $5, $6)`,
		id, adminEmail, "Admin", string(hash), "admin", true,
	)
	if err != nil {
		return fmt.Errorf("failed to insert admin user: %w", err)
	}

	log.Printf("Seeded initial admin user: %s", adminEmail)
	return nil
}
