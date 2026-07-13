// @title Vessl Cloud API
// @version 1.0
// @description Vessl SaaS Cloud API
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	cloudhttp "vessl.dev/vessl/internal/cloud/http"
)

func main() {
	_ = godotenv.Load()

	log.Printf(" Booting Vessl Cloud SaaS...")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	host := os.Getenv("HOST")
	addr := host + ":" + port

	// Initialize the Echo Router
	e := echo.New()

	// Mount the Cloud SaaS API routes
	cloudhttp.MountCloudRoutes(e)

	// Note: In Cloud Mode, the dashboard is deployed separately (e.g., to Vercel or Cloudflare Pages)
	// so the Cloud API server does not need to serve the React assets itself.

	log.Printf(" Vessl Cloud SaaS listening on %s", addr)
	if err := http.ListenAndServe(addr, e); err != nil {
		log.Fatalf(" Server crashed: %v", err)
	}
}
