package httpfile

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"strings"
)

// FileItem ...
type FileItem struct {
	FilePath string
	// application/octet-stream by default
	ContentType string
}

// UploadOptions ...
type UploadOptions struct {
	FileItems []FileItem
	TargetURL string
	Header    map[string]string
	// file by default
	FileField  string
	ExtraField map[string]string
	Stream     bool
}

// UploadResponse ...
type UploadResponse struct {
	Result     []byte
	Header     http.Header
	StatusCode int
}

// NewFileItem ...
func NewFileItem(filePath string, contentType ...string) FileItem {
	var ct string
	if len(contentType) > 0 {
		ct = contentType[0]
	}
	return FileItem{FilePath: filePath, ContentType: ct}
}

// NewFileItems ...
func NewFileItems(filePath string, contentType ...string) []FileItem {
	var ct string
	if len(contentType) > 0 {
		ct = contentType[0]
	}
	return []FileItem{FileItem{FilePath: filePath, ContentType: ct}}
}

// Upload single or multi file to file server by formdata
func Upload(opts UploadOptions) (*UploadResponse, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	if opts.FileField == "" {
		opts.FileField = "file"
	}
	for _, item := range opts.FileItems {
		flieNames := strings.Split(item.FilePath, "/")
		fileName := flieNames[len(flieNames)-1]
		var fileWriter io.Writer
		if item.ContentType == "" {
			item.ContentType = "application/octet-stream"
		}
		fileWriter, _ = createFormFile(bodyWriter, opts.FileField, fileName, item.ContentType)
		if fileWriter == nil {
			return nil, errors.New("error writing to buffer")
		}
		fh, err := os.Open(item.FilePath)
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(fileWriter, fh)
		if err != nil {
			return nil, err
		}
	}
	for key, val := range opts.ExtraField {
		_ = bodyWriter.WriteField(key, val)
	}
	bodyWriter.Close()
	request, err := http.NewRequest(http.MethodPost, opts.TargetURL, bodyBuf)
	if err != nil {
		return nil, err
	}
	for k, v := range opts.Header {
		request.Header.Set(k, v)
	}
	request.Header.Set("Content-Type", bodyWriter.FormDataContentType())
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	res := &UploadResponse{
		Result:     respBody,
		Header:     resp.Header,
		StatusCode: resp.StatusCode,
	}
	return res, err
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

// CreateFormFile is a convenience wrapper around CreatePart. It creates
// a new form-data header with the provided field name and file name.
func createFormFile(w *multipart.Writer, fieldname, filename string, contentType string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(fieldname), escapeQuotes(filename)))
	h.Set("Content-Type", contentType)
	return w.CreatePart(h)
}
func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

// UploadFile ...
func UploadFile(filePath string, targetURL string, Header ...map[string]string) (*UploadResponse, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return UploadReader(file, targetURL, Header...)
}

// UploadReader ...
func UploadReader(body io.Reader, targetURL string, Header ...map[string]string) (*UploadResponse, error) {
	request, err := http.NewRequest(http.MethodPost, targetURL, body)
	request.Header.Set("Content-Type", "binary/octet-stream")
	if len(Header) > 0 {
		for k, v := range Header[0] {
			request.Header.Set(k, v)
		}
	}
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	res := &UploadResponse{
		Result:     respBody,
		Header:     resp.Header,
		StatusCode: resp.StatusCode,
	}
	return res, err
}

// DownloadResponse ...
type DownloadResponse struct {
	FileSize   int64
	Header     http.Header
	StatusCode int
}

// Download ...
func Download(targetURL string, savePath string, Header ...map[string]string) (*DownloadResponse, error) {
	request, err := http.NewRequest(http.MethodGet, targetURL, nil)
	if len(Header) > 0 {
		for k, v := range Header[0] {
			request.Header.Set(k, v)
		}
	}
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	out, err := os.Create(savePath)
	if err != nil {
		return nil, err
	}
	defer out.Close()
	defer resp.Body.Close()
	n, err := io.Copy(out, resp.Body)
	res := &DownloadResponse{
		FileSize:   n,
		Header:     resp.Header,
		StatusCode: resp.StatusCode,
	}
	return res, err
}
