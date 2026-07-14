package middleware

import (
	"log/slog"
	"mime"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/labstack/echo/v5"
)

// MarkdownNegotiation converts HTML responses when an agent requests Markdown.
func MarkdownNegotiation(next echo.HandlerFunc) echo.HandlerFunc {
	return echo.WrapMiddleware(markdownNegotiation)(next)
}

func markdownNegotiation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Accept")
		if !acceptsMarkdown(strings.Join(r.Header.Values("Accept"), ",")) {
			next.ServeHTTP(w, r)
			return
		}

		response := httptest.NewRecorder()
		next.ServeHTTP(response, r)

		contentType, _, err := mime.ParseMediaType(response.Header().Get("Content-Type"))
		if err != nil || !strings.EqualFold(contentType, "text/html") {
			writeBufferedResponse(w, response, response.Body.Bytes())
			return
		}

		markdown, err := htmltomarkdown.ConvertString(response.Body.String())
		if err != nil {
			slog.WarnContext(r.Context(), "failed to convert HTML response to Markdown", "error", err)
			writeBufferedResponse(w, response, response.Body.Bytes())
			return
		}

		body := []byte(markdown)
		header := response.Header()
		header.Set("Content-Type", "text/markdown; charset=utf-8")
		header.Set("Content-Length", strconv.Itoa(len(body)))
		header.Add("Vary", "Accept")
		for _, name := range []string{"Content-Encoding", "Content-Range", "Transfer-Encoding", "ETag", "Last-Modified"} {
			header.Del(name)
		}

		writeBufferedResponse(w, response, body)
	})
}

func acceptsMarkdown(accept string) bool {
	for _, value := range strings.Split(accept, ",") {
		mediaType, params, err := mime.ParseMediaType(strings.TrimSpace(value))
		if err != nil || !strings.EqualFold(mediaType, "text/markdown") {
			continue
		}

		quality, err := strconv.ParseFloat(params["q"], 64)
		if params["q"] == "" || (err == nil && quality > 0) {
			return true
		}
	}

	return false
}

func writeBufferedResponse(w http.ResponseWriter, response *httptest.ResponseRecorder, body []byte) {
	for name, values := range response.Header() {
		w.Header()[name] = values
	}
	w.WriteHeader(response.Code)
	_, _ = w.Write(body)
}
