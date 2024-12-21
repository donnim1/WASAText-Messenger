// cmd/webapi/register-web-ui-stub.go
//go:build !webui

package webapi

import (
	"net/http"
)

// registerWebUI is an empty stub because the `webui` build tag has not been specified.
// In Go, build tags allow you to include or exclude files during the build process.
// When the `webui` tag is not specified, this stub ensures that the application can
// compile without requiring the actual implementation for serving the web UI.
func registerWebUI(hdl http.Handler) (http.Handler, error) {
	// Simply return the provided handler unchanged, as no web UI needs to be registered.
	return hdl, nil
}
