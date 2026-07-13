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

	"vessl.dev/vessl/dashboard"
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

	// Embed the React Dashboard if serving statically
	if os.Getenv("VESSL_STATIC_DIR") != "" {
		e.Static("/*", os.Getenv("VESSL_STATIC_DIR"))
	} else {
		dashboard.RegisterHandlers(e)
	}

	log.Printf(" Vessl Cloud SaaS listening on %s", addr)
	if err := http.ListenAndServe(addr, e); err != nil {
		log.Fatalf(" Server crashed: %v", err)
	}
}
