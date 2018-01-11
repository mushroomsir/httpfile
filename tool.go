package httpfile

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
)

var defaultHTTPClient = &http.Client{}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

// CreateFormFile is a convenience wrapper around CreatePart. It creates
// a new form-data header with the provided field name and file name.
func createFormFile(w *multipart.Writer, filename string, contentType string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, escapeQuotes(filename)))
	h.Set("Content-Type", contentType)
	return w.CreatePart(h)
}
func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}
