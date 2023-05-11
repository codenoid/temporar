package main

import (
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
)

var buildDir = "frontend/build"

func svelteStaticLoader(c *gin.Context) {
	trypath := path.Join(buildDir, c.Request.URL.Path)
	if _, err := f.ReadFile(trypath); err == nil {
		c.Status(http.StatusOK)
		c.FileFromFS(trypath, http.FS(f))
		return
	}

	index := path.Join(buildDir, "index.html")
	if b, err := f.ReadFile(index); err == nil {
		c.Status(http.StatusOK)
		c.Header("Content-Type", "text/html")
		c.Writer.Write(b)
	}
}
