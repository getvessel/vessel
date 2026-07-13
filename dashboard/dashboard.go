package dashboard

import (
	"net/http"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

func RegisterHandlers(e *echo.Echo) {
	e.GET("/*", func(c echo.Context) error {
		reqPath := filepath.Clean(c.Request().URL.Path)
		if reqPath == "/" || reqPath == "." {
			reqPath = "index.html"
		}

		content, err := DistFS.ReadFile("dist/" + reqPath)
		if err != nil {
			indexContent, err := DistFS.ReadFile("dist/index.html")
			if err != nil {
				return c.String(http.StatusNotFound, "Dashboard not built. Please run 'npm run build' in the dashboard directory.")
			}
			return c.HTMLBlob(http.StatusOK, indexContent)
		}

		contentType := http.DetectContentType(content)
		if filepath.Ext(reqPath) == ".css" {
			contentType = "text/css"
		} else if filepath.Ext(reqPath) == ".js" {
			contentType = "application/javascript"
		} else if filepath.Ext(reqPath) == ".svg" {
			contentType = "image/svg+xml"
		}
		return c.Blob(http.StatusOK, contentType, content)
	})
}
