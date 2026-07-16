package utils

import "os"

// IsDryRun returns true if the DEPLOY_DRY_RUN environment variable is set to "true".
func IsDryRun() bool {
	return os.Getenv("DEPLOY_DRY_RUN") == "true"
}
