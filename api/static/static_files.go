package static

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"
)

func StaticHandler(root http.FileSystem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(ctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	}
}
