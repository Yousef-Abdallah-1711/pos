package http

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestProductImagesAreContainedRasterAssets(t *testing.T) {
	directory := t.TempDir()
	folder := filepath.Join(directory, "products", "1", "2")
	if err := os.MkdirAll(folder, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(folder, "image.png"), []byte("synthetic-image"), 0644); err != nil {
		t.Fatal(err)
	}
	outside := filepath.Join(t.TempDir(), "private.txt")
	if err := os.WriteFile(outside, []byte("private"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(folder, "escape.png")); err != nil {
		t.Fatal(err)
	}
	handler, err := NewProductImageHandler(directory)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := handler.root.Close(); err != nil {
			t.Error(err)
		}
	})
	for _, name := range []string{"products/1/2/image.png", "products/1/2/escape.png", "products/1/2/image.html", "../private.txt", "products/1/2/", "products/1/2/missing.png"} {
		request := httptest.NewRequest(http.MethodGet, "/uploads/test", nil)
		request.SetPathValue("path", name)
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		if name != "products/1/2/image.png" {
			if response.Code != 404 {
				t.Errorf("%s: want 404 got %d", name, response.Code)
			}
			continue
		}
		if response.Code != 200 || response.Header().Get("Content-Type") != "image/png" || response.Header().Get("X-Content-Type-Options") != "nosniff" || response.Body.String() != "synthetic-image" {
			t.Fatalf("unsafe image response: status=%d headers=%v", response.Code, response.Header())
		}
	}
}
