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
		if code == http.StatusFound || code == http.StatusMovedPermanently || code == http.StatusSeeOther {
			return c.Redirect(int(code), body)
		}
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

func (i *Integrator) GetQueryParam(ctx interface{}, name string) string {
	ec := ctx.(echo.Context)
	return ec.QueryParam(name)
}

func (i *Integrator) GetPathParam(ctx interface{}, name string) string {
	ec := ctx.(echo.Context)
	return ec.Param(name)
}

func (i *Integrator) GetRequestMethod(ctx interface{}) string {
	ec := ctx.(echo.Context)
	return ec.Request().Method
}

func (i *Integrator) GetFormData(ctx interface{}) map[string][]string {
	ec := ctx.(echo.Context)
	if err := ec.Request().ParseForm(); err != nil {
		return nil
	}
	return ec.Request().Form
}
