package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const scalarHTML = `<!doctype html>
<html>
  <head>
    <title>Hadith API Go - Documentation</title>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <style>
      body { margin: 0; padding: 0; }
    </style>
  </head>
  <body>
    <script id="api-reference" data-url="/openapi.yaml"></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference@latest/dist/browser/standalone.js"></script>
  </body>
</html>
`

type DocsHandler struct{}

func NewDocsHandler() *DocsHandler {
	return &DocsHandler{}
}

func (h *DocsHandler) ServeDocs(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, scalarHTML)
}

func (h *DocsHandler) ServeOpenAPI(c *gin.Context) {
	c.Header("Content-Type", "text/yaml; charset=utf-8")
	c.Header("Access-Control-Allow-Origin", "*")
	c.File("./docs/openapi.yaml")
}
