package adminecho

import (
	"fmt"
	"github.com/go-advanced-admin/admin"
	"github.com/labstack/echo/v4"
	"mime"
	"net/http"
	"path/filepath"
)

type Integrator struct {
	group *echo.Group
}

func NewIntegrator(g *echo.Group) *Integrator {
	return &Integrator{group: g}
}

func (i *Integrator) HandleRoute(method, path string, handler admin.HandlerFunc) {
	i.group.Add(method, path, func(c echo.Context) error {
		code, body := handler(c)
		return c.HTML(int(code), body)
	})
}

func (i *Integrator) ServeAssets(prefix string, renderer admin.TemplateRenderer) {
	i.group.GET(fmt.Sprintf("/%s/*", prefix), func(c echo.Context) error {
		fileName := c.Param("*")
		fileData, err := renderer.GetAsset(fileName)
		if err != nil {
			return c.NoContent(http.StatusNotFound)
		}

		contentType := mime.TypeByExtension(filepath.Ext(fileName))
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		return c.Blob(http.StatusOK, contentType, fileData)
	})
}
