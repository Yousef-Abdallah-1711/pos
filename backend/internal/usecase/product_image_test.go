package usecase

import (
	"bytes"
	"encoding/base64"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func imageHeader(t *testing.T, content []byte) *multipart.FileHeader {
	t.Helper()
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	part, err := writer.CreateFormFile("image", "../../untrusted.html")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = part.Write(content); err != nil {
		t.Fatal(err)
	}
	if err = writer.Close(); err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest("POST", "/", &buffer)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	if err = request.ParseMultipartForm(1024); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := request.MultipartForm.RemoveAll(); err != nil {
			t.Error(err)
		}
	})
	return request.MultipartForm.File["image"][0]
}

func TestProductImageUsesDetectedTypeAndGeneratedName(t *testing.T) {
	t.Chdir(t.TempDir())
	content, err := base64.StdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR4nGNw6TjzHwAFRAKYH64YhwAAAABJRU5ErkJggg==")
	if err != nil {
		t.Fatal(err)
	}
	uc := NewProductUseCase(nil, nil, "https://api.example.test/uploads")
	imageURL, err := uc.uploadProductImage(2, imageHeader(t, content), 1)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(imageURL, "https://api.example.test/uploads/products/1/2/") || !strings.HasSuffix(imageURL, ".png") || strings.Contains(imageURL, "untrusted") {
		t.Fatalf("unsafe generated URL %s", imageURL)
	}
	filename := filepath.Join("uploads", strings.TrimPrefix(imageURL, uc.storageURL+"/"))
	saved, err := os.ReadFile(filename)
	if err != nil || !bytes.Equal(saved, content) {
		t.Fatalf("image content changed: %v", err)
	}
	uc.discardPreparedImage(imageURL)
	if _, err = os.Stat(filename); !os.IsNotExist(err) {
		t.Fatalf("unreferenced image not removed: %v", err)
	}
}

func TestProductImageRejectsDisguisedExecutableAndOversize(t *testing.T) {
	uc := NewProductUseCase(nil, nil, "https://api.example.test/uploads")
	header := imageHeader(t, []byte("<html><script>malicious()</script></html>"))
	header.Header.Set("Content-Type", "image/png")
	if _, err := uc.uploadProductImage(2, header, 1); err == nil {
		t.Fatal("spoofed image was accepted")
	}
	header.Size = maxProductImageSize + 1
	if _, err := uc.uploadProductImage(2, header, 1); err == nil {
		t.Fatal("oversized image was accepted")
	}
}
