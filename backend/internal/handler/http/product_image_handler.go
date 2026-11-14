package http

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"regexp"
)

// Only public raster catalog assets are served; filesystem containment also rejects symlinks escaping uploads.
var productImagePath = regexp.MustCompile(`^products/[1-9][0-9]*/[1-9][0-9]*/[A-Za-z0-9_-][A-Za-z0-9_.-]*\.(png|jpg|jpeg|gif|webp)$`)

type ProductImageHandler struct{ root *os.Root }

func NewProductImageHandler(directory string) (*ProductImageHandler, error) {
	if err := os.MkdirAll(directory, 0755); err != nil {
		return nil, fmt.Errorf("prepare image storage: %w", err)
	}
	root, err := os.OpenRoot(directory)
	if err != nil {
		return nil, fmt.Errorf("open image storage: %w", err)
	}
	return &ProductImageHandler{root: root}, nil
}

func (h *ProductImageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("path")
	if !productImagePath.MatchString(name) {
		http.NotFound(w, r)
		return
	}
	file, err := h.root.Open(name)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil || !info.Mode().IsRegular() || info.Size() > 10*1024*1024 {
		http.NotFound(w, r)
		return
	}
	types := map[string]string{".png": "image/png", ".jpg": "image/jpeg", ".jpeg": "image/jpeg", ".gif": "image/gif", ".webp": "image/webp"}
	w.Header().Set("Content-Type", types[path.Ext(name)])
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; sandbox")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	http.ServeContent(w, r, path.Base(name), info.ModTime(), file)
}
