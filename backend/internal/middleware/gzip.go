package middleware

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"strconv"
	"strings"
)

const gzipMinSizeBytes = 1024

type bufferedResponseWriter struct {
	headers     http.Header
	statusCode  int
	wroteHeader bool
	body        bytes.Buffer
}

func newBufferedResponseWriter() *bufferedResponseWriter {
	return &bufferedResponseWriter{
		headers: make(http.Header),
	}
}

func (w *bufferedResponseWriter) Header() http.Header {
	return w.headers
}

func (w *bufferedResponseWriter) WriteHeader(statusCode int) {
	if w.wroteHeader {
		return
	}
	w.wroteHeader = true
	w.statusCode = statusCode
}

func (w *bufferedResponseWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.body.Write(b)
}

// GzipMiddleware compresses API responses larger than 1KB when client supports gzip.
func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if shouldSkipGzip(r) || !clientAcceptsGzip(r) {
			next.ServeHTTP(rw, r)
			return
		}

		buffered := newBufferedResponseWriter()
		next.ServeHTTP(buffered, r)

		statusCode := buffered.statusCode
		if statusCode == 0 {
			statusCode = http.StatusOK
		}

		copyHeaders(rw.Header(), buffered.headers)
		appendVaryHeader(rw.Header(), "Accept-Encoding")

		allowBody := statusAllowsBody(statusCode) && r.Method != http.MethodHead
		body := buffered.body.Bytes()
		bodyLen := len(body)

		if allowBody && shouldCompressResponse(statusCode, rw.Header().Get("Content-Type"), rw.Header().Get("Content-Encoding"), bodyLen) {
			var compressed bytes.Buffer
			gz, err := gzip.NewWriterLevel(&compressed, gzip.BestSpeed)
			if err == nil {
				_, _ = gz.Write(body)
				_ = gz.Close()
				rw.Header().Set("Content-Encoding", "gzip")
				rw.Header().Set("Content-Length", strconv.Itoa(compressed.Len()))
				rw.WriteHeader(statusCode)
				_, _ = rw.Write(compressed.Bytes())
				return
			}
		}

		if statusAllowsBody(statusCode) {
			rw.Header().Set("Content-Length", strconv.Itoa(bodyLen))
		}
		rw.WriteHeader(statusCode)
		if allowBody {
			_, _ = rw.Write(body)
		}
	})
}

func shouldSkipGzip(r *http.Request) bool {
	if r == nil {
		return true
	}
	if !strings.HasPrefix(strings.TrimSpace(r.URL.Path), "/api/") {
		return true
	}
	if strings.Contains(strings.ToLower(strings.TrimSpace(r.Header.Get("Connection"))), "upgrade") {
		return true
	}
	return false
}

func clientAcceptsGzip(r *http.Request) bool {
	if r == nil {
		return false
	}
	acceptEncoding := strings.ToLower(strings.TrimSpace(r.Header.Get("Accept-Encoding")))
	return strings.Contains(acceptEncoding, "gzip")
}

func shouldCompressResponse(statusCode int, contentType string, contentEncoding string, bodyLen int) bool {
	if !statusAllowsBody(statusCode) {
		return false
	}
	if strings.TrimSpace(contentEncoding) != "" {
		return false
	}
	if bodyLen < gzipMinSizeBytes {
		return false
	}

	normalized := strings.ToLower(strings.TrimSpace(contentType))
	if normalized == "" {
		return true
	}

	return strings.HasPrefix(normalized, "application/json") ||
		strings.HasPrefix(normalized, "text/") ||
		strings.HasPrefix(normalized, "application/javascript") ||
		strings.HasPrefix(normalized, "application/xml")
}

func statusAllowsBody(statusCode int) bool {
	if statusCode >= 100 && statusCode < 200 {
		return false
	}
	return statusCode != http.StatusNoContent && statusCode != http.StatusNotModified
}

func copyHeaders(dst http.Header, src http.Header) {
	for key := range dst {
		dst.Del(key)
	}
	for key, values := range src {
		for _, value := range values {
			dst.Add(key, value)
		}
	}
}

func appendVaryHeader(header http.Header, value string) {
	existing := header.Values("Vary")
	for _, item := range existing {
		for _, part := range strings.Split(item, ",") {
			if strings.EqualFold(strings.TrimSpace(part), value) {
				return
			}
		}
	}
	header.Add("Vary", value)
}
