package usecase

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

const maxProductImageSize = 10 * 1024 * 1024

// Files receive server-generated names and extensions from their actual content.
func (uc *ProductUseCase) uploadProductImage(restaurantID int, header *multipart.FileHeader, tenantID int) (string, error) {
	if tenantID <= 0 || restaurantID <= 0 || header.Size > maxProductImageSize {
		return "", fmt.Errorf("invalid image scope or file size exceeds 10MB limit")
	}
	file, err := header.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open image: %w", err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	prefix, err := reader.Peek(512)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("failed to inspect image: %w", err)
	}
	extensions := map[string]string{"image/jpeg": ".jpg", "image/png": ".png", "image/gif": ".gif", "image/webp": ".webp"}
	extension, valid := extensions[http.DetectContentType(prefix)]
	if !valid {
		return "", fmt.Errorf("invalid image content: allowed jpeg, png, gif, webp")
	}
	relative := fmt.Sprintf("products/%d/%d/%s%s", tenantID, restaurantID, uuid.NewString(), extension)
	if err := writeProductImage(filepath.Join("uploads", relative), reader); err != nil {
		return "", err
	}
	return strings.TrimRight(uc.storageURL, "/") + "/" + relative, nil
}

func writeProductImage(filename string, source io.Reader) error {
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return fmt.Errorf("failed to create image directory: %w", err)
	}
	output, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return fmt.Errorf("failed to create image: %w", err)
	}
	size, copyErr := io.Copy(output, io.LimitReader(source, maxProductImageSize+1))
	closeErr := output.Close()
	if copyErr != nil || closeErr != nil || size > maxProductImageSize {
		if err := os.Remove(filename); err != nil {
			log.Printf("failed to remove incomplete product image: %v", err)
		}
		return fmt.Errorf("failed to save valid image within 10MB limit")
	}
	return nil
}

// Delete only a newly prepared image if its database write fails; existing assets remain.
func (uc *ProductUseCase) discardPreparedImage(imageURL string) {
	prefix := strings.TrimRight(uc.storageURL, "/") + "/products/"
	if !strings.HasPrefix(imageURL, prefix) {
		return
	}
	relative := strings.TrimPrefix(imageURL, strings.TrimRight(uc.storageURL, "/")+"/")
	if err := os.Remove(filepath.Join("uploads", relative)); err != nil {
		log.Printf("failed to remove unreferenced product image: %v", err)
	}
}
