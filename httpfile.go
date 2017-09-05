package httpfile

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

var httpFile = New(nil)

// New ...
func New(client *http.Client) *HTTPFile {
	httpfile := &HTTPFile{client: client}
	if httpfile.client == nil {
		httpfile.client = &http.Client{}
	}
	return httpfile
}

// HTTPFile ...
type HTTPFile struct {
	client *http.Client
}

// Upload ...
func (h *HTTPFile) Upload(opts UploadOptions) (*UploadResponse, error) {
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
		fh.Close()
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
	resp, err := h.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	res := &UploadResponse{
		Result:     respBody,
		Res:        resp,
		Header:     resp.Header,
		StatusCode: resp.StatusCode,
	}
	return res, err
}

// UploadFile ...
func (h *HTTPFile) UploadFile(filePath string, targetURL string, Header ...map[string]string) (*UploadResponse, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return h.UploadReader(file, targetURL, Header...)
}

// UploadReader ...
func (h *HTTPFile) UploadReader(body io.Reader, targetURL string, Header ...map[string]string) (*UploadResponse, error) {
	request, err := http.NewRequest(http.MethodPost, targetURL, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "binary/octet-stream")
	if len(Header) > 0 {
		for k, v := range Header[0] {
			request.Header.Set(k, v)
		}
	}
	resp, err := h.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	res := &UploadResponse{
		Result:     respBody,
		Res:        resp,
		Header:     resp.Header,
		StatusCode: resp.StatusCode,
	}
	return res, err
}

// Download will get filename from 'Content-Disposition' if savePath is empty.
func (h *HTTPFile) Download(targetURL string, savePath string, Header ...map[string]string) (*DownloadResponse, error) {
	request, err := http.NewRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		return nil, err
	}
	if len(Header) > 0 {
		for k, v := range Header[0] {
			request.Header.Set(k, v)
		}
	}
	resp, err := h.client.Do(request)
	if err != nil {
		return nil, err
	}
	if savePath == "" {
		_, params, err := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))
		if err == nil {
			savePath = params["filename"]
		}
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
		Res:        resp,
		Header:     resp.Header,
		StatusCode: resp.StatusCode,
	}
	return res, err
}

// Head ...
func (h *HTTPFile) Head(targetURL string, Header ...map[string]string) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodHead, targetURL, nil)
	if err != nil {
		return nil, err
	}
	if len(Header) > 0 {
		for k, v := range Header[0] {
			request.Header.Set(k, v)
		}
	}
	return h.client.Do(request)
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
	Res        *http.Response
	Result     []byte
	Header     http.Header
	StatusCode int
}

// Upload single or multi file to file server by formdata
func Upload(opts UploadOptions) (*UploadResponse, error) {
	return httpFile.Upload(opts)
}

// UploadFile ...
func UploadFile(filePath string, targetURL string, Header ...map[string]string) (*UploadResponse, error) {
	return httpFile.UploadFile(filePath, targetURL, Header...)
}

// UploadReader ...
func UploadReader(body io.Reader, targetURL string, Header ...map[string]string) (*UploadResponse, error) {
	return httpFile.UploadReader(body, targetURL, Header...)
}

// DownloadResponse ...
type DownloadResponse struct {
	Res        *http.Response
	FileSize   int64
	Header     http.Header
	StatusCode int
}

// Download will get filename from 'Content-Disposition' if savePath is empty.
func Download(targetURL string, savePath string, Header ...map[string]string) (*DownloadResponse, error) {
	return httpFile.Download(targetURL, savePath, Header...)
}

// Head ...
func Head(targetURL string, Header ...map[string]string) (*http.Response, error) {
	return httpFile.Head(targetURL, Header...)
}
