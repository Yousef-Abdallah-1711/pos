package http

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
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
	symlinkFixtureAvailable := true
	if err := os.Symlink(outside, filepath.Join(folder, "escape.png")); err != nil {
		if runtime.GOOS != "windows" {
			t.Fatalf("create external symlink fixture: %v", err)
		}
		symlinkFixtureAvailable = false
		t.Logf("external-symlink assertion is exercised by Linux CI; this Windows host cannot create the fixture: %v", err)
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
	paths := []string{"products/1/2/image.png", "products/1/2/image.html", "../private.txt", "products/1/2/", "products/1/2/missing.png"}
	if symlinkFixtureAvailable {
		paths = append(paths, "products/1/2/escape.png")
	}
	for _, name := range paths {
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
