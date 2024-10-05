package adminpanel

import (
	"fmt"
	"net/http"
)

// GetLogBaseLink returns the base URL path for logs.
func (ap *AdminPanel) GetLogBaseLink() string {
	return "/i/log"
}

// GetFullLogBaseLink returns the full URL path for logs, including the admin prefix.
func (ap *AdminPanel) GetFullLogBaseLink() string {
	return ap.Config.GetLink(ap.GetLogBaseLink())
}

// GetLogHandler returns the HTTP handler function for viewing a log entry.
func (ap *AdminPanel) GetLogHandler() HandlerFunc {
	return func(data interface{}) (uint, string) {
		instanceIDStr := ap.Web.GetPathParam(data, "id")
		if instanceIDStr == "" {
			return GetErrorHTML(http.StatusBadRequest, fmt.Errorf("instance id is required"))
		}

		entry, err := ap.Config.LogStore.GetLogEntry(instanceIDStr)
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}
		if entry == nil {
			return GetErrorHTML(http.StatusNotFound, fmt.Errorf("log entry not found"))
		}

		allowed, err := ap.PermissionChecker.HasLogViewPermission(data, entry.ID)
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}
		if !allowed {
			return GetErrorHTML(http.StatusForbidden, fmt.Errorf("you are not allowed to view this log entry"))
		}

		apps, err := GetAppsWithReadPermissions(ap, data)
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}

		html, err := ap.Config.Renderer.RenderTemplate("log", map[string]interface{}{
			"apps":        apps,
			"navBarItems": ap.Config.GetNavBarItems(data),
			"log":         entry,
		})
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}
		err = ap.CreateLogViewLog(data, *entry)
		if err != nil {
			return GetErrorHTML(http.StatusInternalServerError, err)
		}
		return http.StatusOK, html
	}
}
